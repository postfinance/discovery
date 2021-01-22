package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/labels"
)

func TestFilter(t *testing.T) {
	servers := Servers{
		*NewServer("server1", Labels{"env": "prod", "location": "ch"}),
		*NewServer("server2", Labels{"env": "test"}),
	}

	var tt = []struct {
		name     string
		selector string
		expected []string
	}{
		{
			"empty selector",
			"",
			[]string{"server1", "server2"},
		},
		{
			"matching selector",
			"env=prod",
			[]string{"server1"},
		},
		{
			"not matching selector",
			"env=dev",
			[]string{},
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			sel, err := labels.Parse(tc.selector)
			require.NoError(t, err)
			s := servers.Filter(ServersBySelector(sel))
			assert.Equal(t, s.Names(), tc.expected)
		})
	}
}

func TestSort(t *testing.T) {
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	servers := Servers{
		Server{
			Name:     "server1",
			Modified: now,
		},
		Server{
			Name:     "a-server",
			Modified: oneMinuteAgo,
		},
	}

	t.Run("sort by name", func(t *testing.T) {
		servers.SortByName()
		assert.Equal(t, []string{"a-server", "server1"}, servers.Names())
	})

	t.Run("sort by modification date", func(t *testing.T) {
		servers.SortByDate()
		assert.Equal(t, []string{"server1", "a-server"}, servers.Names())
	})
}
