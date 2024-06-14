package auth

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

// TokenKind represents the two different kind of tokens.
type TokenKind int

const (
	// SelfIssuedToken are tokens which are issued by this library.
	SelfIssuedToken TokenKind = iota
	// ExternalToken are tokens which are issued by an identity provider by OpenID connect.
	ExternalToken
)

// TokenClaims represents the information asserted about a subject.
type TokenClaims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}
