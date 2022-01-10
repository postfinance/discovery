// Package client represents the discovery client command.
package client

import (
	"context"
	"crypto/x509"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/postfinance/discovery/internal/auth"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	"github.com/zbindenren/king"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v3"
)

// CLI is the client command.
type CLI struct {
	Globals
	Server    serverCmd    `cmd:"" help:"Register and unregister servers." aliases:"srv"`
	Login     loginCmd     `cmd:"" help:"Perform OIDC login."`
	Service   serviceCmd   `cmd:"" help:"Register and unregister services." aliases:"svc"`
	Namespace namespaceCmd `cmd:"" help:"Register and unregister namespaces." aliases:"ns"`
	Token     tokenCmd     `cmd:"" help:"Manage access tokens"`
}

// Globals are the global client flags.
type Globals struct {
	Address      string           `short:"a" help:"The address of the discovery grpc endpoint." default:"localhost:3001"`
	Timeout      time.Duration    `help:"The request timeout" default:"5s"`
	Debug        bool             `short:"d" help:"Log debug output."`
	Insecure     bool             `help:"use insecure connection without tls." xor:"tls"`
	ShowConfig   king.ShowConfig  `help:"Show used config files"`
	Version      king.VersionFlag `help:"Show version information"`
	TokenPath    string           `help:"Authentication token" default:"~/.config/discovery/.token"`
	OIDCEndpoint string           `help:"OIDC endpoint URL." required:"true"`
	OIDCClientID string           `help:"OIDC client ID." required:"true"`
	CACert       string           `help:"Path to a custom tls ca pem file. Certificates in this file are added to system cert pool." type:"existingfile" xor:"tls"`
}

func (g Globals) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), g.Timeout)
}

func (g Globals) conn() (*grpc.ClientConn, error) {
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	if !g.Insecure {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		if g.CACert != "" {
			pool, err = auth.AppendCertsToSystemPool(g.CACert)
			if err != nil {
				return nil, err
			}
		}

		creds := credentials.NewClientTLSFromCert(pool, "")
		dialOption = grpc.WithTransportCredentials(creds)
	}

	token, err := g.getToken()
	if err != nil {
		return nil, err
	}

	dialOpts := []grpc.DialOption{dialOption, grpc.WithUnaryInterceptor(buildClientInterceptor(token))}

	return grpc.Dial(g.Address, dialOpts...)
}

func (g Globals) serverClient() (discoveryv1.ServerAPIClient, error) {
	conn, err := g.conn()
	if err != nil {
		return nil, err
	}

	return discoveryv1.NewServerAPIClient(conn), nil
}

func (g Globals) serviceClient() (discoveryv1.ServiceAPIClient, error) {
	conn, err := g.conn()
	if err != nil {
		return nil, err
	}

	return discoveryv1.NewServiceAPIClient(conn), nil
}

func (g Globals) namespaceClient() (discoveryv1.NamespaceAPIClient, error) {
	conn, err := g.conn()
	if err != nil {
		return nil, err
	}

	return discoveryv1.NewNamespaceAPIClient(conn), nil
}

func (g Globals) tokenClient() (discoveryv1.TokenAPIClient, error) {
	conn, err := g.conn()
	if err != nil {
		return nil, err
	}

	return discoveryv1.NewTokenAPIClient(conn), nil
}

func buildClientInterceptor(token string) func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, grpc.UnaryInvoker, ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type token struct {
	MachineToken string `yaml:"machine_token"` // for machines
	RefreshToken string `yaml:"refresh_token"`
	IDToken      string `yaml:"id_token"`
}

func (g Globals) loadToken() (*token, error) {
	path := kong.ExpandPath(g.TokenPath)

	d, err := os.ReadFile(path) //nolint: gosec // just reading token
	if err != nil {
		return nil, err
	}

	t := &token{}

	if err := yaml.Unmarshal(d, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (g Globals) saveToken(t *auth.Token) error {
	path := kong.ExpandPath(g.TokenPath)
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	d, err := yaml.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(path, d, 0o600)
}

// getToken loads or asks for user or machine token. If it is
// a user token it refreshes it if necessary.
func (g Globals) getToken() (string, error) {
	var token string

	t, err := g.loadToken()
	if err != nil {
		return "", errors.Wrap(err, "login required")
	}

	if err == nil && t.MachineToken == "" {
		cli, err := auth.NewClient(g.OIDCEndpoint, g.OIDCClientID)
		if err != nil {
			return "", err
		}

		tkn := &auth.Token{
			RefreshToken: t.RefreshToken,
			IDToken:      t.IDToken,
		}

		nt, err := cli.Refresh(tkn)
		if err != nil {
			return "", err
		}

		token = nt.IDToken

		if nt.IDToken != t.IDToken {
			if err := g.saveToken(nt); err != nil {
				return "", err
			}
		}
	}

	if t.MachineToken != "" {
		token = t.MachineToken
	}

	return token, nil
}
