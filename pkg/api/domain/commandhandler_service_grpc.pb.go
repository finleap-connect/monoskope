// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package domain

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CommandHandlerExtensionsClient is the client API for CommandHandlerExtensions service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandHandlerExtensionsClient interface {
	// Returns roles and scopes available.
	GetPermissionModel(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PermissionModel, error)
}

type commandHandlerExtensionsClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandHandlerExtensionsClient(cc grpc.ClientConnInterface) CommandHandlerExtensionsClient {
	return &commandHandlerExtensionsClient{cc}
}

func (c *commandHandlerExtensionsClient) GetPermissionModel(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PermissionModel, error) {
	out := new(PermissionModel)
	err := c.cc.Invoke(ctx, "/domain.CommandHandlerExtensions/GetPermissionModel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandHandlerExtensionsServer is the server API for CommandHandlerExtensions service.
// All implementations must embed UnimplementedCommandHandlerExtensionsServer
// for forward compatibility
type CommandHandlerExtensionsServer interface {
	// Returns roles and scopes available.
	GetPermissionModel(context.Context, *emptypb.Empty) (*PermissionModel, error)
	mustEmbedUnimplementedCommandHandlerExtensionsServer()
}

// UnimplementedCommandHandlerExtensionsServer must be embedded to have forward compatible implementations.
type UnimplementedCommandHandlerExtensionsServer struct {
}

func (UnimplementedCommandHandlerExtensionsServer) GetPermissionModel(context.Context, *emptypb.Empty) (*PermissionModel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermissionModel not implemented")
}
func (UnimplementedCommandHandlerExtensionsServer) mustEmbedUnimplementedCommandHandlerExtensionsServer() {
}

// UnsafeCommandHandlerExtensionsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandHandlerExtensionsServer will
// result in compilation errors.
type UnsafeCommandHandlerExtensionsServer interface {
	mustEmbedUnimplementedCommandHandlerExtensionsServer()
}

func RegisterCommandHandlerExtensionsServer(s grpc.ServiceRegistrar, srv CommandHandlerExtensionsServer) {
	s.RegisterService(&CommandHandlerExtensions_ServiceDesc, srv)
}

func _CommandHandlerExtensions_GetPermissionModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandHandlerExtensionsServer).GetPermissionModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.CommandHandlerExtensions/GetPermissionModel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandHandlerExtensionsServer).GetPermissionModel(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// CommandHandlerExtensions_ServiceDesc is the grpc.ServiceDesc for CommandHandlerExtensions service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommandHandlerExtensions_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "domain.CommandHandlerExtensions",
	HandlerType: (*CommandHandlerExtensionsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPermissionModel",
			Handler:    _CommandHandlerExtensions_GetPermissionModel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/domain/commandhandler_service.proto",
}
