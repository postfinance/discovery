package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/exporter"
	"github.com/postfinance/single"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/zbindenren/king"
	"go.uber.org/zap"
)

type exporterCmd struct {
	Directory      string        `help:"The destination directory." default:"/tmp/discovery"`
	Server         string        `help:"The server for which services should be exported." required:"true"`
	ResyncInterval time.Duration `help:"The interval in that the exporter resyncs all services to filesystem." default:"1h"`
	HTTPListen     string        `help:"HTTP listen adddress" default:"localhost:3003"`
}

//nolint:interfacer // kong does not work with interfaces
func (e exporterCmd) Run(g *Globals, l *zap.SugaredLogger, app *kong.Context, registry *prometheus.Registry) error {
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	l.Infow("starting exporter",
		king.FlagMap(app, regexp.MustCompile("key"), regexp.MustCompile("password"), regexp.MustCompile("secret")).
			Rm("help", "env-help", "version", "show-config", "etcd-ca", "etcd-cert").
			Register(app.Model.Name, registry).
			List()...)

	b, err := g.backend()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(e.Directory, 0o700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", e.Directory, err)
	}

	one, err := single.New(app.Model.Name, single.WithLockPath("/var/tmp"))
	if err != nil {
		l.Fatal(err)
	}

	if err := one.Lock(); err != nil {
		l.Fatal(err)
	}

	defer func() {
		if err := one.Unlock(); err != nil {
			l.Error(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := e.config(registry)

	exp := exporter.New(b, l, cfg)

	return exp.Start(ctx, e.Server)
}

func (e exporterCmd) config(registry prometheus.Registerer) exporter.Config {
	return exporter.Config{
		Directory:          e.Directory,
		ResyncInterval:     e.ResyncInterval,
		PrometheusRegistry: registry,
		HTTPListenAddr:     e.HTTPListen,
	}
}
