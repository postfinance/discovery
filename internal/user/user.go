// Package user contains user logic.
package user

import (
	"gitlab.pnet.ch/linux/go/auth"
)

// // HasRole returns true if user has one of roles.
// func HasRole(u auth.User, roles ...string) bool {
// 	for _, ur := range u.Roles {
// 		for _, r := range roles {
// 			if r == ur {
// 				return true
// 			}
// 		}
// 	}
//
// 	return false
// }

// HasNamespace returns true if user has one of namespaces.
func HasNamespace(u auth.User, namespaces ...string) bool {
	userNamespaces := auth.MustGetData[[]string](u)
	if userNamespaces == nil {
		return false
	}

	for _, un := range *userNamespaces {
		for _, u := range namespaces {
			if u == un {
				return true
			}
		}
	}

	return false
}
