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
// source: api/domain/queryhandler_service.proto

package domain

import (
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
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

// GetAllRequest is the generic request to query all instances of a certain projection
type GetAllRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IncludeDeleted bool `protobuf:"varint,1,opt,name=include_deleted,json=includeDeleted,proto3" json:"include_deleted,omitempty"`
}

func (x *GetAllRequest) Reset() {
	*x = GetAllRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_queryhandler_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllRequest) ProtoMessage() {}

func (x *GetAllRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_queryhandler_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllRequest.ProtoReflect.Descriptor instead.
func (*GetAllRequest) Descriptor() ([]byte, []int) {
	return file_api_domain_queryhandler_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAllRequest) GetIncludeDeleted() bool {
	if x != nil {
		return x.IncludeDeleted
	}
	return false
}

type GetCertificateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the aggregate referenced (UUID 128-bit number)
	AggregateId string `protobuf:"bytes,1,opt,name=aggregate_id,json=aggregateId,proto3" json:"aggregate_id,omitempty"`
	// Type of the aggregate referenced
	AggregateType string `protobuf:"bytes,2,opt,name=aggregate_type,json=aggregateType,proto3" json:"aggregate_type,omitempty"`
}

func (x *GetCertificateRequest) Reset() {
	*x = GetCertificateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_queryhandler_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCertificateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCertificateRequest) ProtoMessage() {}

func (x *GetCertificateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_queryhandler_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCertificateRequest.ProtoReflect.Descriptor instead.
func (*GetCertificateRequest) Descriptor() ([]byte, []int) {
	return file_api_domain_queryhandler_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetCertificateRequest) GetAggregateId() string {
	if x != nil {
		return x.AggregateId
	}
	return ""
}

func (x *GetCertificateRequest) GetAggregateType() string {
	if x != nil {
		return x.AggregateType
	}
	return ""
}

type GetClusterMappingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId  string `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	ClusterId string `protobuf:"bytes,2,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
}

func (x *GetClusterMappingRequest) Reset() {
	*x = GetClusterMappingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_queryhandler_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClusterMappingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterMappingRequest) ProtoMessage() {}

func (x *GetClusterMappingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_queryhandler_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterMappingRequest.ProtoReflect.Descriptor instead.
func (*GetClusterMappingRequest) Descriptor() ([]byte, []int) {
	return file_api_domain_queryhandler_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetClusterMappingRequest) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *GetClusterMappingRequest) GetClusterId() string {
	if x != nil {
		return x.ClusterId
	}
	return ""
}

type GetCountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IncludeDeleted bool `protobuf:"varint,1,opt,name=include_deleted,json=includeDeleted,proto3" json:"include_deleted,omitempty"`
}

func (x *GetCountRequest) Reset() {
	*x = GetCountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_queryhandler_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCountRequest) ProtoMessage() {}

func (x *GetCountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_queryhandler_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCountRequest.ProtoReflect.Descriptor instead.
func (*GetCountRequest) Descriptor() ([]byte, []int) {
	return file_api_domain_queryhandler_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetCountRequest) GetIncludeDeleted() bool {
	if x != nil {
		return x.IncludeDeleted
	}
	return false
}

type GetCountResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GetCountResult) Reset() {
	*x = GetCountResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_queryhandler_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCountResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCountResult) ProtoMessage() {}

