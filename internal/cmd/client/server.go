package client

import (
	"context"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	"github.com/zbindenren/sfmt"
	"go.uber.org/zap"
)

const registerTimout = 300 * time.Second

type serverCmd struct {
	List       serverList       `cmd:"" help:"List registered servers."`
	Register   serverRegister   `cmd:"" help:"Register a server."`
	UnRegister serverUnRegister `cmd:""  name:"unregister" help:"Unregister a server."`
}

type serverList struct {
	Output  string `short:"o" default:"table" help:"Output formats. Valid formats: json, yaml, csv, table."`
	Headers bool   `short:"H" help:"Show headers."`
}

//nolint:dupl // it does not the same as namespaceList command
func (s serverList) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serverClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.ListServer(ctx, &discoveryv1.ListServerRequest{})
	if err != nil {
		return err
	}

	servers := convert.ServersFromPB(r.GetServers())

	sw := sfmt.SliceWriter{
		Writer:    os.Stdout,
		NoHeaders: !s.Headers,
	}
	f := sfmt.ParseFormat(s.Output)

	sw.Write(f, servers)

	return nil
}

type serverRegister struct {
	Name   string            `arg:"true" help:"Server name." required:"true"`
	Labels map[string]string `short:"l" help:"Labels" mapsep:","`
}

func (s serverRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serverClient()
	if err != nil {
		return err
	}

	l.Infow("increasing command timeout", "timeout", registerTimout)

	ctx, cancel := context.WithTimeout(context.Background(), registerTimout)
	defer cancel()

	_, err = cli.RegisterServer(ctx, &discoveryv1.RegisterServerRequest{
		Name:   s.Name,
		Labels: s.Labels,
	})

	return err
}

type serverUnRegister struct {
	Name string `arg:"true" help:"Server name." required:"true"`
}

func (s serverUnRegister) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.serverClient()
	if err != nil {
		return err
	}

	l.Infow("increasing command timeout", "timeout", registerTimout)

	ctx, cancel := context.WithTimeout(context.Background(), registerTimout)
	defer cancel()

	_, err = cli.UnregisterServer(ctx, &discoveryv1.UnregisterServerRequest{
		Name: s.Name,
	})

	return err
}
