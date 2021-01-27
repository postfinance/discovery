package client

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/alecthomas/kong"
	"github.com/zbindenren/discovery"
	discoveryv1 "github.com/zbindenren/discovery/pkg/discoverypb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type importCmd struct {
	Path     string `arg:"true" help:"Path to yaml file." required:"true"`
	Selector string `short:"s" help:"Kubernetes style selectors (key=value) to select servers with specific labels."`
}

func (i importCmd) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serviceClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	d, err := ioutil.ReadFile(i.Path)
	if err != nil {
		return err
	}

	type service struct {
		discovery.Service `yaml:",inline"`
		Tags              discovery.Labels `yaml:"tags"`
	}

	services := []service{}

	err = yaml.Unmarshal(d, &services)
	if err != nil {
		return err
	}

	for j := range services {
		s := services[j]
		s.Selector = i.Selector
		l.Debugw("import serivce", s.KeyVals()...)

		_, err := cli.RegisterService(ctx, &discoveryv1.RegisterServiceRequest{
			Name:        s.Name,
			Endpoint:    s.Endpoint.String(),
			Description: s.Description,
			Labels:      s.Tags,
			Namespace:   s.Namespace,
			Selector:    s.Selector,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
