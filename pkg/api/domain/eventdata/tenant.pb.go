// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.7
// source: api/domain/eventdata/tenant.proto

package eventdata

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TenantCreatedEventData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the tenant
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Prefix of the tenant
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
}

func (x *TenantCreatedEventData) Reset() {
	*x = TenantCreatedEventData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_tenant_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantCreatedEventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantCreatedEventData) ProtoMessage() {}

func (x *TenantCreatedEventData) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_tenant_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantCreatedEventData.ProtoReflect.Descriptor instead.
func (*TenantCreatedEventData) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_tenant_proto_rawDescGZIP(), []int{0}
}

func (x *TenantCreatedEventData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TenantCreatedEventData) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

type TenantUpdatedEventData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// New name for the tenant
	Name *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *TenantUpdatedEventData) Reset() {
	*x = TenantUpdatedEventData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_tenant_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantUpdatedEventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantUpdatedEventData) ProtoMessage() {}

func (x *TenantUpdatedEventData) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_tenant_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantUpdatedEventData.ProtoReflect.Descriptor instead.
func (*TenantUpdatedEventData) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_tenant_proto_rawDescGZIP(), []int{1}
}

func (x *TenantUpdatedEventData) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

var File_api_domain_eventdata_tenant_proto protoreflect.FileDescriptor

var file_api_domain_eventdata_tenant_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x1e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x44,
	0x0a, 0x16, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x22, 0x4a, 0x0a, 0x16, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x30,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73,
	0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_eventdata_tenant_proto_rawDescOnce sync.Once
	file_api_domain_eventdata_tenant_proto_rawDescData = file_api_domain_eventdata_tenant_proto_rawDesc
)

func file_api_domain_eventdata_tenant_proto_rawDescGZIP() []byte {
	file_api_domain_eventdata_tenant_proto_rawDescOnce.Do(func() {
		file_api_domain_eventdata_tenant_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_eventdata_tenant_proto_rawDescData)
	})
	return file_api_domain_eventdata_tenant_proto_rawDescData
}

var file_api_domain_eventdata_tenant_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_domain_eventdata_tenant_proto_goTypes = []interface{}{
	(*TenantCreatedEventData)(nil), // 0: eventdata.TenantCreatedEventData
	(*TenantUpdatedEventData)(nil), // 1: eventdata.TenantUpdatedEventData
	(*wrapperspb.StringValue)(nil), // 2: google.protobuf.StringValue
}
var file_api_domain_eventdata_tenant_proto_depIdxs = []int32{
	2, // 0: eventdata.TenantUpdatedEventData.name:type_name -> google.protobuf.StringValue
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_domain_eventdata_tenant_proto_init() }
func file_api_domain_eventdata_tenant_proto_init() {
	if File_api_domain_eventdata_tenant_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_domain_eventdata_tenant_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantCreatedEventData); i {
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
		file_api_domain_eventdata_tenant_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantUpdatedEventData); i {
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
			RawDescriptor: file_api_domain_eventdata_tenant_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_domain_eventdata_tenant_proto_goTypes,
		DependencyIndexes: file_api_domain_eventdata_tenant_proto_depIdxs,
		MessageInfos:      file_api_domain_eventdata_tenant_proto_msgTypes,
	}.Build()
	File_api_domain_eventdata_tenant_proto = out.File
	file_api_domain_eventdata_tenant_proto_rawDesc = nil
	file_api_domain_eventdata_tenant_proto_goTypes = nil
	file_api_domain_eventdata_tenant_proto_depIdxs = nil
}