func (x *GetCountResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_queryhandler_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCountResult.ProtoReflect.Descriptor instead.
func (*GetCountResult) Descriptor() ([]byte, []int) {
	return file_api_domain_queryhandler_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetCountResult) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_api_domain_queryhandler_service_proto protoreflect.FileDescriptor

var file_api_domain_queryhandler_service_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x1a,
	0x21, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x23, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x61,
	0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x33, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x62,
	0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x38, 0x0a, 0x0d,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a,
	0x0f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x61, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x21, 0x0a, 0x0c, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65,
	0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x56, 0x0a, 0x18, 0x47, 0x65, 0x74,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x3a, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69,
	0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x26, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0xc9, 0x02, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x34,
	0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x15, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x11, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x30, 0x01, 0x12, 0x3a, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x11, 0x2e,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x3d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x11, 0x2e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x53, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x73, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x30, 0x01, 0x12, 0x3b, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x17, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x32, 0x83, 0x02, 0x0a, 0x06, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x36, 0x0a, 0x06,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x15, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x13, 0x2e,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x12, 0x3e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x13, 0x2e,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x12, 0x43, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x17, 0x2e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x30, 0x01, 0x32, 0x93, 0x02, 0x0a, 0x07, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x15, 0x2e,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x30, 0x01, 0x12, 0x3d, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x4f, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xfc, 0x03,
	0x0a, 0x0d, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x52, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x42, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x30, 0x01, 0x12, 0x50, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x14, 0x2e,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x30, 0x01, 0x12, 0x67, 0x0a, 0x22, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x73, 0x42, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x30, 0x01, 0x12, 0x68,
	0x0a, 0x23, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73, 0x42, 0x79, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x42,
	0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x30, 0x01, 0x12, 0x72, 0x0a, 0x2b, 0x47, 0x65, 0x74, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x42, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x41, 0x6e, 0x64, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x70, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x32, 0x58, 0x0a, 0x0b,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x49, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6e, 0x6c, 0x65, 0x61, 0x70, 0x2d, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_queryhandler_service_proto_rawDescOnce sync.Once
	file_api_domain_queryhandler_service_proto_rawDescData = file_api_domain_queryhandler_service_proto_rawDesc
)

func file_api_domain_queryhandler_service_proto_rawDescGZIP() []byte {
	file_api_domain_queryhandler_service_proto_rawDescOnce.Do(func() {
		file_api_domain_queryhandler_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_queryhandler_service_proto_rawDescData)
	})
	return file_api_domain_queryhandler_service_proto_rawDescData
}

