package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
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

// TokenClaims is like jwt standard claims with additional list of namespaces.
type TokenClaims struct {
	jwt.StandardClaims
	Namespaces []string `json:"namespaces,omitempty"`
}

// Create creates a new token. If expires is 0, it never expires.
func (t *TokenHandler) Create(id string, expires time.Duration, namespaces ...string) (string, error) {
	now := time.Now()

	claims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ID:        id,
			Issuer:    t.issuer,
			IssuedAt:  jwt.At(now),
			NotBefore: jwt.At(now),
		},
		Namespaces: namespaces,
	}

	if expires > 0 {
		claims.ExpiresAt = jwt.At(now.Add(expires))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(t.secret))
}

// Validate validates a token. If successful it returns a machine user.
func (t *TokenHandler) Validate(token string) (*User, error) {
	tknClaims := TokenClaims{}

	tkn, err := jwt.ParseWithClaims(token, &tknClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("wrong signing method in token: %v", token.Method)
		}
		return []byte(t.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("invalid access token")
	}

	if tknClaims.Issuer != t.issuer {
		return nil, errors.Errorf("wrong issuer is '%s', not %s", tknClaims.Issuer, t.issuer)
	}

	u := User{
		Username:   tknClaims.ID,
		Namespaces: tknClaims.Namespaces,
		Kind:       MachineToken,
	}

	if tknClaims.ExpiresAt != nil {
		u.ExpiresAt = tknClaims.ExpiresAt.Time
	}

	return &u, nil
}

// IsMachine checks if token is a machine token issued by lslb service.
func (t *TokenHandler) IsMachine(token string) (bool, error) {
	tknClaims := jwt.StandardClaims{}
	p := new(jwt.Parser)

	if _, _, err := p.ParseUnverified(token, &tknClaims); err != nil {
		return false, err
	}

	return tknClaims.Issuer == t.issuer, nil
}
