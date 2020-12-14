// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package auth

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	// Auth
	GetAuthInformation(ctx context.Context, in *AuthState, opts ...grpc.CallOption) (*AuthInformation, error)
	ExchangeAuthCode(ctx context.Context, in *AuthCode, opts ...grpc.CallOption) (*AuthResponse, error)
	RefreshAuth(ctx context.Context, in *RefreshAuthRequest, opts ...grpc.CallOption) (*AccessToken, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) GetAuthInformation(ctx context.Context, in *AuthState, opts ...grpc.CallOption) (*AuthInformation, error) {
	out := new(AuthInformation)
	err := c.cc.Invoke(ctx, "/gateway.auth.Auth/GetAuthInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) ExchangeAuthCode(ctx context.Context, in *AuthCode, opts ...grpc.CallOption) (*AuthResponse, error) {
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, "/gateway.auth.Auth/ExchangeAuthCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) RefreshAuth(ctx context.Context, in *RefreshAuthRequest, opts ...grpc.CallOption) (*AccessToken, error) {
	out := new(AccessToken)
	err := c.cc.Invoke(ctx, "/gateway.auth.Auth/RefreshAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	// Auth
	GetAuthInformation(context.Context, *AuthState) (*AuthInformation, error)
	ExchangeAuthCode(context.Context, *AuthCode) (*AuthResponse, error)
	RefreshAuth(context.Context, *RefreshAuthRequest) (*AccessToken, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) GetAuthInformation(context.Context, *AuthState) (*AuthInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthInformation not implemented")
}
func (UnimplementedAuthServer) ExchangeAuthCode(context.Context, *AuthCode) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExchangeAuthCode not implemented")
}
func (UnimplementedAuthServer) RefreshAuth(context.Context, *RefreshAuthRequest) (*AccessToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshAuth not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_GetAuthInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthState)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).GetAuthInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.auth.Auth/GetAuthInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).GetAuthInformation(ctx, req.(*AuthState))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_ExchangeAuthCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).ExchangeAuthCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.auth.Auth/ExchangeAuthCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).ExchangeAuthCode(ctx, req.(*AuthCode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_RefreshAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).RefreshAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.auth.Auth/RefreshAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).RefreshAuth(ctx, req.(*RefreshAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthInformation",
			Handler:    _Auth_GetAuthInformation_Handler,
		},
		{
			MethodName: "ExchangeAuthCode",
			Handler:    _Auth_ExchangeAuthCode_Handler,
		},
		{
			MethodName: "RefreshAuth",
			Handler:    _Auth_RefreshAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/gateway/auth/service.proto",
}
