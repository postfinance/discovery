package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortNamespaces(t *testing.T) {
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	namespaces := Servers{
		Server{
			Name:     "namespace1",
			Modified: now,
		},
		Server{
			Name:     "a-namespace",
			Modified: oneMinuteAgo,
		},
	}

	t.Run("sort by name", func(t *testing.T) {
		namespaces.SortByName()
		assert.Equal(t, []string{"a-namespace", "namespace1"}, namespaces.Names())
	})

	t.Run("sort by modification date", func(t *testing.T) {
		namespaces.SortByDate()
		assert.Equal(t, []string{"namespace1", "a-namespace"}, namespaces.Names())
	})
}
