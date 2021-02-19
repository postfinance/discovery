package repo

import (
	"context"
	"testing"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamespaceChan(t *testing.T) {
	c, err := hash.New(hash.WithPrefix("/discovery"))
	require.NoError(t, err)

	r := NewNamespace(c)

	n := discovery.Namespace{
		Name: "namespace",
	}

	errHandler := func(err error) {
		require.NoError(t, err)
	}

	ch := r.Chan(context.Background(), errHandler)

	go func() {
		_, err := r.Save(n)
		assert.NoError(t, err)
		err = r.Delete("namespace")
		assert.NoError(t, err)
	}()

	e := <-ch
	assert.Equal(t, Change, e.Event)
	assert.Equal(t, "namespace", e.Name)

	e = <-ch
	assert.Equal(t, Delete, e.Event)
	assert.Equal(t, "namespace", e.Name)
}
