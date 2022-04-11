// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gateway

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

// GatewayClient is the client API for Gateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayClient interface {
	// PrepareAuthentication returns the URL to call to authenticate against the
	// upstream IDP
	RequestUpstreamAuthentication(ctx context.Context, in *UpstreamAuthenticationRequest, opts ...grpc.CallOption) (*UpstreamAuthenticationResponse, error)
	// RequestAuthentication is called to exchange the authorization code with the
	// upstream IDP and to authenticate with the m8 control plane
	RequestAuthentication(ctx context.Context, in *AuthenticationRequest, opts ...grpc.CallOption) (*AuthenticationResponse, error)
}

type gatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayClient(cc grpc.ClientConnInterface) GatewayClient {
	return &gatewayClient{cc}
}

func (c *gatewayClient) RequestUpstreamAuthentication(ctx context.Context, in *UpstreamAuthenticationRequest, opts ...grpc.CallOption) (*UpstreamAuthenticationResponse, error) {
	out := new(UpstreamAuthenticationResponse)
	err := c.cc.Invoke(ctx, "/gateway.Gateway/RequestUpstreamAuthentication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayClient) RequestAuthentication(ctx context.Context, in *AuthenticationRequest, opts ...grpc.CallOption) (*AuthenticationResponse, error) {
	out := new(AuthenticationResponse)
	err := c.cc.Invoke(ctx, "/gateway.Gateway/RequestAuthentication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServer is the server API for Gateway service.
// All implementations must embed UnimplementedGatewayServer
// for forward compatibility
type GatewayServer interface {
	// PrepareAuthentication returns the URL to call to authenticate against the
	// upstream IDP
	RequestUpstreamAuthentication(context.Context, *UpstreamAuthenticationRequest) (*UpstreamAuthenticationResponse, error)
	// RequestAuthentication is called to exchange the authorization code with the
	// upstream IDP and to authenticate with the m8 control plane
	RequestAuthentication(context.Context, *AuthenticationRequest) (*AuthenticationResponse, error)
	mustEmbedUnimplementedGatewayServer()
}

// UnimplementedGatewayServer must be embedded to have forward compatible implementations.
type UnimplementedGatewayServer struct {
}

func (UnimplementedGatewayServer) RequestUpstreamAuthentication(context.Context, *UpstreamAuthenticationRequest) (*UpstreamAuthenticationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestUpstreamAuthentication not implemented")
}
func (UnimplementedGatewayServer) RequestAuthentication(context.Context, *AuthenticationRequest) (*AuthenticationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestAuthentication not implemented")
}
func (UnimplementedGatewayServer) mustEmbedUnimplementedGatewayServer() {}

// UnsafeGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayServer will
// result in compilation errors.
type UnsafeGatewayServer interface {
	mustEmbedUnimplementedGatewayServer()
}

func RegisterGatewayServer(s grpc.ServiceRegistrar, srv GatewayServer) {
	s.RegisterService(&Gateway_ServiceDesc, srv)
}

func _Gateway_RequestUpstreamAuthentication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpstreamAuthenticationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).RequestUpstreamAuthentication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.Gateway/RequestUpstreamAuthentication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).RequestUpstreamAuthentication(ctx, req.(*UpstreamAuthenticationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gateway_RequestAuthentication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).RequestAuthentication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.Gateway/RequestAuthentication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).RequestAuthentication(ctx, req.(*AuthenticationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gateway_ServiceDesc is the grpc.ServiceDesc for Gateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.Gateway",
	HandlerType: (*GatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestUpstreamAuthentication",
			Handler:    _Gateway_RequestUpstreamAuthentication_Handler,
		},
		{
			MethodName: "RequestAuthentication",
			Handler:    _Gateway_RequestAuthentication_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/gateway/service.proto",
}

// GatewayAuthClient is the client API for GatewayAuth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayAuthClient interface {
	// Performs authorization check based on the attributes associated with the
	// incoming request, and returns status `OK` or not `OK`.
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
}

type gatewayAuthClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayAuthClient(cc grpc.ClientConnInterface) GatewayAuthClient {
	return &gatewayAuthClient{cc}
}

func (c *gatewayAuthClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/gateway.GatewayAuth/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayAuthServer is the server API for GatewayAuth service.
// All implementations must embed UnimplementedGatewayAuthServer
// for forward compatibility
type GatewayAuthServer interface {
	// Performs authorization check based on the attributes associated with the
	// incoming request, and returns status `OK` or not `OK`.
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	mustEmbedUnimplementedGatewayAuthServer()
}

// UnimplementedGatewayAuthServer must be embedded to have forward compatible implementations.
type UnimplementedGatewayAuthServer struct {
}

func (UnimplementedGatewayAuthServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedGatewayAuthServer) mustEmbedUnimplementedGatewayAuthServer() {}

// UnsafeGatewayAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayAuthServer will
// result in compilation errors.
type UnsafeGatewayAuthServer interface {
	mustEmbedUnimplementedGatewayAuthServer()
}

func RegisterGatewayAuthServer(s grpc.ServiceRegistrar, srv GatewayAuthServer) {
	s.RegisterService(&GatewayAuth_ServiceDesc, srv)
}

func _GatewayAuth_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAuthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.GatewayAuth/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAuthServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GatewayAuth_ServiceDesc is the grpc.ServiceDesc for GatewayAuth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GatewayAuth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.GatewayAuth",
	HandlerType: (*GatewayAuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _GatewayAuth_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/gateway/service.proto",
}

// ClusterAuthClient is the client API for ClusterAuth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterAuthClient interface {
	GetAuthToken(ctx context.Context, in *ClusterAuthTokenRequest, opts ...grpc.CallOption) (*ClusterAuthTokenResponse, error)
}

type clusterAuthClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterAuthClient(cc grpc.ClientConnInterface) ClusterAuthClient {
	return &clusterAuthClient{cc}
}

func (c *clusterAuthClient) GetAuthToken(ctx context.Context, in *ClusterAuthTokenRequest, opts ...grpc.CallOption) (*ClusterAuthTokenResponse, error) {
	out := new(ClusterAuthTokenResponse)
	err := c.cc.Invoke(ctx, "/gateway.ClusterAuth/GetAuthToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterAuthServer is the server API for ClusterAuth service.
// All implementations must embed UnimplementedClusterAuthServer
// for forward compatibility
type ClusterAuthServer interface {
	GetAuthToken(context.Context, *ClusterAuthTokenRequest) (*ClusterAuthTokenResponse, error)
	mustEmbedUnimplementedClusterAuthServer()
}

// UnimplementedClusterAuthServer must be embedded to have forward compatible implementations.
type UnimplementedClusterAuthServer struct {
}

func (UnimplementedClusterAuthServer) GetAuthToken(context.Context, *ClusterAuthTokenRequest) (*ClusterAuthTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthToken not implemented")
}
func (UnimplementedClusterAuthServer) mustEmbedUnimplementedClusterAuthServer() {}

// UnsafeClusterAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterAuthServer will
// result in compilation errors.
type UnsafeClusterAuthServer interface {
	mustEmbedUnimplementedClusterAuthServer()
}

func RegisterClusterAuthServer(s grpc.ServiceRegistrar, srv ClusterAuthServer) {
	s.RegisterService(&ClusterAuth_ServiceDesc, srv)
}

func _ClusterAuth_GetAuthToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterAuthTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterAuthServer).GetAuthToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.ClusterAuth/GetAuthToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterAuthServer).GetAuthToken(ctx, req.(*ClusterAuthTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterAuth_ServiceDesc is the grpc.ServiceDesc for ClusterAuth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterAuth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.ClusterAuth",
	HandlerType: (*ClusterAuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthToken",
			Handler:    _ClusterAuth_GetAuthToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/gateway/service.proto",
}

// APITokenClient is the client API for APIToken service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APITokenClient interface {
	RequestAPIToken(ctx context.Context, in *APITokenRequest, opts ...grpc.CallOption) (*APITokenResponse, error)
}

type aPITokenClient struct {
	cc grpc.ClientConnInterface
}

func NewAPITokenClient(cc grpc.ClientConnInterface) APITokenClient {
	return &aPITokenClient{cc}
}

func (c *aPITokenClient) RequestAPIToken(ctx context.Context, in *APITokenRequest, opts ...grpc.CallOption) (*APITokenResponse, error) {
	out := new(APITokenResponse)
	err := c.cc.Invoke(ctx, "/gateway.APIToken/RequestAPIToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APITokenServer is the server API for APIToken service.
// All implementations must embed UnimplementedAPITokenServer
// for forward compatibility
type APITokenServer interface {
	RequestAPIToken(context.Context, *APITokenRequest) (*APITokenResponse, error)
	mustEmbedUnimplementedAPITokenServer()
}

// UnimplementedAPITokenServer must be embedded to have forward compatible implementations.
type UnimplementedAPITokenServer struct {
}

func (UnimplementedAPITokenServer) RequestAPIToken(context.Context, *APITokenRequest) (*APITokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestAPIToken not implemented")
}
func (UnimplementedAPITokenServer) mustEmbedUnimplementedAPITokenServer() {}

// UnsafeAPITokenServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APITokenServer will
// result in compilation errors.
type UnsafeAPITokenServer interface {
	mustEmbedUnimplementedAPITokenServer()
}

func RegisterAPITokenServer(s grpc.ServiceRegistrar, srv APITokenServer) {
	s.RegisterService(&APIToken_ServiceDesc, srv)
}

func _APIToken_RequestAPIToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(APITokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APITokenServer).RequestAPIToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.APIToken/RequestAPIToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APITokenServer).RequestAPIToken(ctx, req.(*APITokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// APIToken_ServiceDesc is the grpc.ServiceDesc for APIToken service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var APIToken_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.APIToken",
	HandlerType: (*APITokenServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestAPIToken",
			Handler:    _APIToken_RequestAPIToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/gateway/service.proto",
}
