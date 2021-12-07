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
// source: api/domain/eventdata/user.proto

package eventdata

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

type UserCreated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Email address of the user
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	// Name of the user
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *UserCreated) Reset() {
	*x = UserCreated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserCreated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserCreated) ProtoMessage() {}

func (x *UserCreated) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserCreated.ProtoReflect.Descriptor instead.
func (*UserCreated) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_user_proto_rawDescGZIP(), []int{0}
}

func (x *UserCreated) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserCreated) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type UserRoleAdded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the user (UUID 128-bit number)
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// Name of the role added to the user
	Role string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	// Scope of the role
	Scope string `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	// Unique identifier of the affected resource (UUID 128-bit number)
	Resource string `protobuf:"bytes,4,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *UserRoleAdded) Reset() {
	*x = UserRoleAdded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserRoleAdded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRoleAdded) ProtoMessage() {}

func (x *UserRoleAdded) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRoleAdded.ProtoReflect.Descriptor instead.
func (*UserRoleAdded) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_user_proto_rawDescGZIP(), []int{1}
}

func (x *UserRoleAdded) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserRoleAdded) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *UserRoleAdded) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *UserRoleAdded) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

type UserUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the user
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *UserUpdated) Reset() {
	*x = UserUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_eventdata_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUpdated) ProtoMessage() {}

func (x *UserUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_eventdata_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUpdated.ProtoReflect.Descriptor instead.
func (*UserUpdated) Descriptor() ([]byte, []int) {
	return file_api_domain_eventdata_user_proto_rawDescGZIP(), []int{2}
}

func (x *UserUpdated) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_api_domain_eventdata_user_proto protoreflect.FileDescriptor

var file_api_domain_eventdata_user_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x22, 0x37, 0x0a, 0x0b,
	0x55, 0x73, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6e, 0x0a, 0x0d, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c,
	0x65, 0x41, 0x64, 0x64, 0x65, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x21, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6e, 0x6c, 0x65, 0x61, 0x70, 0x2d, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_domain_eventdata_user_proto_rawDescOnce sync.Once
	file_api_domain_eventdata_user_proto_rawDescData = file_api_domain_eventdata_user_proto_rawDesc
)

func file_api_domain_eventdata_user_proto_rawDescGZIP() []byte {
	file_api_domain_eventdata_user_proto_rawDescOnce.Do(func() {
		file_api_domain_eventdata_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_eventdata_user_proto_rawDescData)
	})
	return file_api_domain_eventdata_user_proto_rawDescData
}

var file_api_domain_eventdata_user_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_domain_eventdata_user_proto_goTypes = []interface{}{
	(*UserCreated)(nil),   // 0: eventdata.UserCreated
	(*UserRoleAdded)(nil), // 1: eventdata.UserRoleAdded
	(*UserUpdated)(nil),   // 2: eventdata.UserUpdated
}
var file_api_domain_eventdata_user_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_domain_eventdata_user_proto_init() }
func file_api_domain_eventdata_user_proto_init() {
	if File_api_domain_eventdata_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_domain_eventdata_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserCreated); i {
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
		file_api_domain_eventdata_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserRoleAdded); i {
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
		file_api_domain_eventdata_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserUpdated); i {
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
			RawDescriptor: file_api_domain_eventdata_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_domain_eventdata_user_proto_goTypes,
		DependencyIndexes: file_api_domain_eventdata_user_proto_depIdxs,
		MessageInfos:      file_api_domain_eventdata_user_proto_msgTypes,
	}.Build()
	File_api_domain_eventdata_user_proto = out.File
	file_api_domain_eventdata_user_proto_rawDesc = nil
	file_api_domain_eventdata_user_proto_goTypes = nil
	file_api_domain_eventdata_user_proto_depIdxs = nil
}
