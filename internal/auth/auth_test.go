package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func TestFunc(t *testing.T) {
	id := "username"
	badTokenHandler := &TokenHandler{
		secret: "badsecret",
		issuer: "wrongissuer",
	}
	badToken, err := badTokenHandler.Create(id, 0)
	require.NoError(t, err)
	tokenHandler := NewTokenHandler("thesecret", "discovery.postifnance.ch")
	goodToken, err := tokenHandler.Create(id, 0)
	require.NoError(t, err)
	nokVerifier := mockVerifier{ok: false}
	okVerifier := mockVerifier{ok: true}

	t.Run("bad machine token, nok oidc verifier", func(t *testing.T) {
		f := Func(nokVerifier, tokenHandler, zap.New(nil).Sugar())
		m := metadata.MD{}
		m.Set("authorization", "bearer "+badToken)
		ctx := metadata.NewIncomingContext(context.Background(), m)
		_, err = f(ctx)
		assert.Error(t, err)
	})
	t.Run("valid machine token, nok oidc verifier", func(t *testing.T) {
		f := Func(nokVerifier, tokenHandler, zap.New(nil).Sugar())
		m := metadata.MD{}
		m.Set("authorization", "bearer "+goodToken)
		ctx := metadata.NewIncomingContext(context.Background(), m)
		c, err := f(ctx)
		require.NoError(t, err)
		u, ok := UserFromContext(c)
		require.True(t, ok)
		require.False(t, u.IsUser())
		require.Equal(t, id, u.Username)
	})
	t.Run("bad machine token, ok oidc verifier", func(t *testing.T) {
		f := Func(okVerifier, tokenHandler, zap.New(nil).Sugar())
		m := metadata.MD{}
		m.Set("authorization", "bearer "+badToken)
		ctx := metadata.NewIncomingContext(context.Background(), m)
		_, err = f(ctx)
		// claims are private in oidc id token, therefore exactly this error should come
		// if we could set claims here, no error should ocur.
		assert.Contains(t, err.Error(), "oidc: claims not set")
	})
}

type mockVerifier struct {
	ok bool
}

func (m mockVerifier) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	if m.ok {
		return &oidc.IDToken{}, nil
	}
	return nil, errors.New("invalid token - mock")
}
