package client

import (
	"errors"
	"os"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb"
	"github.com/zbindenren/sfmt"
	"go.uber.org/zap"
)

type namespaceCmd struct {
	List       namespaceList       `cmd:"" help:"List registered namespaces."`
	Register   namespaceRegister   `cmd:"" help:"Register a namespace."`
	UnRegister namespaceUnRegister `cmd:""  name:"unregister" help:"Unregister a namespace."`
}

type namespaceList struct {
	Output    string `short:"o" default:"table" help:"Output formats. Valid formats: json, yaml, csv, table."`
	NoHeaders bool   `short:"N" help:"Do not print headers."`
}

//nolint: dupl // it does not the same as serverList command
func (n namespaceList) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.namespaceClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.ListNamespace(ctx, &discoveryv1.ListNamespaceRequest{})
	if err != nil {
		return err
	}

	namespaces := convert.NamespacesFromPB(r.GetNamespaces())

	sw := sfmt.SliceWriter{
		Writer:    os.Stdout,
		NoHeaders: n.NoHeaders,
	}
	f := sfmt.ParseFormat(n.Output)

	sw.Write(f, namespaces)

	return nil
}

type namespaceRegister struct {
	Name         string `arg:"true" help:"Namespace name name." required:"true"`
	ExportConfig string `short:"e" help:"Configures how services get exported. Possible values: blackbox,standard and disabled." enum:"blackbox,standard,disabled" default:"standard"`
}

func (n namespaceRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.namespaceClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	var e discovery.ExportConfig

	switch n.ExportConfig {
	case "standard":
		e = discovery.Standard
	case "blackbox":
		e = discovery.Blackbox
	case "disabled":
		e = discovery.Disabled
	default:
		return errors.New("unsupported export configuration")
	}

	_, err = cli.RegisterNamespace(ctx, &discoveryv1.RegisterNamespaceRequest{
		Name:   n.Name,
		Export: int32(e),
	})

	return err
}

type namespaceUnRegister struct {
	Name string `arg:"true" help:"Server name." required:"true"`
}

func (n namespaceUnRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.namespaceClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	_, err = cli.UnregisterNamespace(ctx, &discoveryv1.UnregisterNamespaceRequest{
		Name: n.Name,
	})

	return err
}
