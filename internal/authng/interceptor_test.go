package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	jwt "github.com/golang-jwt/jwt/v4"
	auth "github.com/postfinance/discovery/internal/authng"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	"github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1/discoveryv1connect"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	testUser = "test"
)

var (
	client      *http.Client
	authzConfig = auth.Configs{
		auth.Config{
			Role: "echo",
			Rules: []auth.Rule{
				{
					Service: discoveryv1connect.ServerAPIName,
					Methods: []string{"ListServer", "RegisterServer", "UnregisterServer"},
				},
			},
		},
	}

	validToken                = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0IiwiaWF0IjoxNjE0NjgwOTkzLjg1NzI4MSwiaXNzIjoiZHVtbXkiLCJuYmYiOjE2MTQ2ODA5OTMuODU3MjgxLCJ1c2VyIjp7Im5hbWUiOiJ0ZXN0Iiwicm9sZXMiOlsiZWNobyJdLCJraW5kIjowfX0.Jca_DpqEwkBLSvRyJxFGd7zZcKdsNTs32nTb2TDnou0"                                  // iss = dummy (same as the verifier), no aud, does have role "echo"
	validTokenAudience        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0IiwiaWF0IjoxNjE0NjgwOTkzLjg1NzI4MSwiaXNzIjoiZHVtbXkiLCJuYmYiOjE2MTQ2ODA5OTMuODU3MjgxLCJhdWQiOiJ5dW1teSIsInVzZXIiOnsibmFtZSI6InRlc3QiLCJyb2xlcyI6WyJlY2hvIl0sImtpbmQiOjB9fQ.KJGgxWnZNkW_p6-2KPAMUdtKuFY8c18qbgtBa9bJDRc"               // iss = dummy (same as the verifier), aud = yummy, does have role "echo"
	unauthorizedToken         = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0IiwiaWF0IjoxNjE0NjgwOTkzLjg1NzI4MSwiaXNzIjoiZHVtbXkiLCJuYmYiOjE2MTQ2ODA5OTMuODU3MjgxLCJ1c2VyIjp7Im5hbWUiOiJ0ZXN0Iiwicm9sZXMiOlsicmVhZGVyIiwid3JpdGVyIl0sImtpbmQiOjB9fQ.1mMuHyEGPd44coov4iTx0ijNXuCDB0KEZ2FQEMt502g"                   // iss = dummy (same as the verifier), no aud, does not have role "echo"
	unauthorizedTokenAudience = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0IiwiaWF0IjoxNjE0NjgwOTkzLjg1NzI4MSwiaXNzIjoiZHVtbXkiLCJuYmYiOjE2MTQ2ODA5OTMuODU3MjgxLCJhdWQiOiJ5dW1teSIsInVzZXIiOnsibmFtZSI6InRlc3QiLCJyb2xlcyI6WyJyZWFkZXIiLCJ3cml0ZXIiXSwia2luZCI6MH19.D7jvClVX3XoLhi4eXx7wac_YZsCtTPx2YqB7suHVusY" // iss = dummy (same as the verifier), aud = yummy, does not have role "echo"
	invalidISSTok             = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0IiwiaWF0IjoxNjE0NjgwOTkzLjg1NzI4MSwiaXNzIjoiaW52YWxpZCIsIm5iZiI6MTYxNDY4MDk5My44NTcyODEsInVzZXIiOnsibmFtZSI6InRlc3QiLCJyb2xlcyI6WyJyZWFkZXIiLCJ3cml0ZXIiXSwia2luZCI6MH19.A_0wNVWthu6JlPkP0JMjlYwRivEpwO0JodCNOje92Uka"                // iss = invalid
)

var _ discoveryv1connect.ServerAPIHandler = (*api)(nil)

type api struct{}

// ListServer implements discoveryv1connect.ServerAPIHandler.
func (a *api) ListServer(context.Context, *connect.Request[discoveryv1.ListServerRequest]) (*connect.Response[discoveryv1.ListServerResponse], error) {
	resp := connect.NewResponse(&discoveryv1.ListServerResponse{
		Servers: []*discoveryv1.Server{
			{
				Name: "server",
			},
		},
	})

	return resp, nil
}

// RegisterServer implements discoveryv1connect.ServerAPIHandler.
func (a *api) RegisterServer(context.Context, *connect.Request[discoveryv1.RegisterServerRequest]) (*connect.Response[discoveryv1.RegisterServerResponse], error) {
	panic("unimplemented")
}

