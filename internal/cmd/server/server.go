// Package server represents the discovery server command.
package server

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zbindenren/discovery/internal/server"
	"github.com/zbindenren/king"
	"go.uber.org/zap"
)

// CLI is the server command.
type CLI struct {
	Globals
	Server   serverCmd   `cmd:"" help:"Start discovery grpc server" default:"1"`
	Exporter exporterCmd `cmd:"" help:"Start exporter server"`
}

type serverCmd struct {
	GRPCListen string `short:"l" help:"GRPC gateway listen adddress" default:"localhost:3001"`
	HTTPListen string `help:"HTTP listen adddress" default:"localhost:3002"`
	Replicas   int    `help:"The number of service replicas." default:"1"`
}

func (s serverCmd) Run(g *Globals, l *zap.SugaredLogger, app *kong.Context) error {
	reg := prometheus.NewRegistry()

	l.Infow("starting grpc server",
		king.FlagMap(app, regexp.MustCompile("key"), regexp.MustCompile("password"), regexp.MustCompile("secret")).
			Rm("help", "env-help", "version").
			Register(app.Model.Name, reg).
			List()...)

	b, err := g.backend()
	if err != nil {
		return err
	}

	ctx := contextWithSignal(context.Background(), func(s os.Signal) {
		l.Infow("stopping server", "signal", s.String())
	}, syscall.SIGINT, syscall.SIGTERM)

	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	reg.MustRegister(prometheus.NewGoCollector())

	srv, err := server.New(b, l, reg, s.Replicas, s.GRPCListen, s.HTTPListen)
	if err != nil {
		return err
	}

	return srv.Run(ctx)
}

func contextWithSignal(ctx context.Context, f func(s os.Signal), s ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, s...)

		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case sig := <-c:
			if f != nil {
				f(sig)
			}

			cancel()
		}
	}()

	return ctx
}
