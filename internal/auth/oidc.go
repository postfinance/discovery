package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/errors"
)

// Verifier verifers an OIDC token.
type Verifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

// NewVerifier creates a new oidc verifier.
func NewVerifier(url, clientID string, timeout time.Duration, transport http.RoundTripper) (*oidc.IDTokenVerifier, error) {
	cli := &http.Client{
		Timeout: timeout,
	}

	if transport != nil {
		cli.Transport = transport
	}

	ctx := oidc.ClientContext(context.Background(), cli)

	provider, err := oidc.NewProvider(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get oidc provider")
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	return provider.Verifier(oidcConfig), nil
}

// AppendCertsToSystemPool adds certificates to system cert pool. If it is not possible to get system pool,
// certificates are added to an emptycert pool.
func AppendCertsToSystemPool(pemFile string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(pemFile) //nolint: gosec // we need to read that file
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", pemFile, err)
	}

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool()
	}

	caCertPool.AppendCertsFromPEM(caCert)

	return caCertPool, nil
}

// NewTLSTransportFromCertPool creates a new *http.Transport form cert pool.
func NewTLSTransportFromCertPool(pool *x509.CertPool) *http.Transport {
	tlsConfig := &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}
	tlsConfig.BuildNameToCertificate()

	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}
}

type claims struct {
	Roles      []string `json:"roles"`
	Name       string   `json:"name"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
}
