package repo

import (
	"context"
	"testing"

	"github.com/postfinance/store/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zbindenren/discovery"
)

func TestChan(t *testing.T) {
	c, err := hash.New()
	require.NoError(t, err)

	r := NewService(c)

	s, err := discovery.NewService("test", "http://example.com/metrics")
	require.NoError(t, err)

	s.Labels = discovery.Labels{"env": "test"}

	errHandler := func(err error) {
		require.NoError(t, err)
	}

	ch := r.Chan(context.Background(), errHandler)

	go func() {
		svc, err := r.Save(*s)
		assert.NoError(t, err)
		err = r.Delete(svc.ID, "default")
		assert.NoError(t, err)
	}()

	e := <-ch
	assert.Equal(t, Change, e.Event)
	assert.NotEmpty(t, e.ID)

	e = <-ch
	assert.Equal(t, Delete, e.Event)
	assert.NotEmpty(t, e.ID)
}
