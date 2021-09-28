package auth

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// TokenHandler creates tokens.
type TokenHandler struct {
	issuer string
	secret string
}

// NewTokenHandler creates a now TokenHandler
func NewTokenHandler(secret, issuer string) *TokenHandler {
	return &TokenHandler{
		issuer: issuer,
		secret: secret,
	}
}

// Create creates a new token. If expires is 0, it never expires.
func (t *TokenHandler) Create(id string, expires time.Duration, namespaces ...string) (string, error) {
	now := time.Now()

	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			Issuer:    t.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Namespaces: namespaces,
	}

	if expires > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(expires))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(t.secret))
}

// Validate validates a token. If successful it returns a machine user.
func (t *TokenHandler) Validate(token string) (*User, error) {
	tknClaims := TokenClaims{}

	tkn, err := jwt.ParseWithClaims(token, &tknClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})

	claims, ok := tkn.Claims.(*TokenClaims)
	if !ok || !tkn.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims.RegisteredClaims.Issuer != t.issuer {
		return nil, fmt.Errorf("wrong issuer is '%s', not %s", tknClaims.Issuer, t.issuer)
	}

	u := User{
		Username:   claims.RegisteredClaims.ID,
		Namespaces: claims.Namespaces,
		Kind:       MachineToken,
	}

	if claims.ExpiresAt != nil {
		u.ExpiresAt = claims.ExpiresAt.Time
	}

	return &u, nil
}

// IsMachine checks if token is a machine token issued by lslb service.
func (t *TokenHandler) IsMachine(token string) (bool, error) {
	tknClaims := jwt.MapClaims{}
	p := new(jwt.Parser)

	if _, _, err := p.ParseUnverified(token, &tknClaims); err != nil {
		return false, err
	}

	return tknClaims["iss"] == t.issuer, nil
}

// TokenClaims is like jwt standard claims with additional list of namespaces.
type TokenClaims struct {
	jwt.RegisteredClaims
	Namespaces []string `json:"namespaces,omitempty"`
}
