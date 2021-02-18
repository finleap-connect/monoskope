// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/domain/commanddata/tenant.proto

package commanddata

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
type CreateTenantCommandData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the tenant
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // required
	// Prefix for namespaces and other resources related to the tenant
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"` // required
}

func (x *CreateTenantCommandData) Reset() {
	*x = CreateTenantCommandData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_commanddata_tenant_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTenantCommandData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTenantCommandData) ProtoMessage() {}

func (x *CreateTenantCommandData) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_commanddata_tenant_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTenantCommandData.ProtoReflect.Descriptor instead.
func (*CreateTenantCommandData) Descriptor() ([]byte, []int) {
	return file_api_domain_commanddata_tenant_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTenantCommandData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateTenantCommandData) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

// Request to update a single tenant
type UpdateTenantCommandData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the tenant to update (UUID 128-bit number)
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // required
	// The information to update
	Update *UpdateTenantCommandData_Update `protobuf:"bytes,2,opt,name=update,proto3" json:"update,omitempty"` // required
}

func (x *UpdateTenantCommandData) Reset() {
	*x = UpdateTenantCommandData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_commanddata_tenant_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTenantCommandData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTenantCommandData) ProtoMessage() {}

func (x *UpdateTenantCommandData) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_commanddata_tenant_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTenantCommandData.ProtoReflect.Descriptor instead.
func (*UpdateTenantCommandData) Descriptor() ([]byte, []int) {
	return file_api_domain_commanddata_tenant_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateTenantCommandData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateTenantCommandData) GetUpdate() *UpdateTenantCommandData_Update {
	if x != nil {
		return x.Update
	}
	return nil
}

// Request to delete a single tenant
type DeleteTenantCommandData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the tenant to delete (UUID 128-bit number)
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // required
}

func (x *DeleteTenantCommandData) Reset() {
	*x = DeleteTenantCommandData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_commanddata_tenant_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTenantCommandData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTenantCommandData) ProtoMessage() {}

func (x *DeleteTenantCommandData) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_commanddata_tenant_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTenantCommandData.ProtoReflect.Descriptor instead.
func (*DeleteTenantCommandData) Descriptor() ([]byte, []int) {
	return file_api_domain_commanddata_tenant_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteTenantCommandData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpdateTenantCommandData_Update struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// New name for the tenant
	Name *wrappers.StringValue `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *UpdateTenantCommandData_Update) Reset() {
	*x = UpdateTenantCommandData_Update{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_commanddata_tenant_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTenantCommandData_Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTenantCommandData_Update) ProtoMessage() {}

func (x *UpdateTenantCommandData_Update) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_commanddata_tenant_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTenantCommandData_Update.ProtoReflect.Descriptor instead.
func (*UpdateTenantCommandData_Update) Descriptor() ([]byte, []int) {
	return file_api_domain_commanddata_tenant_proto_rawDescGZIP(), []int{1, 0}
}

func (x *UpdateTenantCommandData_Update) GetName() *wrappers.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

var File_api_domain_commanddata_tenant_proto protoreflect.FileDescriptor

var file_api_domain_commanddata_tenant_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x64, 0x61,
	0x74, 0x61, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x45, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x22, 0xaa, 0x01, 0x0a, 0x17, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x43, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x3a, 0x0a, 0x06, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x29, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f,
	0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f,
	0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_commanddata_tenant_proto_rawDescOnce sync.Once
	file_api_domain_commanddata_tenant_proto_rawDescData = file_api_domain_commanddata_tenant_proto_rawDesc
)

func file_api_domain_commanddata_tenant_proto_rawDescGZIP() []byte {
	file_api_domain_commanddata_tenant_proto_rawDescOnce.Do(func() {
		file_api_domain_commanddata_tenant_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_commanddata_tenant_proto_rawDescData)
	})
	return file_api_domain_commanddata_tenant_proto_rawDescData
}

var file_api_domain_commanddata_tenant_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_domain_commanddata_tenant_proto_goTypes = []interface{}{
	(*CreateTenantCommandData)(nil),        // 0: commanddata.CreateTenantCommandData
	(*UpdateTenantCommandData)(nil),        // 1: commanddata.UpdateTenantCommandData
	(*DeleteTenantCommandData)(nil),        // 2: commanddata.DeleteTenantCommandData
	(*UpdateTenantCommandData_Update)(nil), // 3: commanddata.UpdateTenantCommandData.Update
	(*wrappers.StringValue)(nil),           // 4: google.protobuf.StringValue
}
var file_api_domain_commanddata_tenant_proto_depIdxs = []int32{
	3, // 0: commanddata.UpdateTenantCommandData.update:type_name -> commanddata.UpdateTenantCommandData.Update
	4, // 1: commanddata.UpdateTenantCommandData.Update.name:type_name -> google.protobuf.StringValue
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_domain_commanddata_tenant_proto_init() }
func file_api_domain_commanddata_tenant_proto_init() {
	if File_api_domain_commanddata_tenant_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_domain_commanddata_tenant_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTenantCommandData); i {
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
		file_api_domain_commanddata_tenant_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTenantCommandData); i {
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
		file_api_domain_commanddata_tenant_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTenantCommandData); i {
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
		file_api_domain_commanddata_tenant_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTenantCommandData_Update); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_domain_commanddata_tenant_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_domain_commanddata_tenant_proto_goTypes,
		DependencyIndexes: file_api_domain_commanddata_tenant_proto_depIdxs,
		MessageInfos:      file_api_domain_commanddata_tenant_proto_msgTypes,
	}.Build()
	File_api_domain_commanddata_tenant_proto = out.File
	file_api_domain_commanddata_tenant_proto_rawDesc = nil
	file_api_domain_commanddata_tenant_proto_goTypes = nil
	file_api_domain_commanddata_tenant_proto_depIdxs = nil
}