var file_api_domain_queryhandler_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_domain_queryhandler_service_proto_goTypes = []interface{}{
	(*GetAllRequest)(nil),                    // 0: domain.GetAllRequest
	(*GetCertificateRequest)(nil),            // 1: domain.GetCertificateRequest
	(*GetClusterMappingRequest)(nil),         // 2: domain.GetClusterMappingRequest
	(*GetCountRequest)(nil),                  // 3: domain.GetCountRequest
	(*GetCountResult)(nil),                   // 4: domain.GetCountResult
	(*wrapperspb.StringValue)(nil),           // 5: google.protobuf.StringValue
	(*projections.User)(nil),                 // 6: projections.User
	(*projections.UserRoleBinding)(nil),      // 7: projections.UserRoleBinding
	(*projections.Tenant)(nil),               // 8: projections.Tenant
	(*projections.TenantUser)(nil),           // 9: projections.TenantUser
	(*projections.Cluster)(nil),              // 10: projections.Cluster
	(*projections.TenantClusterBinding)(nil), // 11: projections.TenantClusterBinding
	(*projections.Certificate)(nil),          // 12: projections.Certificate
}
var file_api_domain_queryhandler_service_proto_depIdxs = []int32{
	0,  // 0: domain.User.GetAll:input_type -> domain.GetAllRequest
	5,  // 1: domain.User.GetById:input_type -> google.protobuf.StringValue
	5,  // 2: domain.User.GetByEmail:input_type -> google.protobuf.StringValue
	5,  // 3: domain.User.GetRoleBindingsById:input_type -> google.protobuf.StringValue
	3,  // 4: domain.User.GetCount:input_type -> domain.GetCountRequest
	0,  // 5: domain.Tenant.GetAll:input_type -> domain.GetAllRequest
	5,  // 6: domain.Tenant.GetById:input_type -> google.protobuf.StringValue
	5,  // 7: domain.Tenant.GetByName:input_type -> google.protobuf.StringValue
	5,  // 8: domain.Tenant.GetUsers:input_type -> google.protobuf.StringValue
	0,  // 9: domain.Cluster.GetAll:input_type -> domain.GetAllRequest
	5,  // 10: domain.Cluster.GetById:input_type -> google.protobuf.StringValue
	5,  // 11: domain.Cluster.GetByName:input_type -> google.protobuf.StringValue
	5,  // 12: domain.Cluster.GetBootstrapToken:input_type -> google.protobuf.StringValue
	5,  // 13: domain.ClusterAccess.GetClusterAccessByTenantId:input_type -> google.protobuf.StringValue
	5,  // 14: domain.ClusterAccess.GetClusterAccessByUserId:input_type -> google.protobuf.StringValue
	5,  // 15: domain.ClusterAccess.GetTenantClusterMappingsByTenantId:input_type -> google.protobuf.StringValue
	5,  // 16: domain.ClusterAccess.GetTenantClusterMappingsByClusterId:input_type -> google.protobuf.StringValue
	2,  // 17: domain.ClusterAccess.GetTenantClusterMappingByTenantAndClusterId:input_type -> domain.GetClusterMappingRequest
	1,  // 18: domain.Certificate.GetCertificate:input_type -> domain.GetCertificateRequest
	6,  // 19: domain.User.GetAll:output_type -> projections.User
	6,  // 20: domain.User.GetById:output_type -> projections.User
	6,  // 21: domain.User.GetByEmail:output_type -> projections.User
	7,  // 22: domain.User.GetRoleBindingsById:output_type -> projections.UserRoleBinding
	4,  // 23: domain.User.GetCount:output_type -> domain.GetCountResult
	8,  // 24: domain.Tenant.GetAll:output_type -> projections.Tenant
	8,  // 25: domain.Tenant.GetById:output_type -> projections.Tenant
	8,  // 26: domain.Tenant.GetByName:output_type -> projections.Tenant
	9,  // 27: domain.Tenant.GetUsers:output_type -> projections.TenantUser
	10, // 28: domain.Cluster.GetAll:output_type -> projections.Cluster
	10, // 29: domain.Cluster.GetById:output_type -> projections.Cluster
	10, // 30: domain.Cluster.GetByName:output_type -> projections.Cluster
	5,  // 31: domain.Cluster.GetBootstrapToken:output_type -> google.protobuf.StringValue
	10, // 32: domain.ClusterAccess.GetClusterAccessByTenantId:output_type -> projections.Cluster
	10, // 33: domain.ClusterAccess.GetClusterAccessByUserId:output_type -> projections.Cluster
	11, // 34: domain.ClusterAccess.GetTenantClusterMappingsByTenantId:output_type -> projections.TenantClusterBinding
	11, // 35: domain.ClusterAccess.GetTenantClusterMappingsByClusterId:output_type -> projections.TenantClusterBinding
	11, // 36: domain.ClusterAccess.GetTenantClusterMappingByTenantAndClusterId:output_type -> projections.TenantClusterBinding
	12, // 37: domain.Certificate.GetCertificate:output_type -> projections.Certificate
	19, // [19:38] is the sub-list for method output_type
	0,  // [0:19] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_domain_queryhandler_service_proto_init() }
func file_api_domain_queryhandler_service_proto_init() {
	if File_api_domain_queryhandler_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_domain_queryhandler_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllRequest); i {
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
		file_api_domain_queryhandler_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCertificateRequest); i {
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
		file_api_domain_queryhandler_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClusterMappingRequest); i {
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
		file_api_domain_queryhandler_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCountRequest); i {
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
		file_api_domain_queryhandler_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCountResult); i {
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
			RawDescriptor: file_api_domain_queryhandler_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   5,
		},
		GoTypes:           file_api_domain_queryhandler_service_proto_goTypes,
		DependencyIndexes: file_api_domain_queryhandler_service_proto_depIdxs,
		MessageInfos:      file_api_domain_queryhandler_service_proto_msgTypes,
	}.Build()
	File_api_domain_queryhandler_service_proto = out.File
	file_api_domain_queryhandler_service_proto_rawDesc = nil
	file_api_domain_queryhandler_service_proto_goTypes = nil
	file_api_domain_queryhandler_service_proto_depIdxs = nil
}
