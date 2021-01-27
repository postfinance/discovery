package client

import (
	"io/ioutil"

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

	ctx, cancel := g.ctx()
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

	for i := range services {
		s := services[i]
		l.Debugw("import serivce", s.KeyVals()...)

		_, err := cli.RegisterService(ctx, &discoveryv1.RegisterServiceRequest{
			Name:        s.Name,
			Endpoint:    s.Endpoint.String(),
			Description: s.Description,
			Labels:      s.Tags,
			Namespace:   s.Namespace,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
