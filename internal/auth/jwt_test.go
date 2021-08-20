package auth

import (
	"fmt"
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

	fmt.Println(token)

	t.Run("valid token", func(t *testing.T) {
		u, err := th.Validate(token)
		require.NoError(t, err)
		assert.Equal(t, u.Username, "username")
		assert.Equal(t, u.Namespaces, []string{"namespace1", "namespace2"})
		assert.Equal(t, u.Kind, MachineToken)
		assert.True(t, u.ExpiresAt.After(time.Now()))
	})

	t.Run("invalid token - wrong issuer", func(t *testing.T) {
		oth := NewTokenHandler("thesecret", "issuer2")
		u, err := oth.Validate(token)
		assert.Error(t, err)
		assert.Nil(t, u)
	})

	t.Run("invalid token - different secret", func(t *testing.T) {
		oth := NewTokenHandler("othersecret", "issuer")
		u, err := oth.Validate(token)
		assert.Error(t, err)
		assert.Nil(t, u)
	})
}

func TestJWTCompatibility(t *testing.T) {

	/*
		{
			"jti": "username",
			"iat": 1629461215.308863,
			"iss": "issuer",
			"nbf": 1629461215.308863,
			"namespaces": [
				"namespace1",
				"namespace2"
			]
		}
	*/
	const oldToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ1c2VybmFtZSIsImlhdCI6MTYyOTQ2MTIxNS4zMDg4NjMsImlzcyI6Imlzc3VlciIsIm5iZiI6MTYyOTQ2MTIxNS4zMDg4NjMsIm5hbWVzcGFjZXMiOlsibmFtZXNwYWNlMSIsIm5hbWVzcGFjZTIiXX0.suZSwZfDuLVdAGYy3rEJE0T3sSvq-qPi9SoOizMLkas`

	/*
		{
			"jti": "username",
			"iat": 1629461023,
			"iss": "issuer",
			"nbf": 1629461023,
			"namespaces": [
				"namespace1",
				"namespace2"
			]
		}
	*/
	const token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ1c2VybmFtZSIsImlhdCI6MTYyOTQ2MTAyMywiaXNzIjoiaXNzdWVyIiwibmJmIjoxNjI5NDYxMDIzLCJuYW1lc3BhY2VzIjpbIm5hbWVzcGFjZTEiLCJuYW1lc3BhY2UyIl19.MEQPTHAQNBQbn4pnqIJQctRgqnHcuJTsiHCiWmK_7ZE`

	th := NewTokenHandler("thesecret", "issuer")

	t.Run("valid old token", func(t *testing.T) {
		_, err := th.Validate(oldToken)
		assert.NoError(t, err)
	})

	t.Run("valid token", func(t *testing.T) {
		_, err := th.Validate(token)
		assert.NoError(t, err)
	})
}
