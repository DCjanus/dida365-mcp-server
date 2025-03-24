// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/oauth.proto

package api

import (
	context "context"
	model "github.com/dcjanus/dida365-mcp-server/gen/model"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Dida365OAuthService_Ping_FullMethodName          = "/api.Dida365oAuthService/Ping"
	Dida365OAuthService_OAuthLogin_FullMethodName    = "/api.Dida365oAuthService/OAuthLogin"
	Dida365OAuthService_OAuthCallback_FullMethodName = "/api.Dida365oAuthService/OAuthCallback"
	Dida365OAuthService_OAuthPrompt_FullMethodName   = "/api.Dida365oAuthService/OAuthPrompt"
)

// Dida365OAuthServiceClient is the client API for Dida365OAuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Dida365 OAuth Service, to gather the access token to use for the Dida365 API
type Dida365OAuthServiceClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
	OAuthLogin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*model.TemporaryRedirectResponse, error)
	OAuthCallback(ctx context.Context, in *OAuthCallbackRequest, opts ...grpc.CallOption) (*model.TemporaryRedirectResponse, error)
	OAuthPrompt(ctx context.Context, in *OAuthPromptRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
}

type dida365OAuthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDida365OAuthServiceClient(cc grpc.ClientConnInterface) Dida365OAuthServiceClient {
	return &dida365OAuthServiceClient{cc}
}

func (c *dida365OAuthServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, Dida365OAuthService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dida365OAuthServiceClient) OAuthLogin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*model.TemporaryRedirectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(model.TemporaryRedirectResponse)
	err := c.cc.Invoke(ctx, Dida365OAuthService_OAuthLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dida365OAuthServiceClient) OAuthCallback(ctx context.Context, in *OAuthCallbackRequest, opts ...grpc.CallOption) (*model.TemporaryRedirectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(model.TemporaryRedirectResponse)
	err := c.cc.Invoke(ctx, Dida365OAuthService_OAuthCallback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dida365OAuthServiceClient) OAuthPrompt(ctx context.Context, in *OAuthPromptRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, Dida365OAuthService_OAuthPrompt_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Dida365OAuthServiceServer is the server API for Dida365OAuthService service.
// All implementations must embed UnimplementedDida365OAuthServiceServer
// for forward compatibility.
//
// Dida365 OAuth Service, to gather the access token to use for the Dida365 API
type Dida365OAuthServiceServer interface {
	Ping(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error)
	OAuthLogin(context.Context, *emptypb.Empty) (*model.TemporaryRedirectResponse, error)
	OAuthCallback(context.Context, *OAuthCallbackRequest) (*model.TemporaryRedirectResponse, error)
	OAuthPrompt(context.Context, *OAuthPromptRequest) (*httpbody.HttpBody, error)
	mustEmbedUnimplementedDida365OAuthServiceServer()
}

// UnimplementedDida365OAuthServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDida365OAuthServiceServer struct{}

func (UnimplementedDida365OAuthServiceServer) Ping(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDida365OAuthServiceServer) OAuthLogin(context.Context, *emptypb.Empty) (*model.TemporaryRedirectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthLogin not implemented")
}
func (UnimplementedDida365OAuthServiceServer) OAuthCallback(context.Context, *OAuthCallbackRequest) (*model.TemporaryRedirectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthCallback not implemented")
}
func (UnimplementedDida365OAuthServiceServer) OAuthPrompt(context.Context, *OAuthPromptRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthPrompt not implemented")
}
func (UnimplementedDida365OAuthServiceServer) mustEmbedUnimplementedDida365OAuthServiceServer() {}
func (UnimplementedDida365OAuthServiceServer) testEmbeddedByValue()                             {}

// UnsafeDida365OAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Dida365OAuthServiceServer will
// result in compilation errors.
type UnsafeDida365OAuthServiceServer interface {
	mustEmbedUnimplementedDida365OAuthServiceServer()
}

func RegisterDida365OAuthServiceServer(s grpc.ServiceRegistrar, srv Dida365OAuthServiceServer) {
	// If the following call pancis, it indicates UnimplementedDida365OAuthServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Dida365OAuthService_ServiceDesc, srv)
}

func _Dida365OAuthService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Dida365OAuthServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dida365OAuthService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Dida365OAuthServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dida365OAuthService_OAuthLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Dida365OAuthServiceServer).OAuthLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dida365OAuthService_OAuthLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Dida365OAuthServiceServer).OAuthLogin(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dida365OAuthService_OAuthCallback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OAuthCallbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Dida365OAuthServiceServer).OAuthCallback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dida365OAuthService_OAuthCallback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Dida365OAuthServiceServer).OAuthCallback(ctx, req.(*OAuthCallbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dida365OAuthService_OAuthPrompt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OAuthPromptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Dida365OAuthServiceServer).OAuthPrompt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dida365OAuthService_OAuthPrompt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Dida365OAuthServiceServer).OAuthPrompt(ctx, req.(*OAuthPromptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Dida365OAuthService_ServiceDesc is the grpc.ServiceDesc for Dida365OAuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Dida365OAuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Dida365oAuthService",
	HandlerType: (*Dida365OAuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Dida365OAuthService_Ping_Handler,
		},
		{
			MethodName: "OAuthLogin",
			Handler:    _Dida365OAuthService_OAuthLogin_Handler,
		},
		{
			MethodName: "OAuthCallback",
			Handler:    _Dida365OAuthService_OAuthCallback_Handler,
		},
		{
			MethodName: "OAuthPrompt",
			Handler:    _Dida365OAuthService_OAuthPrompt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/oauth.proto",
}
