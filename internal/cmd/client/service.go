package client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	"github.com/zbindenren/sfmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/labels"
)

type serviceCmd struct {
	List       serviceList       `cmd:"" help:"List registered services."`
	Register   serviceRegister   `cmd:"" help:"Register a service."`
	UnRegister serviceUnRegister `cmd:"" help:"Unregister a service by ID or endpoint URL." name:"unregister"`
	Import     serviceImport     `cmd:"" help:"Import services from yaml file."`
}

type serviceList struct {
	Output        string `short:"o" default:"table" help:"Output formats. Valid formats: json, yaml, csv, table."`
	SortBy        string `default:"endpoint" help:"Sort services by endpoint or modification date (allowed values: endpoint or date)" enum:"endpoint,date"`
	Headers       bool   `short:"H" help:"Show headers."`
	UnResolved    bool   `short:"u" help:"List only services that cannot be resolved by the local resolver." name:"unresolved"`
	Namespace     string `short:"n" help:"Filter services by namespace."`
	serviceFilter `prefix:"filter-"`
}

func (s serviceList) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.ListService(ctx, &discoveryv1.ListServiceRequest{
		Namespace: s.Namespace,
	})
	if err != nil {
		return err
	}

	services := convert.ServicesFromPB(r.GetServices())

	filters, err := s.serviceFilter.filters()
	if err != nil {
		return err
	}

	if len(filters) > 0 {
		services = services.Filter(filters...)
	}

	switch s.SortBy {
	case "endpoint":
		services.SortByEndpoint()
	case "date":
		services.SortByDate()
	default:
		return fmt.Errorf("unsupported sort type '%s'", s.SortBy)
	}

	if s.UnResolved {
		unresolved, err := services.UnResolved()
		if err != nil {
			return err
		}

		services = unresolved
	}

	sw := sfmt.SliceWriter{
		Writer:    os.Stdout,
		NoHeaders: !s.Headers,
	}
	f := sfmt.ParseFormat(s.Output)

	sw.Write(f, services)

	return nil
}

type serviceRegister struct {
	Endpoints []string         `short:"e" help:"The service endpoint URLs." required:"true"`
	Name      string           `arg:"true" optional:"true" help:"The service name. This will represent the job name in prometheus." env:"DISCOVERY_NAME"`
	Labels    discovery.Labels `short:"l" help:"Labels for the service." mapsep:","`
	Namespace string           `short:"n" help:"The namespace for the service" default:"default" required:"true"`
	Selector  string           `short:"s" help:"Kubernetes style selectors (key=value) to select servers with specific labels."`
}

func (s serviceRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	if s.Name == "" {
		return errors.New("name cannot be empty")
	}

	var lastErr error

	for _, ep := range s.Endpoints {
		ctx, cancel := g.ctx()
		defer cancel()

		lbls := s.Labels

		r, err := cli.RegisterService(ctx, &discoveryv1.RegisterServiceRequest{
			Name:      s.Name,
			Namespace: s.Namespace,
			Endpoint:  ep,
			Labels:    lbls,
			Selector:  s.Selector,
		})
		if err != nil {
			lastErr = err
			l.Errorw("failed to register", "service", ep, "err", err)

			continue
		}

		l.Infow("service registered", "id", r.GetService().GetId())
	}

	if lastErr != nil {
		return errors.New("register service failed")
	}

	return nil
}

type serviceUnRegister struct {
	Endpoints  []string `arg:"true" optional:"true" help:"The service endpoint URLs or IDs." env:"DISCOVERY_ENDPOINTS"`
	Namespace  string   `short:"n" help:"The namespace for the service" xor:"X"`
	UnResolved bool     `short:"u" help:"Unregister all services that cannot be resolved by the local resolver." xor:"X" name:"unresolved" group:"Unresolved"`
	Percent    float64  `short:"p" help:"The maximum number of unresolved services in percent that are allowed to be unregistered at once." default:"1.0" group:"Unresolved"`
}

func (s serviceUnRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	if s.UnResolved && len(s.Endpoints) > 0 {
		return fmt.Errorf("unresolved cannot be used in combination with endpoints")
	}

	var lastErr error

	if s.UnResolved {
		return s.unRegisterUnresolved(g, l, cli)
	}

	if len(s.Endpoints) > 0 && s.Namespace == "" {
		return fmt.Errorf("no namespace specified")
	}

	for _, ep := range s.Endpoints {
		ctx, cancel := g.ctx()
		defer cancel()

		_, err = cli.UnRegisterService(ctx, &discoveryv1.UnRegisterServiceRequest{
			Namespace: s.Namespace,
			Id:        ep,
		})

		if err != nil {
			lastErr = err
			l.Errorw("failed to unregister", "service", ep, "namespace", s.Namespace, "err", err)

			continue
		}
	}

	if lastErr != nil {
		return errors.New("unregister service failed")
	}

	return nil
}

