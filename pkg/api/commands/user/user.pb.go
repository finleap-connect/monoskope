// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/commands/user.proto

package user

import (
	proto "github.com/golang/protobuf/proto"
	common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
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

type CreateUserCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user information
	UserMetadata *common.UserMetadata `protobuf:"bytes,1,opt,name=user_metadata,json=userMetadata,proto3" json:"user_metadata,omitempty"`
}

func (x *CreateUserCommand) Reset() {
	*x = CreateUserCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_commands_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserCommand) ProtoMessage() {}

func (x *CreateUserCommand) ProtoReflect() protoreflect.Message {
	mi := &file_api_commands_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserCommand.ProtoReflect.Descriptor instead.
func (*CreateUserCommand) Descriptor() ([]byte, []int) {
	return file_api_commands_user_proto_rawDescGZIP(), []int{0}
}

func (x *CreateUserCommand) GetUserMetadata() *common.UserMetadata {
	if x != nil {
		return x.UserMetadata
	}
	return nil
}

type AddRoleToUserCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the user (UUID 128-bit number)
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// Name of the role to add
	Role string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	// Scope of the role
	Scope string `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	// Affected resource within scope
	Resource string `protobuf:"bytes,4,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *AddRoleToUserCommand) Reset() {
	*x = AddRoleToUserCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_commands_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRoleToUserCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRoleToUserCommand) ProtoMessage() {}

func (x *AddRoleToUserCommand) ProtoReflect() protoreflect.Message {
	mi := &file_api_commands_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRoleToUserCommand.ProtoReflect.Descriptor instead.
func (*AddRoleToUserCommand) Descriptor() ([]byte, []int) {
	return file_api_commands_user_proto_rawDescGZIP(), []int{1}
}

func (x *AddRoleToUserCommand) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AddRoleToUserCommand) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *AddRoleToUserCommand) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *AddRoleToUserCommand) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

var File_api_commands_user_proto protoreflect.FileDescriptor

var file_api_commands_user_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x1a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x39, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x75, 0x0a, 0x14, 0x41, 0x64, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f,
	0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73,
	0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_commands_user_proto_rawDescOnce sync.Once
	file_api_commands_user_proto_rawDescData = file_api_commands_user_proto_rawDesc
)

func file_api_commands_user_proto_rawDescGZIP() []byte {
	file_api_commands_user_proto_rawDescOnce.Do(func() {
		file_api_commands_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_commands_user_proto_rawDescData)
	})
	return file_api_commands_user_proto_rawDescData
}

var file_api_commands_user_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_commands_user_proto_goTypes = []interface{}{
	(*CreateUserCommand)(nil),    // 0: commands.user.CreateUserCommand
	(*AddRoleToUserCommand)(nil), // 1: commands.user.AddRoleToUserCommand
	(*common.UserMetadata)(nil),  // 2: common.UserMetadata
}
var file_api_commands_user_proto_depIdxs = []int32{
	2, // 0: commands.user.CreateUserCommand.user_metadata:type_name -> common.UserMetadata
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_commands_user_proto_init() }
func file_api_commands_user_proto_init() {
	if File_api_commands_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_commands_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserCommand); i {
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
		file_api_commands_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRoleToUserCommand); i {
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
			RawDescriptor: file_api_commands_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_commands_user_proto_goTypes,
		DependencyIndexes: file_api_commands_user_proto_depIdxs,
		MessageInfos:      file_api_commands_user_proto_msgTypes,
	}.Build()
	File_api_commands_user_proto = out.File
	file_api_commands_user_proto_rawDesc = nil
	file_api_commands_user_proto_goTypes = nil
	file_api_commands_user_proto_depIdxs = nil
}
