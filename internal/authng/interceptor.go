package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"

	"connectrpc.com/connect"
)

// Authorizer can be used to authenticate and authorize ConnectRPC unary calls through the use of the provided interceptors.
type Authorizer struct {
	verifiers    map[string]Verifier
	parser       *jwt.Parser
	public       map[string]bool
	config       Configs
	authCallback AuthCallback
}

// AuthCallback is a callback function which will be executed after successful authentication.
type AuthCallback func(context.Context, User)

// NewAuthorizer returns a configures authorizer which can be used as an interceptor to authenticate and authorize ConnectRPC
// streaming and unary calls.
func NewAuthorizer(c Configs, opts ...func(*Authorizer)) *Authorizer {
	a := Authorizer{
		verifiers: make(map[string]Verifier),
		parser:    new(jwt.Parser),
		config:    c,
	}

	for _, opt := range opts {
		opt(&a)
	}

	if len(a.verifiers) < 1 {
		panic("no token verifier(s) configured, use WithVerifier() option to configure one")
	}

	return &a
}

// WithPublicEndpoints configures public endpoints. The endpoint must be fully qualified, e.g.: /postfinance.echo.v1.EchoAPI/Echo
func WithPublicEndpoints(eps ...string) func(*Authorizer) {
	return func(a *Authorizer) {
		public := make(map[string]bool)
		for _, ep := range eps {
			public[ep] = true
		}

		a.public = public
	}
}

// WithAuthCallback configures a callback function in the authorizer.
func WithAuthCallback(cb AuthCallback) func(*Authorizer) {
	return func(a *Authorizer) {
		a.authCallback = cb
	}
}

// WithVerifier configures a token verifier for the given issuer. Can be provided multiple times with
// different issuers to validate tokens from different sources (eg. self issued and oidc issued tokens).
func WithVerifier(issuer string, verifier Verifier) func(*Authorizer) {
	return func(a *Authorizer) {
		a.verifiers[issuer] = verifier
	}
}

// WithVerifierByIssuerAndClientID configures a token verifier for the given issuer with different ClientIDs
// to validate tokens from the same sources with different ClientIDs.
func WithVerifierByIssuerAndClientID(issuer, clientID string, verifier Verifier) func(*Authorizer) {
	return func(a *Authorizer) {
		a.verifiers[fmt.Sprintf("%s::%s", issuer, clientID)] = verifier
	}
}

// UnaryServerInterceptor returns a ConnectRPC server interceptor to authenticate and authorize unary calls.
func (a *Authorizer) UnaryServerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// public endpoint - needs no authentication or authorization
			if a.public[req.Spec().Procedure] {
				return next(ctx, req)
			}

			// wrap authorization header in context to be compatible with grpcauth 1.2.1
			key := MetadataHeader
			value := req.Header().Get(MetadataHeader)
			vCtx := context.WithValue(context.Background(), key, value) //nolint:staticcheck // keep the logic from grpcauth 1.2.1

			wrappedCtx, err := a.authenticate(vCtx)
			if err != nil {
				return nil, err
			}

			if err := a.authorize(wrappedCtx, req.Spec().Procedure); err != nil {
				return nil, err
			}

			return next(wrappedCtx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}

// Authenticate authenticates a user. The jwt token is taken out of the incoming context and the issuer (ISS) is parsed
// out of the token to determine which token verifier to call. If a verifier for the issuer is found, it will be called
// to verify the token and obtain a user object (if the token is valid). The user is then placed in the outgoing context
// and can safely be used later.
func (a *Authorizer) authenticate(ctx context.Context) (context.Context, error) {
	val := ctx.Value(MetadataHeader).(string)
	if val == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing authorization header"))
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("malformed authorization string"))
	}

	scheme := splits[0]
	token := splits[1]

	if !strings.EqualFold(scheme, MetadataSchema) {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("authentication scheme %s is not supported", scheme))
	}

	var claims jwt.RegisteredClaims
	// Unverified parse, since we're only interested in the issuer of the token, so we can determine which verifier we
	// must use to parse and verify the token correctly.
	if _, _, err := a.parser.ParseUnverified(token, &claims); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("parse jwt: %w", err))
	}

	verifier := a.verifier(&claims)
	if verifier == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("unknown issuer %s and issuer with audience %s::[%s]", claims.Issuer, claims.Issuer, strings.Join(claims.Audience, ",")))
	}

	user, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("authentication failed: %w", err))
	}

	if a.authCallback != nil {
		a.authCallback(ctx, *user)
	}

	// warp the incoming context and put the user object into the new context
	return context.WithValue(ctx, UserCtxKey, user), nil
}

func (a *Authorizer) verifier(claims *jwt.RegisteredClaims) Verifier {
	if claims.Audience == nil {
		return a.verifiers[claims.Issuer]
	}

	for _, audience := range claims.Audience {
		for _, id := range []string{fmt.Sprintf("%s::%s", claims.Issuer, audience), claims.Issuer} {
			verifier, ok := a.verifiers[id]
			if ok {
				return verifier
			}
		}
	}

	return nil
}

// Authorize authorizes a ConnectRPC call. The ConnectRPC method is splited into its group and method, the user information is
// extracted from the incoming context. With this information authorization is performed, based on the roles of a user
// and the interceptors authorization configurations.
func (a *Authorizer) authorize(ctx context.Context, fullMethod string) error {
	// fullMethod in ConnectRPC is in the form /postfinance.burger.v1.NamespaceAPI/Create --> /service/method
	splits := strings.SplitN(fullMethod, "/", 3)
	if len(splits) != 3 {
		return connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("malformed ConnectRPC method %s", fullMethod))
	}

	method := splits[2]
	service := splits[1]

	user, ok := UserFromContext(ctx)
	if !ok {
		return connect.NewError(connect.CodeFailedPrecondition, errors.New("no user information found in metadata"))
	}

	if !a.config.IsAuthorized(service, method, user) {
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user %s with roles %s is not allowed to call %s in service %s",
			user.Name,
			strings.Join(user.Roles, ","),
			method,
			service,
		))
	}

	return nil
}
