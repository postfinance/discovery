// Package server represents the discovery server command.
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/server"
	discoveryv1connect "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1/discoveryv1connect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/zbindenren/king"
	"gitlab.pnet.ch/linux/go/auth"
	"gitlab.pnet.ch/linux/go/auth/oidc"
	"gitlab.pnet.ch/linux/go/auth/self"
	"gitlab.pnet.ch/linux/go/crpcauth"
	"go.uber.org/zap"
)

// CLI is the server command.
type CLI struct {
	Globals
	Server   serverCmd   `cmd:"" help:"Start discovery grpc server" default:"1"`
	Exporter exporterCmd `cmd:"" help:"Start exporter server"`
}

type serverCmd struct {
	ListenAddr  string    `help:"HTTP listen adddress" default:":3001"`
	Replicas    int       `help:"The number of service replicas." default:"1"`
	TokenIssuer string    `help:"The jwt token issuer name. If you change this, alle issued tokens are invalid." default:"discovery.postfinance.ch"`
	TokenSecret string    `help:"The secret key to issue jwt machine tokens. If you change this, alle issued tokens are invalid." required:"true"`
	OIDC        oidcFlags `embed:"true" prefix:"oidc-"`
	CACert      string    `help:"Path to a custom tls ca pem file. Certificates in this file are added to system cert pool." type:"existingfile"`
}

type oidcFlags struct {
	Endpoint      string   `help:"OIDC endpoint URL." required:"true"`
	ClientID      string   `help:"OIDC client ID." required:"true"`
	Roles         []string `help:"The the roles that are allowed to change servers and namespaces and to issue machine tokens." required:"true"`
	UsernameClaim string   `name:"username-claim" help:"The URL to the oidc server." default:"username"`
	RolesClaim    string   `name:"roles-claim" help:"The URL to the oidc server." default:"roles"`
}

//nolint:interfacer // kong does not work with interfaces
func (s serverCmd) Run(g *Globals, l *zap.SugaredLogger, app *kong.Context, registry *prometheus.Registry) error {
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	config, err := s.config(l, registry)
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

	srv, err := server.New(b, l, *config)
	if err != nil {
		return err
	}

	return srv.Run(ctx)
}

func (s serverCmd) config(l *zap.SugaredLogger, registry prometheus.Registerer) (*server.Config, error) {
	tokenHandler := self.NewTokenHandler(s.TokenIssuer, s.TokenSecret)

	cfg := rbacConfig()

	for _, c := range cfg {
		for _, r := range c.Rules {
			l.Debugw("rbac", "role", c.Role, "service", r.Service, "methods", r.Methods)
		}
	}

	v, err := oidc.NewVerifier(s.OIDC.Endpoint, s.OIDC.ClientID, 30*time.Second, &oidc.PfportalClaims{})
	if err != nil {
		return nil, fmt.Errorf("setup oidc token verifier (pfportal): %w", err)
	}

	a := crpcauth.NewAuthorizer(cfg,
		crpcauth.WithVerifier("discovery.postfinance.ch", tokenHandler),
		crpcauth.WithVerifierByIssuerAndClientID("https://p1-auth-oidc.pnet.ch:7048/auth/pfportal/openid", "cop", v),
		crpcauth.WithAuthCallback(func(ctx context.Context, u auth.User) {
			fmt.Println("-------", u.Name, u.Roles)
		}),
		crpcauth.WithPublicEndpoints(
			discoveryv1connect.NamespaceAPIListNamespaceProcedure,
			discoveryv1connect.ServerAPIListServerProcedure,
			discoveryv1connect.ServiceAPIListServiceProcedure,
			discoveryv1connect.TokenAPIInfoProcedure,
			discoveryv1connect.TokenAPICreateProcedure, // TODO: Remove
		),
	)

	return &server.Config{
		PrometheusRegistry: registry,
		NumReplicas:        s.Replicas,
		ListenAddr:         s.ListenAddr,
		TokenHandler:       tokenHandler,
		Interceptors:       []connect.Interceptor{a.UnaryServerInterceptor()},
	}, nil
}
