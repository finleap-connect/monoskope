// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/eventstore/service.proto

package eventstore

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
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

var File_api_eventstore_service_proto protoreflect.FileDescriptor

var file_api_eventstore_service_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x7c, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x11, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x28, 0x01, 0x12, 0x38, 0x0a, 0x08, 0x52, 0x65,
	0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x12, 0x17, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a,
	0x11, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x30, 0x01, 0x42, 0x46, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66,
	0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d,
	0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0xa2, 0x02, 0x02, 0x4d, 0x47, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_api_eventstore_service_proto_goTypes = []interface{}{
	(*Event)(nil),       // 0: eventstore.Event
	(*EventFilter)(nil), // 1: eventstore.EventFilter
	(*empty.Empty)(nil), // 2: google.protobuf.Empty
}
var file_api_eventstore_service_proto_depIdxs = []int32{
	0, // 0: eventstore.EventStore.Store:input_type -> eventstore.Event
	1, // 1: eventstore.EventStore.Retrieve:input_type -> eventstore.EventFilter
	2, // 2: eventstore.EventStore.Store:output_type -> google.protobuf.Empty
	0, // 3: eventstore.EventStore.Retrieve:output_type -> eventstore.Event
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_eventstore_service_proto_init() }
func file_api_eventstore_service_proto_init() {
	if File_api_eventstore_service_proto != nil {
		return
	}
	file_api_eventstore_messages_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_eventstore_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_eventstore_service_proto_goTypes,
		DependencyIndexes: file_api_eventstore_service_proto_depIdxs,
	}.Build()
	File_api_eventstore_service_proto = out.File
	file_api_eventstore_service_proto_rawDesc = nil
	file_api_eventstore_service_proto_goTypes = nil
	file_api_eventstore_service_proto_depIdxs = nil
}
