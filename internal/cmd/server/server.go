// Package server represents the discovery server command.
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/auth"
	"github.com/postfinance/discovery/internal/server"
	"github.com/prometheus/client_golang/prometheus"
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
	GRPCListen   string   `short:"l" help:"GRPC gateway listen adddress" default:"localhost:3001"`
	HTTPListen   string   `help:"HTTP listen adddress" default:"localhost:3002"`
	Replicas     int      `help:"The number of service replicas." default:"1"`
	TokenIssuer  string   `help:"The jwt token issuer name. If you change this, alle issued tokens are invalid." default:"discovery.postfinance.ch"`
	TokenSecret  string   `help:"The secret key to issue jwt machine tokens. If you change this, alle issued tokens are invalid." required:"true"`
	OIDCEndpoint string   `help:"OIDC endpoint URL." required:"true"`
	OIDCClientID string   `help:"OIDC client ID." required:"true"`
	OIDCRoles    []string `help:"The the roles that are allowed to change servers and namespaces and to issue machine tokens." required:"true"`
	CACert       string   `help:"Path to a custom tls ca pem file. Certificates in this file are added to system cert pool." type:"existingfile"`
}

func (s serverCmd) Run(g *Globals, l *zap.SugaredLogger, app *kong.Context) error {
	config, err := s.config()
	if err != nil {
		return err
	}

	l.Infow("starting grpc server",
		king.FlagMap(app, regexp.MustCompile("key"), regexp.MustCompile("password"), regexp.MustCompile("secret")).
			Rm("help", "env-help", "version").
			Register(app.Model.Name, config.PrometheusRegistry).
			List()...)

	b, err := g.backend()
	if err != nil {
		return err
	}

	ctx := contextWithSignal(context.Background(), func(s os.Signal) {
		l.Infow("stopping server", "signal", s.String())
	}, syscall.SIGINT, syscall.SIGTERM)

	config.PrometheusRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	config.PrometheusRegistry.MustRegister(prometheus.NewGoCollector())

	srv, err := server.New(b, l, config)
	if err != nil {
		return err
	}

	return srv.Run(ctx)
}

func (s serverCmd) config() (server.Config, error) {
	var transport http.RoundTripper

	if s.CACert != "" {
		pool, err := auth.AppendCertsToSystemPool(s.CACert)
		if err != nil {
			return server.Config{}, err
		}

		transport = auth.NewTLSTransportFromCertPool(pool)
	}

	return server.Config{
		PrometheusRegistry: prometheus.NewRegistry(),
		NumReplicas:        s.Replicas,
		GRPCListenAddr:     s.GRPCListen,
		HTTPListenAddr:     s.HTTPListen,
		TokenIssuer:        s.TokenIssuer,
		TokenSecretKey:     s.TokenSecret,
		OIDCClient:         s.OIDCClientID,
		OIDCRoles:          s.OIDCRoles,
		OIDCURL:            s.OIDCEndpoint,
		Transport:          transport,
	}, nil
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
