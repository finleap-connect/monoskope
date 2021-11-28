// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.0
// source: api/domain/projections/certificate.proto

package projections

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

// Certficates can exist for any aggregate and are mainted by the control plane.
type Certificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the certificate (UUID 128-bit number)
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Unique identifier of the referenced aggregate (UUID 128-bit number)
	ReferencedAggregateId string `protobuf:"bytes,2,opt,name=referenced_aggregate_id,json=referencedAggregateId,proto3" json:"referenced_aggregate_id,omitempty"`
	// Type of the aggregate referenced
	AggregateType string `protobuf:"bytes,3,opt,name=aggregate_type,json=aggregateType,proto3" json:"aggregate_type,omitempty"`
	// actual certificate
	Certificate []byte `protobuf:"bytes,4,opt,name=certificate,proto3" json:"certificate,omitempty"`
	// bundle of CA certificates for the certificate managed
	CaCertBundle []byte `protobuf:"bytes,5,opt,name=ca_cert_bundle,json=caCertBundle,proto3" json:"ca_cert_bundle,omitempty"`
	// Metadata about the projection
	Metadata *LifecycleMetadata `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *Certificate) Reset() {
	*x = Certificate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_projections_certificate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Certificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certificate) ProtoMessage() {}

func (x *Certificate) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_projections_certificate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certificate.ProtoReflect.Descriptor instead.
func (*Certificate) Descriptor() ([]byte, []int) {
	return file_api_domain_projections_certificate_proto_rawDescGZIP(), []int{0}
}

func (x *Certificate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Certificate) GetReferencedAggregateId() string {
	if x != nil {
		return x.ReferencedAggregateId
	}
	return ""
}

func (x *Certificate) GetAggregateType() string {
	if x != nil {
		return x.AggregateType
	}
	return ""
}

func (x *Certificate) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

func (x *Certificate) GetCaCertBundle() []byte {
	if x != nil {
		return x.CaCertBundle
	}
	return nil
}

func (x *Certificate) GetMetadata() *LifecycleMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_api_domain_projections_certificate_proto protoreflect.FileDescriptor

var file_api_domain_projections_certificate_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb7, 0x02, 0x0a, 0x0b, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x40, 0x0a, 0x17, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x5f,
	0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x15, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x65, 0x49, 0x64, 0x12, 0x48, 0x0a, 0x0e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0xfa, 0x42, 0x1e,
	0x72, 0x1c, 0x28, 0x3c, 0x32, 0x18, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x5d, 0x5b,
	0x41, 0x2d, 0x5a, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x5f, 0x2d, 0x5d, 0x2b, 0x24, 0x52, 0x0d,
	0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12,
	0x24, 0x0a, 0x0e, 0x63, 0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x62, 0x75, 0x6e, 0x64, 0x6c,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x61, 0x43, 0x65, 0x72, 0x74, 0x42,
	0x75, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x3a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x69, 0x6e, 0x6c, 0x65, 0x61, 0x70, 0x2d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_projections_certificate_proto_rawDescOnce sync.Once
	file_api_domain_projections_certificate_proto_rawDescData = file_api_domain_projections_certificate_proto_rawDesc
)

func file_api_domain_projections_certificate_proto_rawDescGZIP() []byte {
	file_api_domain_projections_certificate_proto_rawDescOnce.Do(func() {
		file_api_domain_projections_certificate_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_projections_certificate_proto_rawDescData)
	})
	return file_api_domain_projections_certificate_proto_rawDescData
}

var file_api_domain_projections_certificate_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_domain_projections_certificate_proto_goTypes = []interface{}{
	(*Certificate)(nil),       // 0: projections.Certificate
	(*LifecycleMetadata)(nil), // 1: projections.LifecycleMetadata
}
var file_api_domain_projections_certificate_proto_depIdxs = []int32{
	1, // 0: projections.Certificate.metadata:type_name -> projections.LifecycleMetadata
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_domain_projections_certificate_proto_init() }
func file_api_domain_projections_certificate_proto_init() {
	if File_api_domain_projections_certificate_proto != nil {
		return
	}
	file_api_domain_projections_metadata_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_domain_projections_certificate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Certificate); i {
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
			RawDescriptor: file_api_domain_projections_certificate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_domain_projections_certificate_proto_goTypes,
		DependencyIndexes: file_api_domain_projections_certificate_proto_depIdxs,
		MessageInfos:      file_api_domain_projections_certificate_proto_msgTypes,
	}.Build()
	File_api_domain_projections_certificate_proto = out.File
	file_api_domain_projections_certificate_proto_rawDesc = nil
	file_api_domain_projections_certificate_proto_goTypes = nil
	file_api_domain_projections_certificate_proto_depIdxs = nil
}
