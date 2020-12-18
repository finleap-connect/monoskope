// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/gateway/tenant/messages.proto

package tenant

import (
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Request to create a new tenant
type CreateTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the tenant
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // required
	// Prefix for namespaces and other resources related to the tenant
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"` // required
}

func (x *CreateTenantRequest) Reset() {
	*x = CreateTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_tenant_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTenantRequest) ProtoMessage() {}

func (x *CreateTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_tenant_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTenantRequest.ProtoReflect.Descriptor instead.
func (*CreateTenantRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_tenant_messages_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTenantRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateTenantRequest) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

// Request to get a single tenant
type GetTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to FindBy:
	//	*GetTenantRequest_Id
	//	*GetTenantRequest_Name
	FindBy isGetTenantRequest_FindBy `protobuf_oneof:"find_by"`
}

func (x *GetTenantRequest) Reset() {
	*x = GetTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_tenant_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTenantRequest) ProtoMessage() {}

func (x *GetTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_tenant_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTenantRequest.ProtoReflect.Descriptor instead.
func (*GetTenantRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_tenant_messages_proto_rawDescGZIP(), []int{1}
}

func (m *GetTenantRequest) GetFindBy() isGetTenantRequest_FindBy {
	if m != nil {
		return m.FindBy
	}
	return nil
}

func (x *GetTenantRequest) GetId() string {
	if x, ok := x.GetFindBy().(*GetTenantRequest_Id); ok {
		return x.Id
	}
	return ""
}

func (x *GetTenantRequest) GetName() string {
	if x, ok := x.GetFindBy().(*GetTenantRequest_Name); ok {
		return x.Name
	}
	return ""
}

type isGetTenantRequest_FindBy interface {
	isGetTenantRequest_FindBy()
}

type GetTenantRequest_Id struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3,oneof"`
}

type GetTenantRequest_Name struct {
	Name string `protobuf:"bytes,2,opt,name=name,proto3,oneof"`
}

func (*GetTenantRequest_Id) isGetTenantRequest_FindBy() {}

func (*GetTenantRequest_Name) isGetTenantRequest_FindBy() {}

// Request to update a single tenant
type UpdateTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the tenant to update
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // required
	// The information to update
	Update *UpdateTenantRequest_Update `protobuf:"bytes,2,opt,name=update,proto3" json:"update,omitempty"` // required
}

func (x *UpdateTenantRequest) Reset() {
	*x = UpdateTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_tenant_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTenantRequest) ProtoMessage() {}

func (x *UpdateTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_tenant_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTenantRequest.ProtoReflect.Descriptor instead.
func (*UpdateTenantRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_tenant_messages_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateTenantRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateTenantRequest) GetUpdate() *UpdateTenantRequest_Update {
	if x != nil {
		return x.Update
	}
	return nil
}

// Request to delete a single tenant
type DeleteTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // required
}

func (x *DeleteTenantRequest) Reset() {
	*x = DeleteTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_tenant_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTenantRequest) ProtoMessage() {}

