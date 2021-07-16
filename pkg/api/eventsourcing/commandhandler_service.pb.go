// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.0
// source: api/eventsourcing/commandhandler_service.proto

package eventsourcing

import (
	commands "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
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

type CommandReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// UUID of the referenced aggregate. If this was a "Create*" command, the ID provied
	// with the command is ignored. A valid ID is generated by the command handler and
	// returned here.
	AggregateId string `protobuf:"bytes,1,opt,name=aggregate_id,json=aggregateId,proto3" json:"aggregate_id,omitempty"` // required
	// version of the aggregate at the time of command being received. Any resulting events
	// applied afterwards to the aggregate will increase this.
	Version uint64 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *CommandReply) Reset() {
	*x = CommandReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eventsourcing_commandhandler_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandReply) ProtoMessage() {}

func (x *CommandReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_eventsourcing_commandhandler_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandReply.ProtoReflect.Descriptor instead.
func (*CommandReply) Descriptor() ([]byte, []int) {
	return file_api_eventsourcing_commandhandler_service_proto_rawDescGZIP(), []int{0}
}

func (x *CommandReply) GetAggregateId() string {
	if x != nil {
		return x.AggregateId
	}
	return ""
}

func (x *CommandReply) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_api_eventsourcing_commandhandler_service_proto protoreflect.FileDescriptor

var file_api_eventsourcing_commandhandler_service_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x69, 0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x69, 0x6e, 0x67, 0x1a,
	0x28, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x69,
	0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a, 0x0c, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x32, 0x4b, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x07, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x1b, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69,
	0x67, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f,
	0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_eventsourcing_commandhandler_service_proto_rawDescOnce sync.Once
	file_api_eventsourcing_commandhandler_service_proto_rawDescData = file_api_eventsourcing_commandhandler_service_proto_rawDesc
)

func file_api_eventsourcing_commandhandler_service_proto_rawDescGZIP() []byte {
	file_api_eventsourcing_commandhandler_service_proto_rawDescOnce.Do(func() {
		file_api_eventsourcing_commandhandler_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_eventsourcing_commandhandler_service_proto_rawDescData)
	})
	return file_api_eventsourcing_commandhandler_service_proto_rawDescData
}

var file_api_eventsourcing_commandhandler_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_eventsourcing_commandhandler_service_proto_goTypes = []interface{}{
	(*CommandReply)(nil),     // 0: eventsourcing.CommandReply
	(*commands.Command)(nil), // 1: commands.Command
}
var file_api_eventsourcing_commandhandler_service_proto_depIdxs = []int32{
	1, // 0: eventsourcing.CommandHandler.Execute:input_type -> commands.Command
	0, // 1: eventsourcing.CommandHandler.Execute:output_type -> eventsourcing.CommandReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_eventsourcing_commandhandler_service_proto_init() }
func file_api_eventsourcing_commandhandler_service_proto_init() {
	if File_api_eventsourcing_commandhandler_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_eventsourcing_commandhandler_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandReply); i {
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
			RawDescriptor: file_api_eventsourcing_commandhandler_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_eventsourcing_commandhandler_service_proto_goTypes,
		DependencyIndexes: file_api_eventsourcing_commandhandler_service_proto_depIdxs,
		MessageInfos:      file_api_eventsourcing_commandhandler_service_proto_msgTypes,
	}.Build()
	File_api_eventsourcing_commandhandler_service_proto = out.File
	file_api_eventsourcing_commandhandler_service_proto_rawDesc = nil
	file_api_eventsourcing_commandhandler_service_proto_goTypes = nil
	file_api_eventsourcing_commandhandler_service_proto_depIdxs = nil
}
