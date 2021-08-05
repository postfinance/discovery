package auth

import (
	"fmt"
	"strconv"
	"strings"
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
		CompatibleStandardClaims: CompatibleStandardClaims{
			Id:        id,
			Issuer:    t.issuer,
			IssuedAt:  TokenTime(now.Unix()),
			NotBefore: TokenTime(now.Unix()),
		},
		Namespaces: namespaces,
	}

	if expires > 0 {
		claims.ExpiresAt = TokenTime(now.Add(expires).Unix())
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

	if claims.CompatibleStandardClaims.Issuer != t.issuer {
		return nil, fmt.Errorf("wrong issuer is '%s', not %s", tknClaims.Issuer, t.issuer)
	}

	u := User{
		Username:   claims.CompatibleStandardClaims.Id,
		Namespaces: claims.Namespaces,
		Kind:       MachineToken,
		ExpiresAt:  claims.ExpiresAt.Time(),
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
	CompatibleStandardClaims
	Namespaces []string `json:"namespaces,omitempty"`
}

// CompatibleStandardClaims is the same as jwt.StandardClaims but also
// allows float64 for durations (in order to be backward compatible to old library).
type CompatibleStandardClaims struct {
	Audience  string    `json:"aud,omitempty"`
	ExpiresAt TokenTime `json:"exp,omitempty"`
	Id        string    `json:"jti,omitempty"` //nolint: revive,stylecheck // we to use the same name as the library
	IssuedAt  TokenTime `json:"iat,omitempty"`
	Issuer    string    `json:"iss,omitempty"`
	NotBefore TokenTime `json:"nbf,omitempty"`
	Subject   string    `json:"sub,omitempty"`
}

// Valid validates standard claims.
func (c CompatibleStandardClaims) Valid() error {
	claims := jwt.StandardClaims{
		Audience:  c.Audience,
		ExpiresAt: int64(c.ExpiresAt),
		Id:        c.Id,
		IssuedAt:  int64(c.IssuedAt),
		Issuer:    c.Issuer,
		NotBefore: int64(c.NotBefore),
		Subject:   c.Subject,
	}

	return claims.Valid()
}

// TokenTime represent the token times. It is an int64 with custom marshal and unmarshl
// to also alow float64 as token times for backward compatibility reasons. Old jwt package
// created float64 token times.
type TokenTime int64

// Time converts TokenTime to time.Time.
func (t *TokenTime) Time() time.Time {
	return time.Unix(int64(*t), 0)
}

func (t *TokenTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t)), nil
}

func (t *TokenTime) UnmarshalJSON(b []byte) error {
	stripped := strings.ReplaceAll(string(b), `"`, "")

	f, err := strconv.ParseFloat(stripped, 64)
	if err != nil {
		return err
	}

	*t = TokenTime(f)

	return nil
}
