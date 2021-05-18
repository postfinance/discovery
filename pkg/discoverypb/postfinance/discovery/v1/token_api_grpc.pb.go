// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package discoveryv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TokenAPIClient is the client API for TokenAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenAPIClient interface {
	// Create creates a token.
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	// Info gives token information.
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
}

type tokenAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenAPIClient(cc grpc.ClientConnInterface) TokenAPIClient {
	return &tokenAPIClient{cc}
}

func (c *tokenAPIClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/postfinance.discovery.v1.TokenAPI/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenAPIClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, "/postfinance.discovery.v1.TokenAPI/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenAPIServer is the server API for TokenAPI service.
// All implementations must embed UnimplementedTokenAPIServer
// for forward compatibility
type TokenAPIServer interface {
	// Create creates a token.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Info gives token information.
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
	mustEmbedUnimplementedTokenAPIServer()
}

// UnimplementedTokenAPIServer must be embedded to have forward compatible implementations.
type UnimplementedTokenAPIServer struct {
}

func (UnimplementedTokenAPIServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTokenAPIServer) Info(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedTokenAPIServer) mustEmbedUnimplementedTokenAPIServer() {}

// UnsafeTokenAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenAPIServer will
// result in compilation errors.
type UnsafeTokenAPIServer interface {
	mustEmbedUnimplementedTokenAPIServer()
}

func RegisterTokenAPIServer(s grpc.ServiceRegistrar, srv TokenAPIServer) {
	s.RegisterService(&TokenAPI_ServiceDesc, srv)
}

func _TokenAPI_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenAPIServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/postfinance.discovery.v1.TokenAPI/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenAPIServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenAPI_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenAPIServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/postfinance.discovery.v1.TokenAPI/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenAPIServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenAPI_ServiceDesc is the grpc.ServiceDesc for TokenAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "postfinance.discovery.v1.TokenAPI",
	HandlerType: (*TokenAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _TokenAPI_Create_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _TokenAPI_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "postfinance/discovery/v1/token_api.proto",
}