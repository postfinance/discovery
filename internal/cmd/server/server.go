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
	"github.com/prometheus/client_golang/prometheus/collectors"
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

//nolint: interfacer // kong does not work with interfaces
func (s serverCmd) Run(g *Globals, l *zap.SugaredLogger, app *kong.Context, registry *prometheus.Registry) error {
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	config, err := s.config(registry)
	if err != nil {
		return err
	}

	l.Infow("starting grpc server",
		king.FlagMap(app, regexp.MustCompile("key"), regexp.MustCompile("password"), regexp.MustCompile("secret")).
			Rm("help", "env-help", "version", "show-config", "etcd-ca", "etcd-cert").
			Register(app.Model.Name, registry).
			List()...)

	b, err := g.backend()
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv, err := server.New(b, l, config)
	if err != nil {
		return err
	}

	return srv.Run(ctx)
}

func (s serverCmd) config(registry prometheus.Registerer) (server.Config, error) {
	var transport http.RoundTripper

	if s.CACert != "" {
		pool, err := auth.AppendCertsToSystemPool(s.CACert)
		if err != nil {
			return server.Config{}, err
		}

		transport = auth.NewTLSTransportFromCertPool(pool)
	}

	return server.Config{
		PrometheusRegistry: registry,
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
