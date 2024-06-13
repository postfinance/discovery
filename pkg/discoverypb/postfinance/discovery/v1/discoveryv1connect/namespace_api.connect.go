// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: postfinance/discovery/v1/namespace_api.proto

package discoveryv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	v1 "postfinance/discovery/v1"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// NamespaceAPIName is the fully-qualified name of the NamespaceAPI service.
	NamespaceAPIName = "postfinance.discovery.v1.NamespaceAPI"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// NamespaceAPIRegisterNamespaceProcedure is the fully-qualified name of the NamespaceAPI's
	// RegisterNamespace RPC.
	NamespaceAPIRegisterNamespaceProcedure = "/postfinance.discovery.v1.NamespaceAPI/RegisterNamespace"
	// NamespaceAPIUnregisterNamespaceProcedure is the fully-qualified name of the NamespaceAPI's
	// UnregisterNamespace RPC.
	NamespaceAPIUnregisterNamespaceProcedure = "/postfinance.discovery.v1.NamespaceAPI/UnregisterNamespace"
	// NamespaceAPIListNamespaceProcedure is the fully-qualified name of the NamespaceAPI's
	// ListNamespace RPC.
	NamespaceAPIListNamespaceProcedure = "/postfinance.discovery.v1.NamespaceAPI/ListNamespace"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	namespaceAPIServiceDescriptor                   = v1.File_postfinance_discovery_v1_namespace_api_proto.Services().ByName("NamespaceAPI")
	namespaceAPIRegisterNamespaceMethodDescriptor   = namespaceAPIServiceDescriptor.Methods().ByName("RegisterNamespace")
	namespaceAPIUnregisterNamespaceMethodDescriptor = namespaceAPIServiceDescriptor.Methods().ByName("UnregisterNamespace")
	namespaceAPIListNamespaceMethodDescriptor       = namespaceAPIServiceDescriptor.Methods().ByName("ListNamespace")
)

// NamespaceAPIClient is a client for the postfinance.discovery.v1.NamespaceAPI service.
type NamespaceAPIClient interface {
	// RegisterNamespace registers a namespace.
	RegisterNamespace(context.Context, *connect.Request[v1.RegisterNamespaceRequest]) (*connect.Response[v1.RegisterNamespaceResponse], error)
	// UnRegisterNamespace unregisters a namespace.
	UnregisterNamespace(context.Context, *connect.Request[v1.UnregisterNamespaceRequest]) (*connect.Response[v1.UnregisterNamespaceResponse], error)
	// ListNamespace lists all namespaces.
	ListNamespace(context.Context, *connect.Request[v1.ListNamespaceRequest]) (*connect.Response[v1.ListNamespaceResponse], error)
}

