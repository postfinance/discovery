package client

import (
	"errors"
	"os"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb"
	"github.com/zbindenren/sfmt"
	"go.uber.org/zap"
)

type serviceCmd struct {
	List       serviceList       `cmd:"" help:"List registered services."`
	Register   serviceRegister   `cmd:"" help:"Register a service."`
	UnRegister serviceUnRegister `cmd:"" help:"Unregister a service by ID or endpoint URL." name:"unregister"`
}

type serviceList struct {
	Output    string `short:"o" default:"table" help:"Output formats. Valid formats: json, yaml, csv, table."`
	NoHeaders bool   `short:"N" help:"Do not print headers."`
	Namespace string `short:"n" help:"If not empty only serices for a namespace are listed."`
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

	sw := sfmt.SliceWriter{
		Writer:    os.Stdout,
		NoHeaders: s.NoHeaders,
	}
	f := sfmt.ParseFormat(s.Output)

	sw.Write(f, services)

	return nil
}

type serviceRegister struct {
	Endpoint  string            `short:"e" help:"The endpoint URL." required:"true"`
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
		c.PrintUsage(true)
		return errors.New("name cannot be empty")
	}

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.RegisterService(ctx, &discoveryv1.RegisterServiceRequest{
		Name:      s.Name,
		Namespace: s.Namespace,
		Endpoint:  s.Endpoint,
		Labels:    s.Labels,
		Selector:  s.Selector,
	})
	if err != nil {
		return err
	}

	l.Infow("service registered", "id", r.GetService().GetId())

	return nil
}

type serviceUnRegister struct {
	Services  []string `arg:"true" help:"The service endpoint URL or ID." required:"true"`
	Namespace string   `short:"n" help:"The namespace for the service" default:"default" required:"true"`
}

func (s serviceUnRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	var lastErr error

	for _, svc := range s.Services {
		ctx, cancel := g.ctx()
		defer cancel()

		_, err = cli.UnRegisterService(ctx, &discoveryv1.UnRegisterServiceRequest{
			Namespace: s.Namespace,
			Id:        svc,
		})

		if err != nil {
			lastErr = err
			l.Errorw("failed to unregister", "service", svc, "err", err)
		}
	}

	return lastErr
}
