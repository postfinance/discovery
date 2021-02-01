package repo

import (
	"testing"
	"time"

	"github.com/postfinance/store/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/postfinance/discovery"
)

func TestNamespace(t *testing.T) {
	c, err := hash.New()
	require.NoError(t, err)

	r := NewNamespace(c)
	n := discovery.Namespace{
		Name:     "namespace1",
		Export:   discovery.Standard,
		Modified: time.Now(),
	}

	t.Run("save", func(t *testing.T) {
		ns, err := r.Save(n)
		assert.NoError(t, err)
		assert.True(t, n.Modified.Before(ns.Modified))
	})

	t.Run("list", func(t *testing.T) {
		namespaces, err := r.List()
		assert.NoError(t, err)
		assert.Len(t, namespaces, 1)
	})

	t.Run("delete", func(t *testing.T) {
		err := r.Delete("namespace1")
		assert.NoError(t, err)

		err = r.Delete("namespace1")
		assert.Error(t, err)

		namespaces, err := r.List()
		assert.NoError(t, err)
		assert.Len(t, namespaces, 0)
	})
}
