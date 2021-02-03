package auth

import (
	"context"
	"time"
)

// User is a oicd user.
type User struct {
	Username   string
	Email      string
	Roles      []string
	Namespaces []string
	ExpiresAt  time.Time
	Kind       TokenKind
}

// IsUser returns true if the token corresponds to a user token and
// false if it is a machine token.
func (u User) IsUser() bool {
	return u.Kind == UserToken
}

// IsMachine returns true if the token corresponds to a machine token and
// false if it is a user token.
func (u User) IsMachine() bool {
	return u.Kind == MachineToken
}

// UserFromContext gets user from context.
func UserFromContext(ctx context.Context) (User, bool) {
	userPtr, ok := ctx.Value(userKey).(*User)
	if ok {
		return *userPtr, true
	}

	user, ok := ctx.Value(userKey).(User)

	return user, ok
}

// HasRole returns true if user has one of roles.
func (u User) HasRole(roles ...string) bool {
	for _, ur := range u.Roles {
		for _, r := range roles {
			if r == ur {
				return true
			}
		}
	}

	return false
}

// HasNamespace returns true if user has one of namespaces.
func (u User) HasNamespace(namespaces ...string) bool {
	for _, un := range u.Namespaces {
		for _, u := range namespaces {
			if u == un {
				return true
			}
		}
	}

	return false
}

// TokenKind defines the kind of token. There are two possible tokens: machine and users.
type TokenKind int

// Two possible tokens: machine and users. User tokens are issued by oidc provider, where machine
// tokens are issued by discovery service.
const (
	MachineToken TokenKind = iota // machine
	UserToken                     // user
)
