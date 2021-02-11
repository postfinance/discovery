package repo

import (
	"context"
	"testing"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerChan(t *testing.T) {
	c, err := hash.New(hash.WithPrefix("/discovery"))
	require.NoError(t, err)

	r := NewServer(c)

	s := discovery.NewServer("server1", discovery.Labels{"env": "test"})
	require.NoError(t, err)

	errHandler := func(err error) {
		require.NoError(t, err)
	}

	ch := r.Chan(context.Background(), errHandler)

	go func() {
		_, err := r.Save(*s)
		assert.NoError(t, err)
		err = r.Delete("server1")
		assert.NoError(t, err)
	}()

	e := <-ch
	assert.Equal(t, Change, e.Event)
	assert.Equal(t, "server1", e.Name)

	e = <-ch
	assert.Equal(t, Delete, e.Event)
	assert.Equal(t, "server1", e.Name)
}
