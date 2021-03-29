package client

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb"
	"github.com/zbindenren/sfmt"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
)

type serviceCmd struct {
	List       serviceList       `cmd:"" help:"List registered services."`
	Register   serviceRegister   `cmd:"" help:"Register a service."`
	UnRegister serviceUnRegister `cmd:"" help:"Unregister a service by ID or endpoint URL." name:"unregister"`
}

type serviceList struct {
	Output        string `short:"o" default:"table" help:"Output formats. Valid formats: json, yaml, csv, table."`
	SortBy        string `default:"endpoint" help:"Sort services by endpoint or modification date (allowed values: endpoint or date)" enum:"endpoint,date"`
	Headers       bool   `short:"H" help:"Show headers."`
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

	sw := sfmt.SliceWriter{
		Writer:    os.Stdout,
		NoHeaders: !s.Headers,
	}
	f := sfmt.ParseFormat(s.Output)

	sw.Write(f, services)

	return nil
}

type serviceRegister struct {
	Endpoints []string          `short:"e" help:"The service endpoint URLs." required:"true"`
	Name      string            `arg:"true" optional:"true" help:"The service name. This will represent the job name in prometheus." env:"DISCOVERY_NAME"`
	Labels    map[string]string `short:"l" help:"Labels for the service." mapsep:","`
	Namespace string            `short:"n" help:"The namespace for the service" default:"default" required:"true"`
	Selector  string            `short:"s" help:"Kubernetes style selectors (key=value) to select servers with specific labels."`
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

		r, err := cli.RegisterService(ctx, &discoveryv1.RegisterServiceRequest{
			Name:      s.Name,
			Namespace: s.Namespace,
			Endpoint:  ep,
			Labels:    s.Labels,
			Selector:  s.Selector,
		})
		if err != nil {
			lastErr = err
			l.Errorw("failed to unregister", "service", ep, "err", err)
		}

		l.Infow("service registered", "id", r.GetService().GetId())
	}

	return lastErr
}

type serviceUnRegister struct {
	Endpoints []string `arg:"true" optional:"true" help:"The service endpoint URLs or IDs." env:"DISCOVERY_ENDPOINTS"`
	Namespace string   `short:"n" help:"The namespace for the service" default:"default" required:"true"`
}

func (s serviceUnRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	var lastErr error

	for _, ep := range s.Endpoints {
		ctx, cancel := g.ctx()
		defer cancel()

		_, err = cli.UnRegisterService(ctx, &discoveryv1.UnRegisterServiceRequest{
			Namespace: s.Namespace,
			Id:        ep,
		})

		if err != nil {
			lastErr = err
			l.Errorw("failed to unregister", "service", ep, "err", err)
		}
	}

	return lastErr
}

type serviceFilter struct {
	Name     string `short:"N" help:"Filter services by job name (regular expression)."`
	Server   string `short:"S" help:"Filter services by server name (regular expression)."`
	Endpoint string `short:"e" help:"Filter services by endpoint (regular expression)."`
	Selector string `short:"s" help:"Filter services by selector."`
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
