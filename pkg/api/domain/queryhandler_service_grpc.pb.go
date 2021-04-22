// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package domain

import (
	context "context"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	eventsourcing "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	// GetAll returns all users.
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (User_GetAllClient, error)
	// GetById returns the user found by the given id.
	GetById(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.User, error)
	// GetByEmail returns the user found by the given email address.
	GetByEmail(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.User, error)
	// GetRoleBindingsById returns all role bindings related to the given user id.
	GetRoleBindingsById(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (User_GetRoleBindingsByIdClient, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (User_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &User_ServiceDesc.Streams[0], "/domain.User/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &userGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type User_GetAllClient interface {
	Recv() (*projections.User, error)
	grpc.ClientStream
}

type userGetAllClient struct {
	grpc.ClientStream
}

func (x *userGetAllClient) Recv() (*projections.User, error) {
	m := new(projections.User)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userClient) GetById(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.User, error) {
	out := new(projections.User)
	err := c.cc.Invoke(ctx, "/domain.User/GetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetByEmail(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.User, error) {
	out := new(projections.User)
	err := c.cc.Invoke(ctx, "/domain.User/GetByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetRoleBindingsById(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (User_GetRoleBindingsByIdClient, error) {
	stream, err := c.cc.NewStream(ctx, &User_ServiceDesc.Streams[1], "/domain.User/GetRoleBindingsById", opts...)
	if err != nil {
		return nil, err
	}
	x := &userGetRoleBindingsByIdClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type User_GetRoleBindingsByIdClient interface {
	Recv() (*projections.UserRoleBinding, error)
	grpc.ClientStream
}

type userGetRoleBindingsByIdClient struct {
	grpc.ClientStream
}

func (x *userGetRoleBindingsByIdClient) Recv() (*projections.UserRoleBinding, error) {
	m := new(projections.UserRoleBinding)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	// GetAll returns all users.
	GetAll(*GetAllRequest, User_GetAllServer) error
	// GetById returns the user found by the given id.
	GetById(context.Context, *wrapperspb.StringValue) (*projections.User, error)
	// GetByEmail returns the user found by the given email address.
	GetByEmail(context.Context, *wrapperspb.StringValue) (*projections.User, error)
	// GetRoleBindingsById returns all role bindings related to the given user id.
	GetRoleBindingsById(*wrapperspb.StringValue, User_GetRoleBindingsByIdServer) error
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) GetAll(*GetAllRequest, User_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedUserServer) GetById(context.Context, *wrapperspb.StringValue) (*projections.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedUserServer) GetByEmail(context.Context, *wrapperspb.StringValue) (*projections.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByEmail not implemented")
}
func (UnimplementedUserServer) GetRoleBindingsById(*wrapperspb.StringValue, User_GetRoleBindingsByIdServer) error {
	return status.Errorf(codes.Unimplemented, "method GetRoleBindingsById not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServer).GetAll(m, &userGetAllServer{stream})
}

type User_GetAllServer interface {
	Send(*projections.User) error
	grpc.ServerStream
}

type userGetAllServer struct {
	grpc.ServerStream
}

func (x *userGetAllServer) Send(m *projections.User) error {
	return x.ServerStream.SendMsg(m)
}

func _User_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.User/GetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetById(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.User/GetByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetByEmail(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetRoleBindingsById_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(wrapperspb.StringValue)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServer).GetRoleBindingsById(m, &userGetRoleBindingsByIdServer{stream})
}

type User_GetRoleBindingsByIdServer interface {
	Send(*projections.UserRoleBinding) error
	grpc.ServerStream
}

type userGetRoleBindingsByIdServer struct {
	grpc.ServerStream
}

func (x *userGetRoleBindingsByIdServer) Send(m *projections.UserRoleBinding) error {
	return x.ServerStream.SendMsg(m)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "domain.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetById",
			Handler:    _User_GetById_Handler,
		},
		{
			MethodName: "GetByEmail",
			Handler:    _User_GetByEmail_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _User_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetRoleBindingsById",
			Handler:       _User_GetRoleBindingsById_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/domain/queryhandler_service.proto",
}

// TenantClient is the client API for Tenant service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TenantClient interface {
	// GetAll returns all tenants.
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (Tenant_GetAllClient, error)
	// GetById returns the tenant found by the given id.
	GetById(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.Tenant, error)
	// GetByName returns the tenant found by the given name
	GetByName(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.Tenant, error)
}

type tenantClient struct {
	cc grpc.ClientConnInterface
}

func NewTenantClient(cc grpc.ClientConnInterface) TenantClient {
	return &tenantClient{cc}
}

func (c *tenantClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (Tenant_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tenant_ServiceDesc.Streams[0], "/domain.Tenant/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &tenantGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tenant_GetAllClient interface {
	Recv() (*projections.Tenant, error)
	grpc.ClientStream
}

type tenantGetAllClient struct {
	grpc.ClientStream
}

func (x *tenantGetAllClient) Recv() (*projections.Tenant, error) {
	m := new(projections.Tenant)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tenantClient) GetById(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.Tenant, error) {
	out := new(projections.Tenant)
	err := c.cc.Invoke(ctx, "/domain.Tenant/GetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantClient) GetByName(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*projections.Tenant, error) {
	out := new(projections.Tenant)
	err := c.cc.Invoke(ctx, "/domain.Tenant/GetByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TenantServer is the server API for Tenant service.
// All implementations must embed UnimplementedTenantServer
// for forward compatibility
type TenantServer interface {
	// GetAll returns all tenants.
	GetAll(*GetAllRequest, Tenant_GetAllServer) error
	// GetById returns the tenant found by the given id.
	GetById(context.Context, *wrapperspb.StringValue) (*projections.Tenant, error)
	// GetByName returns the tenant found by the given name
	GetByName(context.Context, *wrapperspb.StringValue) (*projections.Tenant, error)
	mustEmbedUnimplementedTenantServer()
}

// UnimplementedTenantServer must be embedded to have forward compatible implementations.
type UnimplementedTenantServer struct {
}

func (UnimplementedTenantServer) GetAll(*GetAllRequest, Tenant_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedTenantServer) GetById(context.Context, *wrapperspb.StringValue) (*projections.Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedTenantServer) GetByName(context.Context, *wrapperspb.StringValue) (*projections.Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedTenantServer) mustEmbedUnimplementedTenantServer() {}

// UnsafeTenantServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TenantServer will
// result in compilation errors.
type UnsafeTenantServer interface {
	mustEmbedUnimplementedTenantServer()
}

func RegisterTenantServer(s grpc.ServiceRegistrar, srv TenantServer) {
	s.RegisterService(&Tenant_ServiceDesc, srv)
}

func _Tenant_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TenantServer).GetAll(m, &tenantGetAllServer{stream})
}

type Tenant_GetAllServer interface {
	Send(*projections.Tenant) error
	grpc.ServerStream
}

type tenantGetAllServer struct {
	grpc.ServerStream
}

func (x *tenantGetAllServer) Send(m *projections.Tenant) error {
	return x.ServerStream.SendMsg(m)
}

func _Tenant_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.Tenant/GetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServer).GetById(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tenant_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.Tenant/GetByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServer).GetByName(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// Tenant_ServiceDesc is the grpc.ServiceDesc for Tenant service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tenant_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "domain.Tenant",
	HandlerType: (*TenantServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetById",
			Handler:    _Tenant_GetById_Handler,
		},
		{
			MethodName: "GetByName",
			Handler:    _Tenant_GetByName_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _Tenant_GetAll_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/domain/queryhandler_service.proto",
}

// ClusterRegistrationRequestClient is the client API for ClusterRegistrationRequest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterRegistrationRequestClient interface {
	// GetAll returns all registration requests.
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (ClusterRegistrationRequest_GetAllClient, error)
	// GetPending returns all pending registration requests.
	GetPending(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (ClusterRegistrationRequest_GetPendingClient, error)
}

type clusterRegistrationRequestClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterRegistrationRequestClient(cc grpc.ClientConnInterface) ClusterRegistrationRequestClient {
	return &clusterRegistrationRequestClient{cc}
}

func (c *clusterRegistrationRequestClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (ClusterRegistrationRequest_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &ClusterRegistrationRequest_ServiceDesc.Streams[0], "/domain.ClusterRegistrationRequest/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &clusterRegistrationRequestGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ClusterRegistrationRequest_GetAllClient interface {
	Recv() (*projections.ClusterRegistration, error)
	grpc.ClientStream
}

type clusterRegistrationRequestGetAllClient struct {
	grpc.ClientStream
}

func (x *clusterRegistrationRequestGetAllClient) Recv() (*projections.ClusterRegistration, error) {
	m := new(projections.ClusterRegistration)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *clusterRegistrationRequestClient) GetPending(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (ClusterRegistrationRequest_GetPendingClient, error) {
	stream, err := c.cc.NewStream(ctx, &ClusterRegistrationRequest_ServiceDesc.Streams[1], "/domain.ClusterRegistrationRequest/GetPending", opts...)
	if err != nil {
		return nil, err
	}
	x := &clusterRegistrationRequestGetPendingClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ClusterRegistrationRequest_GetPendingClient interface {
	Recv() (*projections.ClusterRegistration, error)
	grpc.ClientStream
}

type clusterRegistrationRequestGetPendingClient struct {
	grpc.ClientStream
}

func (x *clusterRegistrationRequestGetPendingClient) Recv() (*projections.ClusterRegistration, error) {
	m := new(projections.ClusterRegistration)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClusterRegistrationRequestServer is the server API for ClusterRegistrationRequest service.
// All implementations must embed UnimplementedClusterRegistrationRequestServer
// for forward compatibility
type ClusterRegistrationRequestServer interface {
	// GetAll returns all registration requests.
	GetAll(*GetAllRequest, ClusterRegistrationRequest_GetAllServer) error
	// GetPending returns all pending registration requests.
	GetPending(*emptypb.Empty, ClusterRegistrationRequest_GetPendingServer) error
	mustEmbedUnimplementedClusterRegistrationRequestServer()
}

// UnimplementedClusterRegistrationRequestServer must be embedded to have forward compatible implementations.
type UnimplementedClusterRegistrationRequestServer struct {
}

func (UnimplementedClusterRegistrationRequestServer) GetAll(*GetAllRequest, ClusterRegistrationRequest_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedClusterRegistrationRequestServer) GetPending(*emptypb.Empty, ClusterRegistrationRequest_GetPendingServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPending not implemented")
}
func (UnimplementedClusterRegistrationRequestServer) mustEmbedUnimplementedClusterRegistrationRequestServer() {
}

// UnsafeClusterRegistrationRequestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterRegistrationRequestServer will
// result in compilation errors.
type UnsafeClusterRegistrationRequestServer interface {
	mustEmbedUnimplementedClusterRegistrationRequestServer()
}

func RegisterClusterRegistrationRequestServer(s grpc.ServiceRegistrar, srv ClusterRegistrationRequestServer) {
	s.RegisterService(&ClusterRegistrationRequest_ServiceDesc, srv)
}

func _ClusterRegistrationRequest_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ClusterRegistrationRequestServer).GetAll(m, &clusterRegistrationRequestGetAllServer{stream})
}

type ClusterRegistrationRequest_GetAllServer interface {
	Send(*projections.ClusterRegistration) error
	grpc.ServerStream
}

type clusterRegistrationRequestGetAllServer struct {
	grpc.ServerStream
}

func (x *clusterRegistrationRequestGetAllServer) Send(m *projections.ClusterRegistration) error {
	return x.ServerStream.SendMsg(m)
}

func _ClusterRegistrationRequest_GetPending_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ClusterRegistrationRequestServer).GetPending(m, &clusterRegistrationRequestGetPendingServer{stream})
}

type ClusterRegistrationRequest_GetPendingServer interface {
	Send(*projections.ClusterRegistration) error
	grpc.ServerStream
}

type clusterRegistrationRequestGetPendingServer struct {
	grpc.ServerStream
}

func (x *clusterRegistrationRequestGetPendingServer) Send(m *projections.ClusterRegistration) error {
	return x.ServerStream.SendMsg(m)
}

// ClusterRegistrationRequest_ServiceDesc is the grpc.ServiceDesc for ClusterRegistrationRequest service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterRegistrationRequest_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "domain.ClusterRegistrationRequest",
	HandlerType: (*ClusterRegistrationRequestServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _ClusterRegistrationRequest_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetPending",
			Handler:       _ClusterRegistrationRequest_GetPending_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/domain/queryhandler_service.proto",
}

// K8SOperatorClient is the client API for K8SOperator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type K8SOperatorClient interface {
	GetCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*projections.Cluster, error)
	Retrieve(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (K8SOperator_RetrieveClient, error)
}

type k8SOperatorClient struct {
	cc grpc.ClientConnInterface
}

func NewK8SOperatorClient(cc grpc.ClientConnInterface) K8SOperatorClient {
	return &k8SOperatorClient{cc}
}

func (c *k8SOperatorClient) GetCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*projections.Cluster, error) {
	out := new(projections.Cluster)
	err := c.cc.Invoke(ctx, "/domain.K8sOperator/GetCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *k8SOperatorClient) Retrieve(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (K8SOperator_RetrieveClient, error) {
	stream, err := c.cc.NewStream(ctx, &K8SOperator_ServiceDesc.Streams[0], "/domain.K8sOperator/Retrieve", opts...)
	if err != nil {
		return nil, err
	}
	x := &k8SOperatorRetrieveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type K8SOperator_RetrieveClient interface {
	Recv() (*eventsourcing.Event, error)
	grpc.ClientStream
}

type k8SOperatorRetrieveClient struct {
	grpc.ClientStream
}

func (x *k8SOperatorRetrieveClient) Recv() (*eventsourcing.Event, error) {
	m := new(eventsourcing.Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// K8SOperatorServer is the server API for K8SOperator service.
// All implementations must embed UnimplementedK8SOperatorServer
// for forward compatibility
type K8SOperatorServer interface {
	GetCluster(context.Context, *emptypb.Empty) (*projections.Cluster, error)
	Retrieve(*emptypb.Empty, K8SOperator_RetrieveServer) error
	mustEmbedUnimplementedK8SOperatorServer()
}

// UnimplementedK8SOperatorServer must be embedded to have forward compatible implementations.
type UnimplementedK8SOperatorServer struct {
}

func (UnimplementedK8SOperatorServer) GetCluster(context.Context, *emptypb.Empty) (*projections.Cluster, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCluster not implemented")
}
func (UnimplementedK8SOperatorServer) Retrieve(*emptypb.Empty, K8SOperator_RetrieveServer) error {
	return status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}
func (UnimplementedK8SOperatorServer) mustEmbedUnimplementedK8SOperatorServer() {}

// UnsafeK8SOperatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to K8SOperatorServer will
// result in compilation errors.
type UnsafeK8SOperatorServer interface {
	mustEmbedUnimplementedK8SOperatorServer()
}

func RegisterK8SOperatorServer(s grpc.ServiceRegistrar, srv K8SOperatorServer) {
	s.RegisterService(&K8SOperator_ServiceDesc, srv)
}

func _K8SOperator_GetCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(K8SOperatorServer).GetCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.K8sOperator/GetCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(K8SOperatorServer).GetCluster(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _K8SOperator_Retrieve_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(K8SOperatorServer).Retrieve(m, &k8SOperatorRetrieveServer{stream})
}

type K8SOperator_RetrieveServer interface {
	Send(*eventsourcing.Event) error
	grpc.ServerStream
}

type k8SOperatorRetrieveServer struct {
	grpc.ServerStream
}

func (x *k8SOperatorRetrieveServer) Send(m *eventsourcing.Event) error {
	return x.ServerStream.SendMsg(m)
}

// K8SOperator_ServiceDesc is the grpc.ServiceDesc for K8SOperator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var K8SOperator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "domain.K8sOperator",
	HandlerType: (*K8SOperatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCluster",
			Handler:    _K8SOperator_GetCluster_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Retrieve",
			Handler:       _K8SOperator_Retrieve_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/domain/queryhandler_service.proto",
}
