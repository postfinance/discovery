package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/errors"
)

// Verifier verifers an OIDC token.
type Verifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

// NewVerifier creates a new oidc verifier.
func NewVerifier(url, clientID string, timeout time.Duration) (*oidc.IDTokenVerifier, error) {
	cli := &http.Client{
		Timeout: timeout,
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

type claims struct {
	Roles      []string `json:"roles"`
	Name       string   `json:"name"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
}