// UnregisterServer implements discoveryv1connect.ServerAPIHandler.
func (a *api) UnregisterServer(context.Context, *connect.Request[discoveryv1.UnregisterServerRequest]) (*connect.Response[discoveryv1.UnregisterServerResponse], error) {
	panic("unimplemented")
}

func testServer(a *auth.Authorizer, token string) (*httptest.Server, discoveryv1connect.ServerAPIClient) {
	mux := http.NewServeMux()

	tfPath, tfHandler := discoveryv1connect.NewServerAPIHandler(&api{}, connect.WithInterceptors(a.UnaryServerInterceptor()))
	mux.Handle(tfPath, tfHandler)

	ts := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))

	client := discoveryv1connect.NewServerAPIClient(ts.Client(), ts.URL, connect.WithInterceptors(auth.WithToken(token)))

	return ts, client
}

type dummyVerifier struct{}

var _ auth.Verifier = &dummyVerifier{}

func (d *dummyVerifier) Verify(ctx context.Context, token string) (*auth.User, error) {
	p := jwt.NewParser()

	var claims auth.TokenClaims

	if _, _, err := p.ParseUnverified(token, &claims); err != nil {
		return nil, err
	}

	return &claims.User, nil
}

func TestNewAuthorizer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("observed no panic")
		}
	}()

	// must panic, because no verifier is configured
	auth.NewAuthorizer(authzConfig)
}

func TestAuthInterceptor(t *testing.T) {
	t.Run("verifier by issuer without audience", func(t *testing.T) {
		a := auth.NewAuthorizer(authzConfig,
			auth.WithVerifier("dummy", &dummyVerifier{}),
			auth.WithVerifierByIssuerAndClientID("dummy", testUser, nil), // would cause an error if chosen
		)

		tsOK, clientOK := testServer(a, validToken)
		defer tsOK.Close()

		respOK, err := clientOK.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.NoError(t, err)
		require.Len(t, respOK.Msg.GetServers(), 1)

		tsNOK, clientNOK := testServer(a, unauthorizedToken)
		defer tsNOK.Close()

		respNOK, err := clientNOK.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.Error(t, err)
		require.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
		require.Nil(t, respNOK)
	})

	t.Run("verifier by issuer with audience", func(t *testing.T) {
		a := auth.NewAuthorizer(authzConfig,
			auth.WithVerifier("dummy", &dummyVerifier{}),
			auth.WithVerifierByIssuerAndClientID("dummy", testUser, nil), // would cause an error if chosen
		)

		tsOK, clientOK := testServer(a, validTokenAudience)
		defer tsOK.Close()

		respOK, err := clientOK.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.NoError(t, err)
		require.Len(t, respOK.Msg.GetServers(), 1)

		tsNOK, clientNOK := testServer(a, unauthorizedTokenAudience)
		defer tsNOK.Close()

		respNOK, err := clientNOK.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.Error(t, err)
		require.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
		require.Nil(t, respNOK)
	})

	t.Run("verifier by issuer and client id with audience", func(t *testing.T) {
		a := auth.NewAuthorizer(authzConfig,
			auth.WithVerifier("dummy", nil),                              // would cause an error if chosen
			auth.WithVerifierByIssuerAndClientID("dummy", testUser, nil), // would cause an error if chosen
			auth.WithVerifierByIssuerAndClientID("dummy", "yummy", &dummyVerifier{}),
		)

		ts, client := testServer(a, validTokenAudience)
		defer ts.Close()

		resp, err := client.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.NoError(t, err)
		require.Len(t, resp.Msg.GetServers(), 1)
	})

	t.Run("invalid token", func(t *testing.T) {
		a := auth.NewAuthorizer(authzConfig,
			auth.WithVerifier("dummy", &dummyVerifier{}),
		)

		ts, client := testServer(a, invalidISSTok)
		defer ts.Close()

		resp, err := client.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.Error(t, err)
		require.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
		require.Nil(t, resp)
	})

	t.Run("invalid token with public endpoint", func(t *testing.T) {
		a := auth.NewAuthorizer(authzConfig,
			auth.WithVerifier("dummy", &dummyVerifier{}),
			auth.WithPublicEndpoints(discoveryv1connect.ServerAPIListServerProcedure),
		)

		ts, client := testServer(a, invalidISSTok)
		defer ts.Close()

		resp, err := client.ListServer(context.TODO(), connect.NewRequest(&discoveryv1.ListServerRequest{}))

		require.NoError(t, err)
		require.Len(t, resp.Msg.GetServers(), 1)
	})
}
