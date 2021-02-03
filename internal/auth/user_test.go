package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasRoles(t *testing.T) {
	var tt = []struct {
		roles    []string
		u        User
		expected bool
	}{
		{
			[]string{"a", "b"},
			User{},
			false,
		},
		{
			[]string{"a", "b"},
			User{
				Roles: []string{"a"},
			},
			true,
		},
		{
			[]string{"a", "b"},
			User{
				Roles: []string{"c"},
			},
			false,
		},
		{
			[]string{},
			User{
				Roles: []string{"c"},
			},
			false,
		},
		{
			[]string{},
			User{},
			false,
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.u.HasRole(tc.roles...))
		})
	}
}

func TestHasNamespace(t *testing.T) {
	var tt = []struct {
		namespaces []string
		u          User
		expected   bool
	}{
		{
			[]string{"a", "b"},
			User{},
			false,
		},
		{
			[]string{"a", "b"},
			User{
				Namespaces: []string{"a"},
			},
			true,
		},
		{
			[]string{"a", "b"},
			User{
				Namespaces: []string{"c"},
			},
			false,
		},
		{
			[]string{},
			User{
				Namespaces: []string{"c"},
			},
			false,
		},
		{
			[]string{},
			User{},
			false,
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.u.HasNamespace(tc.namespaces...))
		})
	}
}