func (s serviceUnRegister) unRegisterUnresolved(g *Globals, l *zap.SugaredLogger, cli discoveryv1.ServiceAPIClient) error {
	var lastErr error

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.ListService(ctx, &discoveryv1.ListServiceRequest{})
	if err != nil {
		return err
	}

	services := convert.ServicesFromPB(r.GetServices())

	unresolved, err := services.UnResolved()
	if err != nil {
		return err
	}

	if len(unresolved) == 0 {
		l.Info("no unresolved services found")
	}

	unResolvedPercent := float64(len(unresolved)) / float64(len(services)) * 100

	if unResolvedPercent > s.Percent {
		return fmt.Errorf("maximum %f%% services are allowed to be unregistered", s.Percent)
	}

	for i := range unresolved {
		svc := unresolved[i]

		l.Debugw("unregister service", svc.KeyVals()...)

		ctx, cancel := g.ctx()
		defer cancel()

		_, err = cli.UnRegisterService(ctx, &discoveryv1.UnRegisterServiceRequest{
			Namespace: svc.Namespace,
			Id:        svc.Endpoint.String(),
		})

		if err != nil {
			lastErr = err
			l.Errorw("failed to unregister", "service", svc.Endpoint.String(), "namespace", svc.Namespace, "err", err)

			continue
		}
	}

	if lastErr != nil {
		return errors.New("unregister service failed")
	}

	return nil
}

type serviceFilter struct {
	Name     string `short:"N" help:"Filter services by job name (regular expression)."`
	Server   string `short:"S" help:"Filter services by server name (regular expression)."`
	Endpoint string `short:"e" help:"Filter services by endpoint (regular expression)."`
	Selector string `short:"s" help:"Filter services by label selector."`
}

func (s serviceFilter) filters() ([]discovery.FilterFunc, error) {
	filters := []discovery.FilterFunc{}

	if s.Name != "" {
		r, err := regexp.Compile(s.Name)
		if err != nil {
			return nil, err
		}

		filters = append(filters, discovery.ServiceByName(r))
	}

	if s.Endpoint != "" {
		r, err := regexp.Compile(s.Endpoint)
		if err != nil {
			return nil, err
		}

		filters = append(filters, discovery.ServiceByEndpoint(r))
	}

	if s.Server != "" {
		r, err := regexp.Compile(s.Server)
		if err != nil {
			return nil, err
		}

		filters = append(filters, discovery.ServiceByServer(r))
	}

	if s.Selector != "" {
		sel, err := labels.Parse(s.Selector)
		if err != nil {
			return nil, err
		}

		filters = append(filters, discovery.ServiceBySelector(sel))
	}

	return filters, nil
}

type serviceImport struct {
	Path string `arg:"true" help:"Path to yaml file." required:"true"`
}

func (s serviceImport) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	d, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}

	services := discovery.Services{}

	err = yaml.Unmarshal(d, &services)
	if err != nil {
		return err
	}

	failed := discovery.Services{}

	for j := range services {
		s := services[j]
		l.Debugw("import serivce", s.KeyVals()...)

		_, err := cli.RegisterService(ctx, &discoveryv1.RegisterServiceRequest{
			Name:        s.Name,
			Endpoint:    s.Endpoint.String(),
			Description: s.Description,
			Labels:      s.Labels,
			Namespace:   s.Namespace,
			Selector:    s.Selector,
		})

		if err != nil {
			failed = append(failed, s)
			msg := s.KeyVals()
			msg = append(msg, "err", err)
			l.Errorw("failed to import", msg...)

			continue
		}
	}

	if len(failed) > 0 {
		file, err := ioutil.TempFile("", "discovery-import-failed*.yaml")
		if err != nil {
			return fmt.Errorf("failed to create tmp file: %w", err)
		}

		l.Infow("saving failed to import services", "path", file.Name())

		d, err := yaml.Marshal(failed)
		if err != nil {
			return fmt.Errorf("failed to marshal failed services: %w", err)
		}

		if err := ioutil.WriteFile(file.Name(), d, 0600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Name(), err)
		}
	}

	return nil
}
