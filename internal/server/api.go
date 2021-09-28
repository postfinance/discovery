package server

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/auth"
	"github.com/postfinance/discovery/internal/exporter"
	"github.com/postfinance/discovery/internal/registry"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// API implements the GRPC API.
type API struct {
	discoveryv1.UnsafeNamespaceAPIServer // requires you to implement all gRPC services
	discoveryv1.UnsafeServerAPIServer    // requires you to implement all gRPC services
	discoveryv1.UnsafeServiceAPIServer   // requires you to implement all gRPC services
	discoveryv1.UnsafeTokenAPIServer     // requires you to implement all gRPC services
	r                                    *registry.Registry
	tokenHandler                         *auth.TokenHandler
}

// RegisterServer registers a server.
func (a *API) RegisterServer(_ context.Context, req *discoveryv1.RegisterServerRequest) (*discoveryv1.RegisterServerResponse, error) {
	s, err := a.r.RegisterServer(req.GetName(), req.GetLabels())
	if err != nil {
		if registry.IsValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err)
		}

		return nil, status.Errorf(codes.Internal, "could not register server %s in store: %s", req.GetName(), err)
	}

	return &discoveryv1.RegisterServerResponse{
		Server: convert.ServerToPB(s),
	}, nil
}

// UnregisterServer unregisters a server.
func (a *API) UnregisterServer(_ context.Context, req *discoveryv1.UnregisterServerRequest) (*discoveryv1.UnregisterServerResponse, error) {
	if err := a.r.UnRegisterServer(req.GetName()); err != nil {
		return nil, status.Errorf(codes.Internal, "could not unregister server %s in store: %s", req.GetName(), err)
	}

	return &discoveryv1.UnregisterServerResponse{}, nil
}

// ListServer lists all servers.
func (a *API) ListServer(_ context.Context, _ *discoveryv1.ListServerRequest) (*discoveryv1.ListServerResponse, error) {
	s, err := a.r.ListServer("")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list server: %s", err)
	}

	return &discoveryv1.ListServerResponse{
		Servers: convert.ServersToPB(s),
	}, nil
}

// RegisterService registers a service.
func (a *API) RegisterService(ctx context.Context, req *discoveryv1.RegisterServiceRequest) (*discoveryv1.RegisterServiceResponse, error) {
	if err := verifyUser(ctx, req.GetNamespace()); err != nil {
		return nil, err
	}

	s, err := discovery.NewService(req.GetName(), req.GetEndpoint())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "service with endpoint %s is invalid: %s", req.GetEndpoint(), err)
	}

	s.Labels = req.GetLabels()
	s.Description = req.GetDescription()
	s.Selector = req.GetSelector()

	if req.Namespace != "" {
		s.Namespace = req.Namespace
	}

	svc, err := a.r.RegisterService(*s)
	if err != nil {
		if registry.IsServersNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "no server found for selector '%s'", req.GetSelector())
		}

		if registry.IsNamespaceNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "namespace '%s' not found", req.GetNamespace())
		}

		if registry.IsValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err)
		}

		return nil, status.Errorf(codes.Internal, "could not register service %s in store: %s", req.GetEndpoint(), err)
	}

	return &discoveryv1.RegisterServiceResponse{
		Service: convert.ServiceToPB(svc),
	}, nil
}

// UnRegisterService unregisters a service.
func (a *API) UnRegisterService(ctx context.Context, req *discoveryv1.UnRegisterServiceRequest) (*discoveryv1.UnRegisterServiceResponse, error) {
	if err := verifyUser(ctx, req.GetNamespace()); err != nil {
		return nil, err
	}

	if err := a.r.UnRegisterService(req.GetId(), req.GetNamespace()); err != nil {
		return nil, status.Errorf(codes.Internal, "could not unregister service %s in namespace %s: %s", req.GetId(), req.GetNamespace(), err)
	}

	return &discoveryv1.UnRegisterServiceResponse{}, nil
}

