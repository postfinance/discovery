package repo

import (
	"testing"

	"github.com/postfinance/store/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zbindenren/discovery"
)

func TestServer(t *testing.T) {
	c, err := hash.New()
	require.NoError(t, err)

	r := NewServer(c)
	s := discovery.NewServer("server1", discovery.Labels{"env": "prod"})

	t.Run("save", func(t *testing.T) {
		server, err := r.Save(*s)
		assert.NoError(t, err)
		assert.True(t, s.Modified.Before(server.Modified))
	})

	t.Run("list", func(t *testing.T) {
		servers, err := r.List("")
		assert.NoError(t, err)
		assert.Len(t, servers, 1)
	})

	t.Run("list with selector", func(t *testing.T) {
		servers, err := r.List("env=test")
		assert.NoError(t, err)
		assert.Len(t, servers, 0)

		servers, err = r.List("env=prod")
		assert.NoError(t, err)
		assert.Len(t, servers, 1)
	})

	t.Run("delete", func(t *testing.T) {
		err := r.Delete("server1")
		assert.NoError(t, err)

		err = r.Delete("server1")
		assert.Error(t, err)

		servers, err := r.List("")
		assert.NoError(t, err)
		assert.Len(t, servers, 0)
	})
}
