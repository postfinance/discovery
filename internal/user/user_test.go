package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.pnet.ch/linux/go/auth"
)

// func TestHasRoles(t *testing.T) {
// 	tt := []struct {
// 		roles    []string
// 		u        auth.User
// 		expected bool
// 	}{
// 		{
// 			[]string{"a", "b"},
// 			auth.User{},
// 			false,
// 		},
// 		{
// 			[]string{"a", "b"},
// 			auth.User{
// 				Roles: []string{"a"},
// 			},
// 			true,
// 		},
// 		{
// 			[]string{"a", "b"},
// 			auth.User{
// 				Roles: []string{"c"},
// 			},
// 			false,
// 		},
// 		{
// 			[]string{},
// 			auth.User{
// 				Roles: []string{"c"},
// 			},
// 			false,
// 		},
// 		{
// 			[]string{},
// 			auth.User{},
// 			false,
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		t.Run("", func(t *testing.T) {
// 			assert.Equal(t, tc.expected, HasRole(tc.u, tc.roles...))
// 		})
// 	}
// }

func TestHasNamespace(t *testing.T) {
	tt := []struct {
		namespaces []string
		u          auth.User
		expected   bool
	}{
		{
			[]string{"a", "b"},
			auth.User{},
			false,
		},
		{
			[]string{"a", "b"},
			auth.User{
				Data: []string{"a"},
			},
			true,
		},
		{
			[]string{"a", "b"},
			auth.User{
				Data: []string{"c"},
			},
			false,
		},
		{
			[]string{},
			auth.User{
				Data: []string{"c"},
			},
			false,
		},
		{
			[]string{},
			auth.User{},
			false,
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tc.expected, HasNamespace(tc.u, tc.namespaces...))
		})
	}
}
