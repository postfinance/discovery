package server

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/exporter"
	"github.com/postfinance/single"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zbindenren/king"
	"go.uber.org/zap"
)

type exporterCmd struct {
	Directory      string        `help:"The destination directory." default:"/tmp/discovery"`
	Server         string        `help:"The server for which services should be exported." required:"true"`
	ResyncInterval time.Duration `help:"The interval in that the exporter resyncs all services to filesystem." default:"1h"`
}

func (e exporterCmd) Run(g *Globals, l *zap.SugaredLogger, app *kong.Context) error {
	reg := prometheus.NewRegistry()

	l.Infow("starting exporter",
		king.FlagMap(app, regexp.MustCompile("key"), regexp.MustCompile("password"), regexp.MustCompile("secret")).
			Rm("help", "env-help", "version").
			Register(app.Model.Name, reg).
			List()...)

	b, err := g.backend()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(e.Directory, 0700); err != nil {
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

	ctx := contextWithSignal(context.Background(), func(s os.Signal) {
		l.Infow("stopping server", "signal", s.String())
	}, syscall.SIGINT, syscall.SIGTERM)

	exp := exporter.New(b, l, e.Directory)

	return exp.Start(ctx, e.Server, e.ResyncInterval)
}
