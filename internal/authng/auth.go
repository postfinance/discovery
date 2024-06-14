// Package auth provides helpers to authenticate and authorize RPC calls in a Kubernetes RBAC like approach.
// JWT tokens are used to transport identity information. Tokens can either be issued by methods
// provided by this package or originate from an openid connect provider.
package auth

import (
	"context"
	"slices"
)

const (
	// MetadataHeader is the name of the header
	MetadataHeader = "authorization"
	// MetadataSchema is the authorization schema
	MetadataSchema = "Bearer"
)

// Verifier verifies a JWT token.
type Verifier interface {
	Verify(ctx context.Context, rawIDToken string) (*User, error)
}

// Config describes which roles are allowed (authorized) to access certain services with the given methods.
// Example config:
// ---
//   - role: reader
//     rules:
//   - service: postfinance.burger.namespace.v1.NamespaceAPI
//     methods:
//   - Get
//   - Read
//   - service: postfinance.burger.namespace.v1.DeploymentAPI
//     methods:
//   - Get
//   - Read
type Config struct {
	Role  string `yaml:"role"`
	Rules []Rule `yaml:"rules"`
}

// Rule is an authorization rule matching the given api group and method(s).
type Rule struct {
	Service string   `yaml:"service"`
	Methods []string `yaml:"methods"`
}

// Configs is a slice of authorization configurations.
type Configs []Config

// IsAuthorized returns true if the user has the permission to access the service
func (ac Configs) IsAuthorized(service, method string, user User) bool {
	for idx := range ac {
		authz := ac[idx]

		if slices.Contains(user.Roles, authz.Role) {
			for _, rule := range authz.Rules {
				if rule.Service == service {
					if slices.Contains(rule.Methods, method) {
						return true
					}
				}
			}
		}
	}

	return false
}
