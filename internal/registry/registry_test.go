package registry

import (
	"testing"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store/hash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRegistry(t *testing.T) {
	c, err := hash.New(hash.WithPrefix("/disovery"))
	require.NoError(t, err)

	r, err := New(c, prometheus.NewRegistry(), zap.NewNop().Sugar(), 2)
	require.NoError(t, err)

	t.Run("register three servers", func(t *testing.T) {
		_, err := r.RegisterServer("server1", discovery.Labels{"env": "prod"})
		require.NoError(t, err)
		_, err = r.RegisterServer("server2", discovery.Labels{"env": "prod"})
		require.NoError(t, err)
		_, err = r.RegisterServer("server3", discovery.Labels{"env": "prod"})
		require.NoError(t, err)
		_, err = r.RegisterServer("server4", discovery.Labels{"env": "test"})
		require.NoError(t, err)
	})

	t.Run("register namespace", func(t *testing.T) {
		ns := discovery.DefaultNamespace()
		_, err := r.RegisterNamespace(*ns)
		require.NoError(t, err)
	})

	t.Run("list namespace", func(t *testing.T) {
		n, err := r.ListNamespaces()
		require.NoError(t, err)
		require.Len(t, n, 1)
	})

	t.Run("register service with selector", func(t *testing.T) {
		s := discovery.MustNewService("prod", "http://p.example.com/metrics")
		s.Selector = "env=prod"
		ns, err := r.RegisterService(*s)
		require.NoError(t, err)
		assert.Equal(t, []string{"server2", "server3"}, ns.Servers)
	})

	t.Run("register new server", func(t *testing.T) {
		_, err = r.RegisterServer("server5", discovery.Labels{"env": "prod"})
		require.NoError(t, err)
		l, err := r.serviceRepo.List("", "")
		require.NoError(t, err)
		require.Len(t, l, 1)
		assert.Equal(t, []string{"server2", "server3"}, l[0].Servers)
	})

	t.Run("unregister new server", func(t *testing.T) {
		err := r.UnRegisterServer("server5")
		require.NoError(t, err)
		l, err := r.serviceRepo.List("", "")
		require.NoError(t, err)
		require.Len(t, l, 1)
		assert.Equal(t, []string{"server2", "server3"}, l[0].Servers)
	})

	t.Run("register service with selector", func(t *testing.T) {
		s := discovery.MustNewService("test", "http://t.example.com/metrics")
		s.Selector = "env=test"
		ns, err := r.RegisterService(*s)
		require.NoError(t, err)
		assert.Equal(t, []string{"server4"}, ns.Servers)
	})

	t.Run("unregister service by endpoint", func(t *testing.T) {
		err := r.UnRegisterService("http://p.example.com/metrics", "")
		require.NoError(t, err)
		l, err := r.serviceRepo.List("", "")
		require.NoError(t, err)
		assert.Len(t, l, 1)
	})

	t.Run("unregister namespace ", func(t *testing.T) {
		err := r.UnRegisterNamespace(discovery.DefaultNamespace().Name)
		require.NoError(t, err)
		l, err := r.serviceRepo.List("", "")
		require.NoError(t, err)
		assert.Len(t, l, 0)
	})
}
