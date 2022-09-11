// Copyright 2022 Monoskope Authors
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
// source: api/domain/eventdata/certificate.proto

package eventdata

import (
	common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
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

type CertificateRequested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the aggregate referenced (UUID 128-bit number)
	ReferencedAggregateId string `protobuf:"bytes,1,opt,name=referenced_aggregate_id,json=referencedAggregateId,proto3" json:"referenced_aggregate_id,omitempty"`
	// Type of the aggregate referenced
	ReferencedAggregateType string `protobuf:"bytes,4,opt,name=referenced_aggregate_type,json=referencedAggregateType,proto3" json:"referenced_aggregate_type,omitempty"`
	// Certificate signing request
	SigningRequest []byte `protobuf:"bytes,3,opt,name=signing_request,json=signingRequest,proto3" json:"signing_request,omitempty"`
}

func (x *CertificateRequested) Reset() {
	*x = CertificateRequested{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_certificate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CertificateRequested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateRequested) ProtoMessage() {}

func (x *CertificateRequested) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_certificate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateRequested.ProtoReflect.Descriptor instead.
func (*CertificateRequested) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_certificate_proto_rawDescGZIP(), []int{0}
}

func (x *CertificateRequested) GetReferencedAggregateId() string {
	if x != nil {
		return x.ReferencedAggregateId
	}
	return ""
}

func (x *CertificateRequested) GetReferencedAggregateType() string {
	if x != nil {
		return x.ReferencedAggregateType
	}
	return ""
}

func (x *CertificateRequested) GetSigningRequest() []byte {
	if x != nil {
		return x.SigningRequest
	}
	return nil
}

type CertificateIssued struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Certificate *common.CertificateChain `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"`
}

func (x *CertificateIssued) Reset() {
	*x = CertificateIssued{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_certificate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CertificateIssued) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateIssued) ProtoMessage() {}

func (x *CertificateIssued) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_certificate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateIssued.ProtoReflect.Descriptor instead.
func (*CertificateIssued) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_certificate_proto_rawDescGZIP(), []int{1}
}

func (x *CertificateIssued) GetCertificate() *common.CertificateChain {
	if x != nil {
		return x.Certificate
	}
	return nil
}

var File_api_domain_eventdata_certificate_proto protoreflect.FileDescriptor

var file_api_domain_eventdata_certificate_proto_rawDesc = []byte{
	0x0a, 0x26, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x64,
	0x61, 0x74, 0x61, 0x1a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb3, 0x01, 0x0a, 0x14, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x36,
	0x0a, 0x17, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x5f, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x15, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x41, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x19, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x64, 0x5f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x72, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x64, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x73, 0x69, 0x67,
	0x6e, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4f, 0x0a, 0x11, 0x43,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64,
	0x12, 0x3a, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x52,
	0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x42, 0x3f, 0x5a, 0x3d,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6e, 0x6c, 0x65,
	0x61, 0x70, 0x2d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73,
	0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_eventdata_certificate_proto_rawDescOnce sync.Once
	file_api_domain_eventdata_certificate_proto_rawDescData = file_api_domain_eventdata_certificate_proto_rawDesc
)

func file_api_domain_eventdata_certificate_proto_rawDescGZIP() []byte {
	file_api_domain_eventdata_certificate_proto_rawDescOnce.Do(func() {
		file_api_domain_eventdata_certificate_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_eventdata_certificate_proto_rawDescData)
	})
	return file_api_domain_eventdata_certificate_proto_rawDescData
}

var file_api_domain_eventdata_certificate_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_domain_eventdata_certificate_proto_goTypes = []interface{}{
	(*CertificateRequested)(nil),    // 0: eventdata.CertificateRequested
	(*CertificateIssued)(nil),       // 1: eventdata.CertificateIssued
	(*common.CertificateChain)(nil), // 2: common.CertificateChain
}
var file_api_domain_eventdata_certificate_proto_depIdxs = []int32{
	2, // 0: eventdata.CertificateIssued.certificate:type_name -> common.CertificateChain
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_domain_eventdata_certificate_proto_init() }
func file_api_domain_eventdata_certificate_proto_init() {
	if File_api_domain_eventdata_certificate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_domain_eventdata_certificate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CertificateRequested); i {
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
		file_api_domain_eventdata_certificate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CertificateIssued); i {
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
			RawDescriptor: file_api_domain_eventdata_certificate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_domain_eventdata_certificate_proto_goTypes,
		DependencyIndexes: file_api_domain_eventdata_certificate_proto_depIdxs,
		MessageInfos:      file_api_domain_eventdata_certificate_proto_msgTypes,
	}.Build()
	File_api_domain_eventdata_certificate_proto = out.File
	file_api_domain_eventdata_certificate_proto_rawDesc = nil
	file_api_domain_eventdata_certificate_proto_goTypes = nil
	file_api_domain_eventdata_certificate_proto_depIdxs = nil
}