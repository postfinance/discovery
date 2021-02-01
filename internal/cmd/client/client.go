// Package client represents the discovery client command.
package client

import (
	"context"
	"crypto/x509"
	"time"

	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb"
	"github.com/zbindenren/king"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// CLI is the client command.
type CLI struct {
	Globals
	Server    serverCmd    `cmd:"" help:"Register and unregister servers."`
	Service   serviceCmd   `cmd:"" help:"Register and unregister services."`
	Namespace namespaceCmd `cmd:"" help:"Register and unregister namespaces."`
	Import    importCmd    `cmd:"" help:"Import new services"`
}

// Globals are the global client flags.
type Globals struct {
	Address    string           `short:"a" help:"The address of the discovery grpc endpoint." default:"localhost:3001"`
	Timeout    time.Duration    `help:"The request timeout" default:"5s"`
	Debug      bool             `short:"d" help:"Log debug output."`
	Insecure   bool             `help:"use insecure connection without tls."`
	ShowConfig king.ShowConfig  `help:"Show used config files"`
	Version    king.VersionFlag `help:"Show version information"`
	// Token      string           `help:"Authentication token"`
}

func (g Globals) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), g.Timeout)
}

func (g Globals) conn() (*grpc.ClientConn, error) {
	dialOption := grpc.WithInsecure()

	if !g.Insecure {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		creds := credentials.NewClientTLSFromCert(pool, "")
		dialOption = grpc.WithTransportCredentials(creds)
	}

	dialOpts := []grpc.DialOption{dialOption}

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
