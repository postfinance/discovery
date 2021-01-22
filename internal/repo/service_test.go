package repo

import (
	"testing"

	"github.com/postfinance/store/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zbindenren/discovery"
)

func TestService(t *testing.T) {
	c, err := hash.New()
	require.NoError(t, err)

	r := NewService(c)

	s, err := discovery.NewService("test", "http://example.com/metrics")
	require.NoError(t, err)

	var id string

	s.Labels = discovery.Labels{"env": "test"}

	t.Run("save", func(t *testing.T) {
		svc, err := r.Save(*s)
		assert.NoError(t, err)
		assert.True(t, s.Modified.Before(svc.Modified))
		assert.NotEmpty(t, svc.ID)
		id = svc.ID
	})

	t.Run("update", func(t *testing.T) {
		s, err = discovery.NewService("test", "https://www.example.com/metrics")
		require.NoError(t, err)
		s.Labels = discovery.Labels{"env": "test"}
		s.ID = id
		svc, err := r.Save(*s)
		assert.NoError(t, err)
		assert.True(t, s.Modified.Before(svc.Modified))
		assert.NotEmpty(t, svc.ID)
		id = svc.ID
	})

	t.Run("list", func(t *testing.T) {
		svcs, err := r.List("", "")
		assert.NoError(t, err)
		assert.Len(t, svcs, 1)
	})

	t.Run("list with selector", func(t *testing.T) {
		svcs, err := r.List("", "env=prod")
		assert.NoError(t, err)
		assert.Len(t, svcs, 0)

		svcs, err = r.List("", "env=test")
		assert.NoError(t, err)
		assert.Len(t, svcs, 1)
	})

	t.Run("list with namespace", func(t *testing.T) {
		svcs, err := r.List("other", "")
		assert.NoError(t, err)
		assert.Len(t, svcs, 0)

		svcs, err = r.List("default", "")
		assert.NoError(t, err)
		assert.Len(t, svcs, 1)
	})

	t.Run("delete", func(t *testing.T) {
		err := r.Delete(id, "default")
		assert.NoError(t, err)

		err = r.Delete("notexit", id)
		assert.Error(t, err)

		err = r.Delete("default", id)
		assert.Error(t, err)

		servers, err := r.List("", "")
		assert.NoError(t, err)
		assert.Len(t, servers, 0)
	})
}
