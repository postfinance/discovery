package client

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/zbindenren/discovery/internal/server/convert"
	discoveryv1 "github.com/zbindenren/discovery/pkg/discoverypb"
	"github.com/zbindenren/sfmt"
	"go.uber.org/zap"
)

type serverCmd struct {
	List       serverList       `cmd:"" help:"List registered servers."`
	Register   serverRegister   `cmd:"" help:"Register a server."`
	UnRegister serverUnRegister `cmd:""  name:"unregister" help:"Unregister a server."`
	Enable     serverEnable     `cmd:""  name:"enable" help:"Enable a server."`
	Disable    serverDisable    `cmd:""  name:"disable" help:"Disable a server."`
}

type serverList struct {
	Output    string `short:"o" default:"table" help:"Output formats. Valid formats: json, yaml, csv, table."`
	NoHeaders bool   `short:"N" help:"Do not print headers."`
}

//nolint: dupl // it does not the same as namespaceList command
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
		NoHeaders: s.NoHeaders,
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

	ctx, cancel := g.ctx()
	defer cancel()

	_, err = cli.RegisterServer(ctx, &discoveryv1.RegisterServerRequest{
		Name:   s.Name,
		Labels: s.Labels,
		Status: discoveryv1.RegisterServerRequest_Undefined,
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

	ctx, cancel := g.ctx()
	defer cancel()

	_, err = cli.UnregisterServer(ctx, &discoveryv1.UnregisterServerRequest{
		Name: s.Name,
	})

	return err
}

type serverEnable struct {
	Name string `arg:"true" help:"Server name." required:"true"`
}

func (s serverEnable) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	return setServerStatus(true, s.Name, g)
}

type serverDisable struct {
	Name string `arg:"true" help:"Server name." required:"true"`
}

func (s serverDisable) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	return setServerStatus(false, s.Name, g)
}

func setServerStatus(enabled bool, name string, g *Globals) error {
	cli, err := g.serverClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	status := discoveryv1.RegisterServerRequest_Disabled
	if enabled {
		status = discoveryv1.RegisterServerRequest_Enabled
	}

	_, err = cli.RegisterServer(ctx, &discoveryv1.RegisterServerRequest{
		Name:   name,
		Status: status,
	})

	return err
}
