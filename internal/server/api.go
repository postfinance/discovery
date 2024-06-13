package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/auth"
	"github.com/postfinance/discovery/internal/registry"
	"github.com/postfinance/discovery/internal/repo"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	discoveryv1connect "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1/discoveryv1connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ discoveryv1connect.ServerAPIHandler    = (*API)(nil)
	_ discoveryv1connect.NamespaceAPIHandler = (*API)(nil)
	_ discoveryv1connect.TokenAPIHandler     = (*API)(nil)
)

// API implements the GRPC API.
type API struct {
	r            *registry.Registry
	tokenHandler *auth.TokenHandler
}

// Create implements discoveryv1connect.TokenAPIHandler.
func (a *API) Create(_ context.Context, req *connect.Request[discoveryv1.CreateRequest]) (*connect.Response[discoveryv1.CreateResponse], error) {
	var expiry time.Duration

	if req.Msg.GetExpires() != "" {
		d, err := time.ParseDuration(req.Msg.GetExpires())
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid expiry duration %s: %w", req.Msg.GetExpires(), err))
		}

		expiry = d
	}

	token, err := a.tokenHandler.Create(req.Msg.GetId(), expiry, req.Msg.GetNamespaces()...)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create token: %w", err))
	}

	resp := connect.NewResponse[discoveryv1.CreateResponse](&discoveryv1.CreateResponse{
		Token: token,
	})

	return resp, nil
}

// Info implements discoveryv1connect.TokenAPIHandler.
func (a *API) Info(_ context.Context, in *connect.Request[discoveryv1.InfoRequest]) (*connect.Response[discoveryv1.InfoResponse], error) {
	u, err := a.tokenHandler.Validate(in.Msg.GetToken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("token %s is not valid: %w", in.Msg.GetToken(), err))
	}

	resp := connect.NewResponse[discoveryv1.InfoResponse](&discoveryv1.InfoResponse{
		Tokeninfo: &discoveryv1.TokenInfo{
			Id:         u.Username,
			Namespaces: u.Namespaces,
			ExpiresAt:  convert.TimeToPB(&u.ExpiresAt),
		},
	})

	return resp, nil
}

// ListServer implements discoveryv1connect.ServerAPIHandler.
func (a *API) ListServer(_ context.Context, req *connect.Request[discoveryv1.ListServerRequest]) (*connect.Response[discoveryv1.ListServerResponse], error) {
	s, err := a.r.ListServer("")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("could not list server: %w", err))
	}

	resp := connect.NewResponse[discoveryv1.ListServerResponse](&discoveryv1.ListServerResponse{
		Servers: convert.ServersToPB(s),
	})

	return resp, nil
}

// RegisterServer implements discoveryv1connect.ServerAPIHandler.
func (a *API) RegisterServer(_ context.Context, req *connect.Request[discoveryv1.RegisterServerRequest]) (*connect.Response[discoveryv1.RegisterServerResponse], error) {
	s, err := a.r.RegisterServer(req.Msg.GetName(), req.Msg.GetLabels())
	if err != nil {
		if registry.IsValidationError(err) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("could not register server %s in store: %w", req.Msg.GetName(), err))
	}

	resp := connect.NewResponse[discoveryv1.RegisterServerResponse](&discoveryv1.RegisterServerResponse{
		Server: convert.ServerToPB(s),
	})

	return resp, nil
}

// UnregisterServer implements discoveryv1connect.ServerAPIHandler.
func (a *API) UnregisterServer(_ context.Context, req *connect.Request[discoveryv1.UnregisterServerRequest]) (*connect.Response[discoveryv1.UnregisterServerResponse], error) {
	if err := a.r.UnRegisterServer(req.Msg.GetName()); err != nil {
		c := connect.CodeInternal
		if errors.Is(err, repo.ErrNotFound) {
			c = connect.CodeNotFound
		}

		return nil, connect.NewError(c, fmt.Errorf("could not unregister server %s in store: %w", req.Msg.GetName(), err))
	}

	resp := connect.NewResponse[discoveryv1.UnregisterServerResponse](&discoveryv1.UnregisterServerResponse{})

	return resp, nil
}

// ListNamespace implements discoveryv1connect.NamespaceAPIHandler.
func (a *API) ListNamespace(context.Context, *connect.Request[discoveryv1.ListNamespaceRequest]) (*connect.Response[discoveryv1.ListNamespaceResponse], error) {
	namespaces, err := a.r.ListNamespaces()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("could not list namespaces: %w", err))
	}

	resp := connect.NewResponse[discoveryv1.ListNamespaceResponse](&discoveryv1.ListNamespaceResponse{
		Namespaces: convert.NamespacesToPB(namespaces),
	})

	return resp, nil
}

// RegisterNamespace implements discoveryv1connect.NamespaceAPIHandler.
func (a *API) RegisterNamespace(_ context.Context, req *connect.Request[discoveryv1.RegisterNamespaceRequest]) (*connect.Response[discoveryv1.RegisterNamespaceResponse], error) {
	n, err := a.r.RegisterNamespace(discovery.Namespace{
		Name:     req.Msg.GetName(),
		Export:   discovery.ExportConfig(req.Msg.GetExport()),
		Modified: time.Now(),
	})
	if err != nil {
		if registry.IsValidationError(err) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("could not register namespace %s in store: %w", req.Msg.GetName(), err))
	}

	resp := connect.NewResponse[discoveryv1.RegisterNamespaceResponse](&discoveryv1.RegisterNamespaceResponse{
		Namespace: convert.NamespaceToPB(n),
	})

	return resp, nil
}

// UnregisterNamespace implements discoveryv1connect.NamespaceAPIHandler.
func (a *API) UnregisterNamespace(_ context.Context, req *connect.Request[discoveryv1.UnregisterNamespaceRequest]) (*connect.Response[discoveryv1.UnregisterNamespaceResponse], error) {
	if err := a.r.UnRegisterNamespace(req.Msg.GetName()); err != nil {
		c := connect.CodeInternal
		if errors.Is(err, repo.ErrNotFound) {
			c = connect.CodeNotFound
		}

		return nil, connect.NewError(c, fmt.Errorf("could not unregister namespace %s in store: %w", req.Msg.GetName(), err))
	}

	resp := connect.NewResponse[discoveryv1.UnregisterNamespaceResponse](&discoveryv1.UnregisterNamespaceResponse{})

	return resp, nil
}

func verifyUser(ctx context.Context, namespace string) error {
	u, ok := auth.UserFromContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "unauthententicated user")
	}

	if u.IsMachine() && !u.HasNamespace(namespace) {
		return status.Errorf(codes.PermissionDenied, "machine token %s (%s) is not allowed to change service in %s namespace", u.Username, strings.Join(u.Namespaces, ","), namespace)
	}

	return nil
}
