package auth

import (
	"context"

	"connectrpc.com/connect"
)

// WithToken configures a token authenticator for use in connect.WithInterceptors(...).
func WithToken(token string) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			req.Header().Set(MetadataHeader, MetadataSchema+" "+token)

			return next(ctx, req)
		})
	})
}
