// Package client represents the discovery client command.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"connectrpc.com/connect"
	"github.com/alecthomas/kong"
	"github.com/hashicorp/go-cleanhttp"
	discoveryv1connect "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1/discoveryv1connect"
	"github.com/zbindenren/king"
	"gitlab.pnet.ch/linux/go/crpcauth"
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
	Address    string           `short:"a" help:"The address of the discovery grpc endpoint." default:"http://localhost:3001"`
	Timeout    time.Duration    `help:"The request timeout" default:"15s"`
	Debug      bool             `short:"d" help:"Log debug output."`
	Insecure   bool             `help:"use insecure connection without tls." xor:"tls"`
	ShowConfig king.ShowConfig  `help:"Show used config files"`
	Version    king.VersionFlag `help:"Show version information"`
	TokenPath  string           `help:"Authentication token" default:"~/.config/discovery/.token"`
	OIDC       oidc             `embed:"true" prefix:"oidc-"`
	CACert     string           `help:"Path to a custom tls ca pem file. Certificates in this file are added to system cert pool." type:"existingfile" xor:"tls"`
}

type oidc struct {
	Endpoint         string `help:"OIDC endpoint URL."`
	ClientID         string `help:"OIDC client ID."`
	ExternalLoginCmd string `help:"If not empty, this command is printed out for the login sub command. The command should create a id_token in token-path." required:"true"`
}

func (g Globals) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), g.Timeout)
}

func (g Globals) serverClient() (discoveryv1connect.ServerAPIClient, error) {
	c, err := g.httpClient()
	if err != nil {
		return nil, err
	}

	token, err := g.loadToken()
	if err != nil {
		return nil, err
	}

	return discoveryv1connect.NewServerAPIClient(c, g.Address, connect.WithInterceptors(crpcauth.WithToken(token))), nil
}

func (g Globals) serviceClient() (discoveryv1connect.ServiceAPIClient, error) {
	c, err := g.httpClient()
	if err != nil {
		return nil, err
	}

	token, err := g.loadToken()
	if err != nil {
		return nil, err
	}

	return discoveryv1connect.NewServiceAPIClient(c, g.Address, connect.WithInterceptors(crpcauth.WithToken(token))), nil
}

func (g Globals) namespaceClient() (discoveryv1connect.NamespaceAPIClient, error) {
	c, err := g.httpClient()
	if err != nil {
		return nil, err
	}

	token, err := g.loadToken()
	if err != nil {
		return nil, err
	}

	return discoveryv1connect.NewNamespaceAPIClient(c, g.Address, connect.WithInterceptors(crpcauth.WithToken(token))), nil
}

func (g Globals) tokenClient() (discoveryv1connect.TokenAPIClient, error) {
	c, err := g.httpClient()
	if err != nil {
		return nil, err
	}

	return discoveryv1connect.NewTokenAPIClient(c, g.Address), nil
}

type token struct {
	MachineToken string `yaml:"machine_token"` // for machines
	RefreshToken string `yaml:"refresh_token"`
	IDToken      string `yaml:"id_token"`
}

func (g Globals) loadToken() (string, error) {
	path := kong.ExpandPath(g.TokenPath)

	d, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	t := &token{}

	if err := yaml.Unmarshal(d, t); err != nil {
		return "", err
	}

	if t.MachineToken == "" && t.IDToken == "" {
		return "", fmt.Errorf("no machine_token or id_token found in %s", filepath.Clean(path))
	}

	if t.MachineToken != "" {
		return t.MachineToken, nil
	}

	return t.IDToken, nil
}

func (g Globals) httpClient() (*http.Client, error) {
	clnt := cleanhttp.DefaultClient()
	clnt.Timeout = g.Timeout

	tlsConfig := &tls.Config{
		InsecureSkipVerify: g.Insecure, //nolint:gosec // configurable
	}

	if !g.Insecure {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		tlsConfig.RootCAs = pool
	}

	clnt.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return clnt, nil
}
