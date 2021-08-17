// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.0
// source: api/domain/projections/cluster.proto

package projections

import (
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

// Cluster is the information the Control Plane has about a cluster
type Cluster struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the cluster (UUID 128-bit number)
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Unique name of the cluster, to be utilized for generating unique labels
	// and symbols, e.g. with metrics.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Display name of the cluster
	DisplayName string `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// Address of the clusters KubeAPIServer
	ApiServerAddress string `protobuf:"bytes,4,opt,name=api_server_address,json=apiServerAddress,proto3" json:"api_server_address,omitempty"`
	// CA certificates of the cluster
	CaCertBundle []byte `protobuf:"bytes,5,opt,name=ca_cert_bundle,json=caCertBundle,proto3" json:"ca_cert_bundle,omitempty"`
	// Metadata about the projection
	Metadata *LifecycleMetadata `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Bootstrap token for cluster authentication
	BootstrapToken string `protobuf:"bytes,7,opt,name=bootstrap_token,json=bootstrapToken,proto3" json:"bootstrap_token,omitempty"`
}

func (x *Cluster) Reset() {
	*x = Cluster{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_projections_cluster_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cluster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cluster) ProtoMessage() {}

func (x *Cluster) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_projections_cluster_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cluster.ProtoReflect.Descriptor instead.
func (*Cluster) Descriptor() ([]byte, []int) {
	return file_api_domain_projections_cluster_proto_rawDescGZIP(), []int{0}
}

func (x *Cluster) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Cluster) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Cluster) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Cluster) GetApiServerAddress() string {
	if x != nil {
		return x.ApiServerAddress
	}
	return ""
}

func (x *Cluster) GetCaCertBundle() []byte {
	if x != nil {
		return x.CaCertBundle
	}
	return nil
}

func (x *Cluster) GetMetadata() *LifecycleMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Cluster) GetBootstrapToken() string {
	if x != nil {
		return x.BootstrapToken
	}
	return ""
}

var File_api_domain_projections_cluster_proto protoreflect.FileDescriptor

var file_api_domain_projections_cluster_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x1a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x02, 0x0a, 0x07, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a,
	0x12, 0x61, 0x70, 0x69, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x61, 0x70, 0x69, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a, 0x0e, 0x63,
	0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x62, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x61, 0x43, 0x65, 0x72, 0x74, 0x42, 0x75, 0x6e, 0x64, 0x6c,
	0x65, 0x12, 0x3a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x27, 0x0a,
	0x0f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61,
	0x70, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65,
	0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_projections_cluster_proto_rawDescOnce sync.Once
	file_api_domain_projections_cluster_proto_rawDescData = file_api_domain_projections_cluster_proto_rawDesc
)

func file_api_domain_projections_cluster_proto_rawDescGZIP() []byte {
	file_api_domain_projections_cluster_proto_rawDescOnce.Do(func() {
		file_api_domain_projections_cluster_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_projections_cluster_proto_rawDescData)
	})
	return file_api_domain_projections_cluster_proto_rawDescData
}

var file_api_domain_projections_cluster_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_domain_projections_cluster_proto_goTypes = []interface{}{
	(*Cluster)(nil),           // 0: projections.Cluster
	(*LifecycleMetadata)(nil), // 1: projections.LifecycleMetadata
}
var file_api_domain_projections_cluster_proto_depIdxs = []int32{
	1, // 0: projections.Cluster.metadata:type_name -> projections.LifecycleMetadata
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_domain_projections_cluster_proto_init() }
func file_api_domain_projections_cluster_proto_init() {
	if File_api_domain_projections_cluster_proto != nil {
		return
	}
	file_api_domain_projections_metadata_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_domain_projections_cluster_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cluster); i {
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
			RawDescriptor: file_api_domain_projections_cluster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_domain_projections_cluster_proto_goTypes,
		DependencyIndexes: file_api_domain_projections_cluster_proto_depIdxs,
		MessageInfos:      file_api_domain_projections_cluster_proto_msgTypes,
	}.Build()
	File_api_domain_projections_cluster_proto = out.File
	file_api_domain_projections_cluster_proto_rawDesc = nil
	file_api_domain_projections_cluster_proto_goTypes = nil
	file_api_domain_projections_cluster_proto_depIdxs = nil
}
