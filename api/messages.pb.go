// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/messages.proto

package api

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// Tenant within Monoskope
type Tenant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the tenant
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the tenant
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Display name of the tenant.
	// Can contain blanks and other special characters.
	DisplayName string `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// Prefix for namespaces and other resources related to the tenant.
	// DNS compatibility is ensured on validation. E.g. no more than 12
	// characters.
	Prefix string `protobuf:"bytes,4,opt,name=prefix,proto3" json:"prefix,omitempty"`
	// When tenant has been created
	Created *timestamp.Timestamp `protobuf:"bytes,5,opt,name=Created,proto3" json:"Created,omitempty"`
	// By whom the tenant has been created
	CreatedBy *User `protobuf:"bytes,6,opt,name=CreatedBy,proto3" json:"CreatedBy,omitempty"`
	// When tenant has been last modified
	LastModified *timestamp.Timestamp `protobuf:"bytes,7,opt,name=LastModified,proto3" json:"LastModified,omitempty"`
	// By whom the tenant has been last modified
	LastModifiedBy *User `protobuf:"bytes,8,opt,name=LastModifiedBy,proto3" json:"LastModifiedBy,omitempty"`
	// When tenant has been deleted
	Deleted *timestamp.Timestamp `protobuf:"bytes,9,opt,name=Deleted,proto3" json:"Deleted,omitempty"`
	// By whom the tenant has been deleted
	DeletedBy *User `protobuf:"bytes,10,opt,name=DeletedBy,proto3" json:"DeletedBy,omitempty"`
}

func (x *Tenant) Reset() {
	*x = Tenant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tenant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tenant) ProtoMessage() {}

func (x *Tenant) ProtoReflect() protoreflect.Message {
	mi := &file_api_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tenant.ProtoReflect.Descriptor instead.
func (*Tenant) Descriptor() ([]byte, []int) {
	return file_api_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Tenant) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tenant) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tenant) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Tenant) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *Tenant) GetCreated() *timestamp.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Tenant) GetCreatedBy() *User {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Tenant) GetLastModified() *timestamp.Timestamp {
	if x != nil {
		return x.LastModified
	}
	return nil
}

func (x *Tenant) GetLastModifiedBy() *User {
	if x != nil {
		return x.LastModifiedBy
	}
	return nil
}

func (x *Tenant) GetDeleted() *timestamp.Timestamp {
	if x != nil {
		return x.Deleted
	}
	return nil
}

func (x *Tenant) GetDeletedBy() *User {
	if x != nil {
		return x.DeletedBy
	}
	return nil
}

// User within Monoskope
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the user
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the user
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_api_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_api_messages_proto_rawDescGZIP(), []int{1}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_api_messages_proto protoreflect.FileDescriptor

var file_api_messages_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x03, 0x0a, 0x06, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x12, 0x34, 0x0a, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x27, 0x0a, 0x09, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x42, 0x79, 0x12, 0x3e, 0x0a, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x12, 0x31, 0x0a, 0x0e, 0x4c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x42, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x0e, 0x4c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69,
	0x66, 0x69, 0x65, 0x64, 0x42, 0x79, 0x12, 0x34, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x27, 0x0a, 0x09,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x42, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x09, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x42, 0x79, 0x22, 0x2a, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x42, 0x3b, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f,
	0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f,
	0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x61, 0x70, 0x69, 0xa2, 0x02, 0x02, 0x4d, 0x54, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_messages_proto_rawDescOnce sync.Once
	file_api_messages_proto_rawDescData = file_api_messages_proto_rawDesc
)

func file_api_messages_proto_rawDescGZIP() []byte {
	file_api_messages_proto_rawDescOnce.Do(func() {
		file_api_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_messages_proto_rawDescData)
	})
	return file_api_messages_proto_rawDescData
}

var file_api_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_messages_proto_goTypes = []interface{}{
	(*Tenant)(nil),              // 0: api.Tenant
	(*User)(nil),                // 1: api.User
	(*timestamp.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_api_messages_proto_depIdxs = []int32{
	2, // 0: api.Tenant.Created:type_name -> google.protobuf.Timestamp
	1, // 1: api.Tenant.CreatedBy:type_name -> api.User
	2, // 2: api.Tenant.LastModified:type_name -> google.protobuf.Timestamp
	1, // 3: api.Tenant.LastModifiedBy:type_name -> api.User
	2, // 4: api.Tenant.Deleted:type_name -> google.protobuf.Timestamp
	1, // 5: api.Tenant.DeletedBy:type_name -> api.User
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_api_messages_proto_init() }
func file_api_messages_proto_init() {
	if File_api_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tenant); i {
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
		file_api_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
			RawDescriptor: file_api_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_messages_proto_goTypes,
		DependencyIndexes: file_api_messages_proto_depIdxs,
		MessageInfos:      file_api_messages_proto_msgTypes,
	}.Build()
	File_api_messages_proto = out.File
	file_api_messages_proto_rawDesc = nil
	file_api_messages_proto_goTypes = nil
	file_api_messages_proto_depIdxs = nil
}
