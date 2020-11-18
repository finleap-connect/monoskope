// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/gateway/auth/messages.proto

package auth

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type AuthState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CallbackURL string `protobuf:"bytes,1,opt,name=CallbackURL,proto3" json:"CallbackURL,omitempty"`
}

func (x *AuthState) Reset() {
	*x = AuthState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_auth_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthState) ProtoMessage() {}

func (x *AuthState) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_auth_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthState.ProtoReflect.Descriptor instead.
func (*AuthState) Descriptor() ([]byte, []int) {
	return file_api_gateway_auth_messages_proto_rawDescGZIP(), []int{0}
}

func (x *AuthState) GetCallbackURL() string {
	if x != nil {
		return x.CallbackURL
	}
	return ""
}

type AuthInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthCodeURL string `protobuf:"bytes,1,opt,name=AuthCodeURL,proto3" json:"AuthCodeURL,omitempty"`
	State       string `protobuf:"bytes,2,opt,name=State,proto3" json:"State,omitempty"`
}

func (x *AuthInformation) Reset() {
	*x = AuthInformation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_auth_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthInformation) ProtoMessage() {}

func (x *AuthInformation) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_auth_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthInformation.ProtoReflect.Descriptor instead.
func (*AuthInformation) Descriptor() ([]byte, []int) {
	return file_api_gateway_auth_messages_proto_rawDescGZIP(), []int{1}
}

func (x *AuthInformation) GetAuthCodeURL() string {
	if x != nil {
		return x.AuthCodeURL
	}
	return ""
}

func (x *AuthInformation) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type AuthCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code        string `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty"`
	State       string `protobuf:"bytes,2,opt,name=State,proto3" json:"State,omitempty"`
	CallbackURL string `protobuf:"bytes,3,opt,name=CallbackURL,proto3" json:"CallbackURL,omitempty"`
}

func (x *AuthCode) Reset() {
	*x = AuthCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_auth_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthCode) ProtoMessage() {}

func (x *AuthCode) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_auth_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthCode.ProtoReflect.Descriptor instead.
func (*AuthCode) Descriptor() ([]byte, []int) {
	return file_api_gateway_auth_messages_proto_rawDescGZIP(), []int{2}
}

func (x *AuthCode) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *AuthCode) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *AuthCode) GetCallbackURL() string {
	if x != nil {
		return x.CallbackURL
	}
	return ""
}

type UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthResponse *AuthResponse `protobuf:"bytes,1,opt,name=AuthResponse,proto3" json:"AuthResponse,omitempty"`
	Email        string        `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
}

func (x *UserInfo) Reset() {
	*x = UserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_auth_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfo) ProtoMessage() {}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_auth_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfo.ProtoReflect.Descriptor instead.
func (*UserInfo) Descriptor() ([]byte, []int) {
	return file_api_gateway_auth_messages_proto_rawDescGZIP(), []int{3}
}

func (x *UserInfo) GetAuthResponse() *AuthResponse {
	if x != nil {
		return x.AuthResponse
	}
	return nil
}

func (x *UserInfo) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken  *AccessToken `protobuf:"bytes,1,opt,name=AccessToken,proto3" json:"AccessToken,omitempty"`
	RefreshToken string       `protobuf:"bytes,2,opt,name=RefreshToken,proto3" json:"RefreshToken,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_auth_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResponse) ProtoMessage() {}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_auth_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResponse.ProtoReflect.Descriptor instead.
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return file_api_gateway_auth_messages_proto_rawDescGZIP(), []int{4}
}

func (x *AuthResponse) GetAccessToken() *AccessToken {
	if x != nil {
		return x.AccessToken
	}
	return nil
}