// ListService lists all services.
func (a *API) ListService(_ context.Context, req *discoveryv1.ListServiceRequest) (*discoveryv1.ListServiceResponse, error) {
	s, err := a.r.ListService(req.GetNamespace(), "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list services: %s", err)
	}

	return &discoveryv1.ListServiceResponse{
		Services: convert.ServicesToPB(s),
	}, nil
}

// ListTargetGroup converts services to prometheus target groups.
func (a *API) ListTargetGroup(ctx context.Context, in *discoveryv1.ListTargetGroupRequest) (*discoveryv1.ListTargetGroupResponse, error) {
	var config discovery.ExportConfig

	switch in.Config {
	case "standard":
		config = discovery.Standard
	case "blackbox":
		config = discovery.Blackbox
	case "":
		config = discovery.Standard
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid exporter config: '%s'", in.Config)
	}

	s, err := a.r.ListService(in.Namespace, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list services: %s", err)
	}

	serverFilter, err := regexp.Compile(fmt.Sprintf(`^%s$`, in.Server))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid regular expression: '%s'", in.Server)
	}

	s = s.Filter(discovery.ServiceByServer(serverFilter))

	t := make([]*discoveryv1.TargetGroup, 0, len(s))

	for i := range s {
		tg := exporter.NewTargetGroup(s[i], config)
		t = append(t, convert.TargetGroupToPB(&tg))
	}

	return &discoveryv1.ListTargetGroupResponse{
		Targetgroups: t,
	}, nil
}

// RegisterNamespace registers a server.
func (a *API) RegisterNamespace(_ context.Context, req *discoveryv1.RegisterNamespaceRequest) (*discoveryv1.RegisterNamespaceResponse, error) {
	n, err := a.r.RegisterNamespace(discovery.Namespace{
		Name:     req.Name,
		Export:   discovery.ExportConfig(req.Export),
		Modified: time.Now(),
	})

	if err != nil {
		if registry.IsValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err)
		}

		return nil, status.Errorf(codes.Internal, "could not register namespace %s in store: %s", req.GetName(), err)
	}

	return &discoveryv1.RegisterNamespaceResponse{
		Namespace: convert.NamespaceToPB(n),
	}, nil
}

// UnregisterNamespace unregisters a namespace.
func (a *API) UnregisterNamespace(_ context.Context, req *discoveryv1.UnregisterNamespaceRequest) (*discoveryv1.UnregisterNamespaceResponse, error) {
	if err := a.r.UnRegisterNamespace(req.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "could not unregister namespace %s in store: %s", req.GetName(), err)
	}

	return &discoveryv1.UnregisterNamespaceResponse{}, nil
}

// ListNamespace lists all namespaces.
func (a *API) ListNamespace(_ context.Context, _ *discoveryv1.ListNamespaceRequest) (*discoveryv1.ListNamespaceResponse, error) {
	namespaces, err := a.r.ListNamespaces()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list namespaces: %s", err)
	}

	return &discoveryv1.ListNamespaceResponse{
		Namespaces: convert.NamespacesToPB(namespaces),
	}, nil
}

// Create creates an access token.
func (a *API) Create(_ context.Context, in *discoveryv1.CreateRequest) (*discoveryv1.CreateResponse, error) {
	var expiry time.Duration

	if in.GetExpires() != "" {
		d, err := time.ParseDuration(in.GetExpires())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid expiry duration %s: %s", in.GetExpires(), err)
		}

		expiry = d
	}

	token, err := a.tokenHandler.Create(in.Id, expiry, in.Namespaces...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token: %s", err)
	}

	return &discoveryv1.CreateResponse{
		Token: token,
	}, nil
}

// Info gives token information.
func (a *API) Info(_ context.Context, in *discoveryv1.InfoRequest) (*discoveryv1.InfoResponse, error) {
	u, err := a.tokenHandler.Validate(in.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "token %s is not valid: %s", in.GetToken(), err)
	}

	return &discoveryv1.InfoResponse{
		Tokeninfo: &discoveryv1.TokenInfo{
			Id:         u.Username,
			Namespaces: u.Namespaces,
			ExpiresAt:  convert.TimeToPB(&u.ExpiresAt),
		},
	}, nil
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
