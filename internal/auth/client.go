package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

const (
	dfltTimeout = 5 * time.Second
)

// Token holds all necessary token info.
type Token struct {
	Username     string    `yaml:"-"`
	ISS          string    `yaml:"-"`
	AUD          string    `yaml:"-"`
	RefreshToken string    `yaml:"refresh_token"`
	IDToken      string    `yaml:"id_token"`
	AccessToken  string    `yaml:"access_token"`
	Expiry       time.Time `yaml:"-"`
}

// ClientOption is a functional option to configure the Client.
type ClientOption func(a *Client)

// Client handles requests to the keycloak server.
type Client struct {
	cli           *http.Client
	endPoint      string
	clientID      string
	tokenEndpoint string
}

// NewClient creates a new client with a configured token endpoint.
func NewClient(endPoint, clientID string, opts ...ClientOption) (*Client, error) {
	c := Client{
		cli: &http.Client{
			Timeout: dfltTimeout,
		},
		endPoint: endPoint,
		clientID: clientID,
	}

	for _, opt := range opts {
		opt(&c)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.endPoint, ".well-known/openid-configuration"), nil)
	if err != nil {
		return nil, fmt.Errorf("creating token endpoint request: %w", err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	d := json.NewDecoder(resp.Body)

	if err := d.Decode(&m); err != nil {
		return nil, fmt.Errorf("decoding body: %w", err)
	}

	c.tokenEndpoint = fmt.Sprintf("%s", m["token_endpoint"])

	return &c, nil
}

// WithTimeout overrides the default timeout of the httpclient.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.cli.Timeout = timeout
	}
}

// WithTransport overrides the default transport of the httpclient.
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.cli.Transport = transport
	}
}

// Token returns an OAUTH 2.0 token with Password Grant type.
func (c *Client) Token(username, password string, out io.Writer) (*Token, error) {
	t, err := c.getToken(username, password)
	if err != nil {
		return nil, err
	}

	return t, err
}

// Refresh refreshes a token if it expires in 10 seconds from now.
func (c *Client) Refresh(t *Token) (*Token, error) {
	if t == nil {
		return nil, nil
	}

	if t.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token provided")
	}

	if !t.isExpiredIn(10 * time.Second) {
		return t, nil
	}

	ts, err := c.tokenSource(*t)
	if err != nil {
		return nil, err
	}

	ot, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}

	newToken := Token{
		IDToken:      ot.AccessToken,
		RefreshToken: t.RefreshToken,
	}

	if err := newToken.parse(); err != nil {
		return t, err
	}

	return &newToken, err
}

// NewToken creates a new Token from idToken and refreshToken.
func NewToken(idToken, refreshToken string) (*Token, error) {
	if idToken == "" {
		return nil, fmt.Errorf("id token must not be empty")
	}

	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token must not be empty")
	}

	t := Token{
		IDToken:      idToken,
		RefreshToken: refreshToken,
	}
	if err := t.parse(); err != nil {
		return nil, err
	}

	return &t, nil
}

// tokenSource returns a oauth2 TokenSource that returns a
// token until it expires, automatically refreshing it as necessary
// using the provided context.
func (c *Client) tokenSource(t Token) (oauth2.TokenSource, error) {
	if t.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token found")
	}

	ctx := oidc.ClientContext(context.Background(), c.cli)

	provider, err := oidc.NewProvider(ctx, c.endPoint)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID: c.clientID,
		Endpoint: provider.Endpoint(),

		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	ts := oauth2Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: t.RefreshToken,
	})

	return ts, nil
}

// getToken uses the Password Grant Flow
func (c *Client) getToken(username, password string) (*Token, error) {
	data := url.Values{"scope": {"openid"}, "username": {username}, "password": {password}, "grant_type": {"password"}}

	return c.requestToken(data)
}

func (c *Client) requestToken(data url.Values) (*Token, error) {
	req, err := http.NewRequest(http.MethodPost, c.tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}

	req.SetBasicAuth(c.clientID, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting token: %w", err)
	}

	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body from token request: %w", err)
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshaling body: %w body: %s", err, string(body))
	}

	var t Token

	if val, ok := m["id_token"].(string); ok {
		t.IDToken = val
	}

	if val, ok := m["access_token"].(string); ok {
		t.AccessToken = val
	}

	if val, ok := m["refresh_token"].(string); ok {
		t.RefreshToken = val
	}

	if err := t.parse(); err != nil {
		return nil, err
	}

	return &t, nil
}

// isExpiredIn returns true if token will expire in d.
func (t Token) isExpiredIn(d time.Duration) bool {
	return t.Expiry.Before(time.Now().Add(d))
}

func (t *Token) parse() error {
	var (
		parser = new(jwt.Parser)
		claims = make(jwt.MapClaims)
		token  string
	)

	if t.IDToken == "" {
		token = t.AccessToken
	} else {
		token = t.IDToken
	}

	if _, _, err := parser.ParseUnverified(token, claims); err != nil {
		return fmt.Errorf("parsing jwt token: %w", err)
	}

	t.AUD = fmt.Sprintf("%s", claims["aud"])
	t.ISS = fmt.Sprintf("%s", claims["iss"])
	t.Username = fmt.Sprintf("%s", claims["username"])

	i, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("expiry is not of type float64: %v", claims["exp"])
	}

	t.Expiry = time.Unix(int64(i), 0)

	return nil
}

func handleResponse(r *http.Response) error {
	if r.StatusCode < 200 || r.StatusCode > 399 {
		body, _ := ioutil.ReadAll(r.Body)

		// reset the response body to the original unread state
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		var e errResponse
		if err := json.Unmarshal(body, &e); err != nil {
			return fmt.Errorf("request failed - status %s: %s", r.Status, string(body))
		}

		return fmt.Errorf("request %s: %s: %s", r.Status, e.Error, e.ErrorDescription)
	}

	return nil
}

type errResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