func (x *AuthResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type AccessToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token  string               `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
	Expiry *timestamp.Timestamp `protobuf:"bytes,2,opt,name=Expiry,proto3" json:"Expiry,omitempty"`
}

func (x *AccessToken) Reset() {
	*x = AccessToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_auth_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessToken) ProtoMessage() {}

func (x *AccessToken) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_auth_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessToken.ProtoReflect.Descriptor instead.
func (*AccessToken) Descriptor() ([]byte, []int) {
	return file_api_gateway_auth_messages_proto_rawDescGZIP(), []int{5}
}

func (x *AccessToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *AccessToken) GetExpiry() *timestamp.Timestamp {
	if x != nil {
		return x.Expiry
	}
	return nil
}

var File_api_gateway_auth_messages_proto protoreflect.FileDescriptor

var file_api_gateway_auth_messages_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2d, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55, 0x52, 0x4c, 0x22,
	0x49, 0x0a, 0x0f, 0x41, 0x75, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x55, 0x52,
	0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64,
	0x65, 0x55, 0x52, 0x4c, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x56, 0x0a, 0x08, 0x41, 0x75,
	0x74, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55, 0x52, 0x4c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55,
	0x52, 0x4c, 0x22, 0x60, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3e,
	0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x52, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x22, 0x6f, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x57, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x32, 0x0a, 0x06, 0x45, 0x78,
	0x70, 0x69, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x45, 0x78, 0x70, 0x69, 0x72, 0x79, 0x42, 0x49,
	0x5a, 0x41, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d,
	0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f,
	0x70, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0xa2, 0x02, 0x03, 0x4d, 0x47, 0x41, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_gateway_auth_messages_proto_rawDescOnce sync.Once
	file_api_gateway_auth_messages_proto_rawDescData = file_api_gateway_auth_messages_proto_rawDesc
)

func file_api_gateway_auth_messages_proto_rawDescGZIP() []byte {
	file_api_gateway_auth_messages_proto_rawDescOnce.Do(func() {
		file_api_gateway_auth_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_gateway_auth_messages_proto_rawDescData)
	})
	return file_api_gateway_auth_messages_proto_rawDescData
}

var file_api_gateway_auth_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_gateway_auth_messages_proto_goTypes = []interface{}{
	(*AuthState)(nil),           // 0: gateway.auth.AuthState
	(*AuthInformation)(nil),     // 1: gateway.auth.AuthInformation
	(*AuthCode)(nil),            // 2: gateway.auth.AuthCode
	(*UserInfo)(nil),            // 3: gateway.auth.UserInfo
	(*AuthResponse)(nil),        // 4: gateway.auth.AuthResponse
	(*AccessToken)(nil),         // 5: gateway.auth.AccessToken
	(*timestamp.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_api_gateway_auth_messages_proto_depIdxs = []int32{
	4, // 0: gateway.auth.UserInfo.AuthResponse:type_name -> gateway.auth.AuthResponse
	5, // 1: gateway.auth.AuthResponse.AccessToken:type_name -> gateway.auth.AccessToken
	6, // 2: gateway.auth.AccessToken.Expiry:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_gateway_auth_messages_proto_init() }
func file_api_gateway_auth_messages_proto_init() {
	if File_api_gateway_auth_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_gateway_auth_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthState); i {
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
		file_api_gateway_auth_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthInformation); i {
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
		file_api_gateway_auth_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthCode); i {
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
		file_api_gateway_auth_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfo); i {
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
		file_api_gateway_auth_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthResponse); i {
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
		file_api_gateway_auth_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessToken); i {
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
			RawDescriptor: file_api_gateway_auth_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_gateway_auth_messages_proto_goTypes,
		DependencyIndexes: file_api_gateway_auth_messages_proto_depIdxs,
		MessageInfos:      file_api_gateway_auth_messages_proto_msgTypes,
	}.Build()
	File_api_gateway_auth_messages_proto = out.File
	file_api_gateway_auth_messages_proto_rawDesc = nil
	file_api_gateway_auth_messages_proto_goTypes = nil
	file_api_gateway_auth_messages_proto_depIdxs = nil
}
