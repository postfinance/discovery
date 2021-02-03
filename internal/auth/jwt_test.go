package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	th := NewTokenHandler("thesecret", "issuer")

	token, err := th.Create("username", 1*time.Hour, "namespace1", "namespace2")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	u, err := th.Validate(token)
	require.NoError(t, err)
	assert.Equal(t, u.Username, "username")
	assert.Equal(t, u.Namespaces, []string{"namespace1", "namespace2"})
	assert.Equal(t, u.Kind, MachineToken)
	assert.True(t, u.ExpiresAt.After(time.Now()))
}
