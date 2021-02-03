// Package auth handles openid connect and jwt (for access tokens) authentication and
// authorization.
package auth

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Func creates a authentication function that can be used in combination with grpc middleware.
//
// The function verifies two kinds of tokens:
// First, it verfifies with the given key, if the token is a valid jwt token. If the token is valid access is granted
// for machines.
//
// If the above fails it checks if the token is a valid oidc token. If successful access is granted to a user.
//
// In both (successful) cases it extracts the user and adds it in the current context.
//
// Reflection and list requests are not authorized.
func Func(verifier Verifier, th *TokenHandler, l *zap.SugaredLogger) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		methodName := methodNameFromContext(ctx)
		if strings.HasPrefix(methodName, "/grpc.reflection") {
			return ctx, nil
		}

		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "failed to get authentication token: %s", err)
		}

		// machine tokens
		ok, err := th.IsMachine(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "failed to parse machine token:  %s", err)
		}

		if ok {
			u, err := th.Validate(token)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "machine token is invalid: %s", err)
			}

			l.Debugw("grpc authentication",
				"methodName", methodName,
				"isUserToken", u.IsUser(),
				"name", u.Username,
			)

			return context.WithValue(ctx, userKey, u), nil
		}

		// personal personal
		idToken, err := verifier.Verify(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "personal token is not valid: %s", err)
		}

		c := claims{}

		err = idToken.Claims(&c)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get claims: %s", err)
		}

		u := User{
			Username: c.Username,
			Roles:    c.Roles,
			Email:    c.Email,
			Kind:     UserToken,
		}

		l.Infow("grpc authentication",
			"methodName", methodName,
			"isUserToken", u.IsUser(),
			"name", u.Username,
			"roles", strings.Join(u.Roles, ","),
		)

		return context.WithValue(ctx, userKey, u), nil
	}
}

// UnaryAuthorizeInterceptor authorizes GRPC requests.
func UnaryAuthorizeInterceptor(rwRoles ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := authorizeUser(ctx, rwRoles...); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// UnaryMethodNameInterceptor adds GRPC method name to context.
func UnaryMethodNameInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		wrappedCtx := context.WithValue(ctx, methodNameKey, info.FullMethod)

		return handler(wrappedCtx, req)
	}
}

// StreamMethodNameInterceptor adds GRPC method name to context.
func StreamMethodNameInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = context.WithValue(stream.Context(), methodNameKey, info.FullMethod)

		return handler(srv, wrapped)
	}
}

// StreamAuthorizeInterceptor authorizes GRPC streams.
func StreamAuthorizeInterceptor(rwRoles ...string) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := authorizeUser(stream.Context(), rwRoles...); err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

//nolint: gocyclo // maybe move rules to other function
func authorizeUser(ctx context.Context, rwRoles ...string) error {
	fullMethod := methodNameFromContext(ctx)
	if strings.HasPrefix(fullMethod, "/grpc.reflection") {
		return nil
	}

	u, ok := UserFromContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "unauthententicated user")
	}

	// following methods are only allowed when use has one of rwRoles.
	switch fullMethod {
	case "/postfinance.discovery.v1.NamespaceAPI/RegisterNamespace":
		if u.IsMachine() || !u.HasRole(rwRoles...) {
			return status.Errorf(codes.PermissionDenied, "%s token for %s is not allowed to register a namespace", u.Kind.String(), u.Username)
		}
	case "/postfinance.discovery.v1.NamespaceAPI/UnregisterNamespace":
		if u.IsMachine() || !u.HasRole(rwRoles...) {
			return status.Errorf(codes.PermissionDenied, "%s token for %s is not allowed to unregister a namespace", u.Kind.String(), u.Username)
		}
	case "/postfinance.discovery.v1.ServerAPI/RegisterServer":
		if u.IsMachine() || !u.HasRole(rwRoles...) {
			return status.Errorf(codes.PermissionDenied, "%s token for %s is not allowed to register a server", u.Kind.String(), u.Username)
		}
	case "/postfinance.discovery.v1.ServerAPI/UnregisterServer":
		if u.IsMachine() || !u.HasRole(rwRoles...) {
			return status.Errorf(codes.PermissionDenied, "%s token for %s is not allowed to unregister a server", u.Kind.String(), u.Username)
		}
	case "/postfinance.discovery.v1.TokenAPI/Create":
		if u.IsMachine() || !u.HasRole(rwRoles...) {
			return status.Errorf(codes.PermissionDenied, "%s token for %s is not allowed to create a token", u.Kind.String(), u.Username)
		}
	}

	return nil
}

func methodNameFromContext(ctx context.Context) string {
	m, ok := ctx.Value(methodNameKey).(string)
	if !ok {
		return ""
	}

	return m
}

type key int

const (
	userKey key = iota
	methodNameKey
)
