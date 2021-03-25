// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.6
// source: api/eventsourcing/commands/command.proto

package commands

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
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

// Command is a command to be executed by the CommandHandler
type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the aggregate the command applies to (UUID 128-bit
	// number)
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // required
	// Type of the command
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// Command type specific data
	Data *anypb.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eventsourcing_commands_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_api_eventsourcing_commands_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_api_eventsourcing_commands_command_proto_rawDescGZIP(), []int{0}
}

func (x *Command) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Command) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Command) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

type TestCommandData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Test      string `protobuf:"bytes,1,opt,name=test,proto3" json:"test,omitempty"`
	TestCount int32  `protobuf:"varint,2,opt,name=test_count,json=testCount,proto3" json:"test_count,omitempty"`
}

func (x *TestCommandData) Reset() {
	*x = TestCommandData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eventsourcing_commands_command_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCommandData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCommandData) ProtoMessage() {}

func (x *TestCommandData) ProtoReflect() protoreflect.Message {
	mi := &file_api_eventsourcing_commands_command_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCommandData.ProtoReflect.Descriptor instead.
func (*TestCommandData) Descriptor() ([]byte, []int) {
	return file_api_eventsourcing_commands_command_proto_rawDescGZIP(), []int{1}
}

func (x *TestCommandData) GetTest() string {
	if x != nil {
		return x.Test
	}
	return ""
}

func (x *TestCommandData) GetTestCount() int32 {
	if x != nil {
		return x.TestCount
	}
	return 0
}

var File_api_eventsourcing_commands_command_proto protoreflect.FileDescriptor

var file_api_eventsourcing_commands_command_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x69, 0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x73, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x57, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x28,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x44, 0x0a, 0x0f, 0x54, 0x65, 0x73, 0x74,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x73, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x51,
	0x5a, 0x4f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d,
	0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f,
	0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_eventsourcing_commands_command_proto_rawDescOnce sync.Once
	file_api_eventsourcing_commands_command_proto_rawDescData = file_api_eventsourcing_commands_command_proto_rawDesc
)

func file_api_eventsourcing_commands_command_proto_rawDescGZIP() []byte {
	file_api_eventsourcing_commands_command_proto_rawDescOnce.Do(func() {
		file_api_eventsourcing_commands_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_eventsourcing_commands_command_proto_rawDescData)
	})
	return file_api_eventsourcing_commands_command_proto_rawDescData
}

var file_api_eventsourcing_commands_command_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_eventsourcing_commands_command_proto_goTypes = []interface{}{
	(*Command)(nil),         // 0: commands.Command
	(*TestCommandData)(nil), // 1: commands.TestCommandData
	(*anypb.Any)(nil),       // 2: google.protobuf.Any
}
var file_api_eventsourcing_commands_command_proto_depIdxs = []int32{
	2, // 0: commands.Command.data:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_eventsourcing_commands_command_proto_init() }
func file_api_eventsourcing_commands_command_proto_init() {
	if File_api_eventsourcing_commands_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_eventsourcing_commands_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
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
		file_api_eventsourcing_commands_command_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCommandData); i {
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
			RawDescriptor: file_api_eventsourcing_commands_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_eventsourcing_commands_command_proto_goTypes,
		DependencyIndexes: file_api_eventsourcing_commands_command_proto_depIdxs,
		MessageInfos:      file_api_eventsourcing_commands_command_proto_msgTypes,
	}.Build()
	File_api_eventsourcing_commands_command_proto = out.File
	file_api_eventsourcing_commands_command_proto_rawDesc = nil
	file_api_eventsourcing_commands_command_proto_goTypes = nil
	file_api_eventsourcing_commands_command_proto_depIdxs = nil
}