func (x *DeleteTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_tenant_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTenantRequest.ProtoReflect.Descriptor instead.
func (*DeleteTenantRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_tenant_messages_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteTenantRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Request to list all tenants
type ListTenantsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Include tenants which have previously been deleted
	IncludeDeleted bool `protobuf:"varint,1,opt,name=include_deleted,json=includeDeleted,proto3" json:"include_deleted,omitempty"`
}

func (x *ListTenantsRequest) Reset() {
	*x = ListTenantsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_tenant_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTenantsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTenantsRequest) ProtoMessage() {}

func (x *ListTenantsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_tenant_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTenantsRequest.ProtoReflect.Descriptor instead.
func (*ListTenantsRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_tenant_messages_proto_rawDescGZIP(), []int{4}
}

func (x *ListTenantsRequest) GetIncludeDeleted() bool {
	if x != nil {
		return x.IncludeDeleted
	}
	return false
}

type UpdateTenantRequest_Update struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// New name for the tenant
	Name *wrappers.StringValue `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // optional
}

func (x *UpdateTenantRequest_Update) Reset() {
	*x = UpdateTenantRequest_Update{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_tenant_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTenantRequest_Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTenantRequest_Update) ProtoMessage() {}

func (x *UpdateTenantRequest_Update) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_tenant_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTenantRequest_Update.ProtoReflect.Descriptor instead.
func (*UpdateTenantRequest_Update) Descriptor() ([]byte, []int) {
	return file_api_gateway_tenant_messages_proto_rawDescGZIP(), []int{2, 0}
}

func (x *UpdateTenantRequest_Update) GetName() *wrappers.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

var File_api_gateway_tenant_messages_proto protoreflect.FileDescriptor

var file_api_gateway_tenant_messages_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x22, 0x45, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x66, 0x69, 0x6e, 0x64, 0x5f, 0x62, 0x79, 0x22, 0xa5, 0x01,
	0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x3a, 0x0a, 0x06, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x25, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3d, 0x0a, 0x12,
	0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x6e, 0x63,
	0x6c, 0x75, 0x64, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x42, 0x4f, 0x5a, 0x47, 0x67,
	0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f,
	0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0xa2, 0x02, 0x03, 0x4d, 0x47, 0x54, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_gateway_tenant_messages_proto_rawDescOnce sync.Once
	file_api_gateway_tenant_messages_proto_rawDescData = file_api_gateway_tenant_messages_proto_rawDesc
)

func file_api_gateway_tenant_messages_proto_rawDescGZIP() []byte {
	file_api_gateway_tenant_messages_proto_rawDescOnce.Do(func() {
		file_api_gateway_tenant_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_gateway_tenant_messages_proto_rawDescData)
	})
	return file_api_gateway_tenant_messages_proto_rawDescData
}

var file_api_gateway_tenant_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_gateway_tenant_messages_proto_goTypes = []interface{}{
	(*CreateTenantRequest)(nil),        // 0: gateway.tenant.CreateTenantRequest
	(*GetTenantRequest)(nil),           // 1: gateway.tenant.GetTenantRequest
	(*UpdateTenantRequest)(nil),        // 2: gateway.tenant.UpdateTenantRequest
	(*DeleteTenantRequest)(nil),        // 3: gateway.tenant.DeleteTenantRequest
	(*ListTenantsRequest)(nil),         // 4: gateway.tenant.ListTenantsRequest
	(*UpdateTenantRequest_Update)(nil), // 5: gateway.tenant.UpdateTenantRequest.Update
	(*wrappers.StringValue)(nil),       // 6: google.protobuf.StringValue
}
var file_api_gateway_tenant_messages_proto_depIdxs = []int32{
	5, // 0: gateway.tenant.UpdateTenantRequest.update:type_name -> gateway.tenant.UpdateTenantRequest.Update
	6, // 1: gateway.tenant.UpdateTenantRequest.Update.name:type_name -> google.protobuf.StringValue
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_gateway_tenant_messages_proto_init() }
func file_api_gateway_tenant_messages_proto_init() {
	if File_api_gateway_tenant_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_gateway_tenant_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTenantRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_tenant_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTenantRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_tenant_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTenantRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_tenant_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTenantRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_tenant_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTenantsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_tenant_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTenantRequest_Update); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_api_gateway_tenant_messages_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*GetTenantRequest_Id)(nil),
		(*GetTenantRequest_Name)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_gateway_tenant_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_gateway_tenant_messages_proto_goTypes,
		DependencyIndexes: file_api_gateway_tenant_messages_proto_depIdxs,
		MessageInfos:      file_api_gateway_tenant_messages_proto_msgTypes,
	}.Build()
	File_api_gateway_tenant_messages_proto = out.File
	file_api_gateway_tenant_messages_proto_rawDesc = nil
	file_api_gateway_tenant_messages_proto_goTypes = nil
	file_api_gateway_tenant_messages_proto_depIdxs = nil
}
