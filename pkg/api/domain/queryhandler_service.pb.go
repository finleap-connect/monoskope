// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: api/domain/queryhandler_service.proto

package domain

import (
	proto "github.com/golang/protobuf/proto"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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

var File_api_domain_queryhandler_service_proto protoreflect.FileDescriptor

var file_api_domain_queryhandler_service_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x61, 0x70,
	0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0xe2, 0x02, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x30, 0x01, 0x12, 0x3a, 0x0a, 0x07, 0x47, 0x65,
	0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x1a, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x79, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x1a, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x53, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65,
	0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c,
	0x65, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x30, 0x01, 0x12, 0x4c, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x41, 0x75, 0x64, 0x69, 0x74, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x42, 0x79, 0x49, 0x64, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x75, 0x64, 0x69,
	0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x32, 0x94, 0x02, 0x0a, 0x0d, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x13, 0x2e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x12, 0x3e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x13, 0x2e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x12, 0x4c, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41, 0x75, 0x64, 0x69, 0x74, 0x54, 0x72, 0x61,
	0x69, 0x6c, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x42,
	0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f,
	0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b,
	0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_domain_queryhandler_service_proto_goTypes = []interface{}{
	(*emptypb.Empty)(nil),               // 0: google.protobuf.Empty
	(*wrapperspb.StringValue)(nil),      // 1: google.protobuf.StringValue
	(*projections.User)(nil),            // 2: projections.User
	(*projections.UserRoleBinding)(nil), // 3: projections.UserRoleBinding
	(*projections.AuditEvent)(nil),      // 4: projections.AuditEvent
	(*projections.Tenant)(nil),          // 5: projections.Tenant
}
var file_api_domain_queryhandler_service_proto_depIdxs = []int32{
	0, // 0: domain.UserService.GetAll:input_type -> google.protobuf.Empty
	1, // 1: domain.UserService.GetById:input_type -> google.protobuf.StringValue
	1, // 2: domain.UserService.GetByEmail:input_type -> google.protobuf.StringValue
	1, // 3: domain.UserService.GetRoleBindingsById:input_type -> google.protobuf.StringValue
	1, // 4: domain.UserService.GetAuditTrailById:input_type -> google.protobuf.StringValue
	0, // 5: domain.TenantService.GetAll:input_type -> google.protobuf.Empty
	1, // 6: domain.TenantService.GetById:input_type -> google.protobuf.StringValue
	1, // 7: domain.TenantService.GetByName:input_type -> google.protobuf.StringValue
	1, // 8: domain.TenantService.GetAuditTrailById:input_type -> google.protobuf.StringValue
	2, // 9: domain.UserService.GetAll:output_type -> projections.User
	2, // 10: domain.UserService.GetById:output_type -> projections.User
	2, // 11: domain.UserService.GetByEmail:output_type -> projections.User
	3, // 12: domain.UserService.GetRoleBindingsById:output_type -> projections.UserRoleBinding
	4, // 13: domain.UserService.GetAuditTrailById:output_type -> projections.AuditEvent
	5, // 14: domain.TenantService.GetAll:output_type -> projections.Tenant
	5, // 15: domain.TenantService.GetById:output_type -> projections.Tenant
	5, // 16: domain.TenantService.GetByName:output_type -> projections.Tenant
	4, // 17: domain.TenantService.GetAuditTrailById:output_type -> projections.AuditEvent
	9, // [9:18] is the sub-list for method output_type
	0, // [0:9] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_domain_queryhandler_service_proto_init() }
func file_api_domain_queryhandler_service_proto_init() {
	if File_api_domain_queryhandler_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_domain_queryhandler_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_domain_queryhandler_service_proto_goTypes,
		DependencyIndexes: file_api_domain_queryhandler_service_proto_depIdxs,
	}.Build()
	File_api_domain_queryhandler_service_proto = out.File
	file_api_domain_queryhandler_service_proto_rawDesc = nil
	file_api_domain_queryhandler_service_proto_goTypes = nil
	file_api_domain_queryhandler_service_proto_depIdxs = nil
}
