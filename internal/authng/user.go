package auth

import "context"

type contextKey int

// UserCtxKey represents the context key
const (
	UserCtxKey contextKey = iota
)

// User is an oidc user.
type User struct {
	Name  string    `json:"name"`
	Roles []string  `json:"roles"`
	Kind  TokenKind `json:"kind"`
}

// UserFromContext extracts the user information from the incoming context.
func UserFromContext(ctx context.Context) (User, bool) {
	userPtr, ok := ctx.Value(UserCtxKey).(*User)
	if ok {
		return *userPtr, true
	}

	user, ok := ctx.Value(UserCtxKey).(User)

	return user, ok
}
