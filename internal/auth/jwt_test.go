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
			"exp": 1628181960.502017,
			"jti":"username",
			"iat": 1628178360.502017,
			"iss": "issuer",
			"nbf": 1628178360.502017,
			"namespaces": [
				"name space1",
				"namespace2"
			]
		}
	*/
	const oldToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjgxODE5NjAuNTAyMDE3LCJqdGkiOiJ1c2VybmFtZSIsImlhdCI6MTYyODE3ODM2MC41MDIwMTcsImlzcyI6Imlzc3VlciIsIm5iZiI6MTYyODE3ODM2MC41MDIwMTcsIm5hbWVzcGFjZXMiOlsibmFtZXNwYWNlMSIsIm5hbWVzcGFjZTIiXX0.ILa9euUbxTsVMBcJwCrhyTdbfV6-PO0c9jB9SFO3KD0`

	/*
		{
			"exp": 1628182190,
			"jti": "username",
			"iat": 1628178590,
			"iss": "issuer",
			"nbf": 1628178590,
			"namespaces": [
				"namespace1",
				"namespace2"
			]
		}
	*/
	const token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjgxODIxOTAsImp0aSI6InVzZXJuYW1lIiwiaWF0IjoxNjI4MTc4NTkwLCJpc3MiOiJpc3N1ZXIiLCJuYmYiOjE2MjgxNzg1OTAsIm5hbWVzcGFjZXMiOlsibmFtZXNwYWNlMSIsIm5hbWVzcGFjZTIiXX0.bOiL4Bz4uQcF4t2FDqI081jrpO9AS3RyRNAhN63mnxY`

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
