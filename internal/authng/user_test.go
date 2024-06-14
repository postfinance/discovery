package auth_test

import (
	"context"
	"testing"

	auth "github.com/postfinance/discovery/internal/authng"
	"github.com/stretchr/testify/require"
)

func TestUserFromContext(t *testing.T) {
	tests := map[string]struct {
		ctx    context.Context
		expect bool
	}{
		"not found": {
			ctx:    context.WithValue(context.TODO(), auth.UserCtxKey, nil),
			expect: false,
		},
		"found": {
			ctx:    context.WithValue(context.TODO(), auth.UserCtxKey, auth.User{Name: "test"}),
			expect: true,
		},
		"found ptr": {
			ctx:    context.WithValue(context.TODO(), auth.UserCtxKey, &auth.User{Name: "test"}),
			expect: true,
		},
	}

	for name, tc := range tests {
		t.Run("UserFromContext "+name, func(t *testing.T) {
			r := require.New(t)

			user, found := auth.UserFromContext(tc.ctx)
			r.Equal(tc.expect, found)
			if tc.expect {
				r.NotEmpty(user)
			} else {
				r.Empty(user)
			}
		})
	}
}