// NewNamespaceAPIClient constructs a client for the postfinance.discovery.v1.NamespaceAPI service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewNamespaceAPIClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) NamespaceAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &namespaceAPIClient{
		registerNamespace: connect.NewClient[v1.RegisterNamespaceRequest, v1.RegisterNamespaceResponse](
			httpClient,
			baseURL+NamespaceAPIRegisterNamespaceProcedure,
			connect.WithSchema(namespaceAPIRegisterNamespaceMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		unregisterNamespace: connect.NewClient[v1.UnregisterNamespaceRequest, v1.UnregisterNamespaceResponse](
			httpClient,
			baseURL+NamespaceAPIUnregisterNamespaceProcedure,
			connect.WithSchema(namespaceAPIUnregisterNamespaceMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listNamespace: connect.NewClient[v1.ListNamespaceRequest, v1.ListNamespaceResponse](
			httpClient,
			baseURL+NamespaceAPIListNamespaceProcedure,
			connect.WithSchema(namespaceAPIListNamespaceMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// namespaceAPIClient implements NamespaceAPIClient.
type namespaceAPIClient struct {
	registerNamespace   *connect.Client[v1.RegisterNamespaceRequest, v1.RegisterNamespaceResponse]
	unregisterNamespace *connect.Client[v1.UnregisterNamespaceRequest, v1.UnregisterNamespaceResponse]
	listNamespace       *connect.Client[v1.ListNamespaceRequest, v1.ListNamespaceResponse]
}

// RegisterNamespace calls postfinance.discovery.v1.NamespaceAPI.RegisterNamespace.
func (c *namespaceAPIClient) RegisterNamespace(ctx context.Context, req *connect.Request[v1.RegisterNamespaceRequest]) (*connect.Response[v1.RegisterNamespaceResponse], error) {
	return c.registerNamespace.CallUnary(ctx, req)
}

// UnregisterNamespace calls postfinance.discovery.v1.NamespaceAPI.UnregisterNamespace.
func (c *namespaceAPIClient) UnregisterNamespace(ctx context.Context, req *connect.Request[v1.UnregisterNamespaceRequest]) (*connect.Response[v1.UnregisterNamespaceResponse], error) {
	return c.unregisterNamespace.CallUnary(ctx, req)
}

// ListNamespace calls postfinance.discovery.v1.NamespaceAPI.ListNamespace.
func (c *namespaceAPIClient) ListNamespace(ctx context.Context, req *connect.Request[v1.ListNamespaceRequest]) (*connect.Response[v1.ListNamespaceResponse], error) {
	return c.listNamespace.CallUnary(ctx, req)
}

// NamespaceAPIHandler is an implementation of the postfinance.discovery.v1.NamespaceAPI service.
type NamespaceAPIHandler interface {
	// RegisterNamespace registers a namespace.
	RegisterNamespace(context.Context, *connect.Request[v1.RegisterNamespaceRequest]) (*connect.Response[v1.RegisterNamespaceResponse], error)
	// UnRegisterNamespace unregisters a namespace.
	UnregisterNamespace(context.Context, *connect.Request[v1.UnregisterNamespaceRequest]) (*connect.Response[v1.UnregisterNamespaceResponse], error)
	// ListNamespace lists all namespaces.
	ListNamespace(context.Context, *connect.Request[v1.ListNamespaceRequest]) (*connect.Response[v1.ListNamespaceResponse], error)
}

// NewNamespaceAPIHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewNamespaceAPIHandler(svc NamespaceAPIHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	namespaceAPIRegisterNamespaceHandler := connect.NewUnaryHandler(
		NamespaceAPIRegisterNamespaceProcedure,
		svc.RegisterNamespace,
		connect.WithSchema(namespaceAPIRegisterNamespaceMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	namespaceAPIUnregisterNamespaceHandler := connect.NewUnaryHandler(
		NamespaceAPIUnregisterNamespaceProcedure,
		svc.UnregisterNamespace,
		connect.WithSchema(namespaceAPIUnregisterNamespaceMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	namespaceAPIListNamespaceHandler := connect.NewUnaryHandler(
		NamespaceAPIListNamespaceProcedure,
		svc.ListNamespace,
		connect.WithSchema(namespaceAPIListNamespaceMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/postfinance.discovery.v1.NamespaceAPI/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case NamespaceAPIRegisterNamespaceProcedure:
			namespaceAPIRegisterNamespaceHandler.ServeHTTP(w, r)
		case NamespaceAPIUnregisterNamespaceProcedure:
			namespaceAPIUnregisterNamespaceHandler.ServeHTTP(w, r)
		case NamespaceAPIListNamespaceProcedure:
			namespaceAPIListNamespaceHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedNamespaceAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedNamespaceAPIHandler struct{}

func (UnimplementedNamespaceAPIHandler) RegisterNamespace(context.Context, *connect.Request[v1.RegisterNamespaceRequest]) (*connect.Response[v1.RegisterNamespaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("postfinance.discovery.v1.NamespaceAPI.RegisterNamespace is not implemented"))
}

func (UnimplementedNamespaceAPIHandler) UnregisterNamespace(context.Context, *connect.Request[v1.UnregisterNamespaceRequest]) (*connect.Response[v1.UnregisterNamespaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("postfinance.discovery.v1.NamespaceAPI.UnregisterNamespace is not implemented"))
}

func (UnimplementedNamespaceAPIHandler) ListNamespace(context.Context, *connect.Request[v1.ListNamespaceRequest]) (*connect.Response[v1.ListNamespaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("postfinance.discovery.v1.NamespaceAPI.ListNamespace is not implemented"))
}
