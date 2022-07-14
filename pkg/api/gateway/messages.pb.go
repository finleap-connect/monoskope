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
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: api/gateway/messages.proto

package gateway

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// AuthorizationScope is an enum defining the available API scopes.
type AuthorizationScope int32

const (
	AuthorizationScope_NONE              AuthorizationScope = 0 // Dummy to prevent accidents
	AuthorizationScope_API               AuthorizationScope = 1 // Read-write for the complete API
	AuthorizationScope_WRITE_SCIM        AuthorizationScope = 2 // Read-write for endpoints under "/scim"
	AuthorizationScope_WRITE_K8SOPERATOR AuthorizationScope = 3 // Read-write for K8sOperator endpoints
)

// Enum value maps for AuthorizationScope.
var (
	AuthorizationScope_name = map[int32]string{
		0: "NONE",
		1: "API",
		2: "WRITE_SCIM",
		3: "WRITE_K8SOPERATOR",
	}
	AuthorizationScope_value = map[string]int32{
		"NONE":              0,
		"API":               1,
		"WRITE_SCIM":        2,
		"WRITE_K8SOPERATOR": 3,
	}
)

func (x AuthorizationScope) Enum() *AuthorizationScope {
	p := new(AuthorizationScope)
	*p = x
	return p
}

func (x AuthorizationScope) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthorizationScope) Descriptor() protoreflect.EnumDescriptor {
	return file_api_gateway_messages_proto_enumTypes[0].Descriptor()
}

func (AuthorizationScope) Type() protoreflect.EnumType {
	return &file_api_gateway_messages_proto_enumTypes[0]
}

func (x AuthorizationScope) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthorizationScope.Descriptor instead.
func (AuthorizationScope) EnumDescriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{0}
}

type UpstreamAuthenticationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// callback_url is the URL where the authorization code
	// will be redirected to by the upstream IDP
	CallbackUrl string `protobuf:"bytes,1,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
}

func (x *UpstreamAuthenticationRequest) Reset() {
	*x = UpstreamAuthenticationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpstreamAuthenticationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpstreamAuthenticationRequest) ProtoMessage() {}

func (x *UpstreamAuthenticationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpstreamAuthenticationRequest.ProtoReflect.Descriptor instead.
func (*UpstreamAuthenticationRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{0}
}

func (x *UpstreamAuthenticationRequest) GetCallbackUrl() string {
	if x != nil {
		return x.CallbackUrl
	}
	return ""
}

type UpstreamAuthenticationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// upstream_idp_redirect is the URL of the IDP to authenticate against
	UpstreamIdpRedirect string `protobuf:"bytes,1,opt,name=upstream_idp_redirect,json=upstreamIdpRedirect,proto3" json:"upstream_idp_redirect,omitempty"`
	// state is the encoded, server-side nonced state containing the callback.
	// This has to be presented to the server along with the actual m8
	// AuthenticationRequest.
	State string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *UpstreamAuthenticationResponse) Reset() {
	*x = UpstreamAuthenticationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpstreamAuthenticationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpstreamAuthenticationResponse) ProtoMessage() {}

func (x *UpstreamAuthenticationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpstreamAuthenticationResponse.ProtoReflect.Descriptor instead.
func (*UpstreamAuthenticationResponse) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{1}
}

func (x *UpstreamAuthenticationResponse) GetUpstreamIdpRedirect() string {
	if x != nil {
		return x.UpstreamIdpRedirect
	}
	return ""
}

func (x *UpstreamAuthenticationResponse) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type AuthenticationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// code is the auth code received by the IDP
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	// state is the encoded, nonced AuthState
	State string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *AuthenticationRequest) Reset() {
	*x = AuthenticationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticationRequest) ProtoMessage() {}

func (x *AuthenticationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticationRequest.ProtoReflect.Descriptor instead.
func (*AuthenticationRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{2}
}

func (x *AuthenticationRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *AuthenticationRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type AuthenticationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// access_token is a JWT to authenticate against the m8 API
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// expiry is the timestamp when the token expires
	Expiry *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiry,proto3" json:"expiry,omitempty"`
	// username is the username known the m8 control plane
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *AuthenticationResponse) Reset() {
	*x = AuthenticationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticationResponse) ProtoMessage() {}

func (x *AuthenticationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticationResponse.ProtoReflect.Descriptor instead.
func (*AuthenticationResponse) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{3}
}

func (x *AuthenticationResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AuthenticationResponse) GetExpiry() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiry
	}
	return nil
}

func (x *AuthenticationResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// ClusterAuthTokenRequest is send in order to retrieve an auth token valid to
// authenticate against a certain cluster with a specific role.
type ClusterAuthTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the cluster (UUID 128-bit number)
	ClusterId string `protobuf:"bytes,1,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	// Kubernetes role name
	Role string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *ClusterAuthTokenRequest) Reset() {
	*x = ClusterAuthTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterAuthTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterAuthTokenRequest) ProtoMessage() {}

func (x *ClusterAuthTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterAuthTokenRequest.ProtoReflect.Descriptor instead.
func (*ClusterAuthTokenRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{4}
}

func (x *ClusterAuthTokenRequest) GetClusterId() string {
	if x != nil {
		return x.ClusterId
	}
	return ""
}

func (x *ClusterAuthTokenRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

// ClusterAuthTokenResponse contains an auth token valid to
// authenticate against a certain cluster with a specific role.
type ClusterAuthTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// JWT to authenticate against a K8s cluster
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// Timestamp when the token expires
	Expiry *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiry,proto3" json:"expiry,omitempty"`
}

func (x *ClusterAuthTokenResponse) Reset() {
	*x = ClusterAuthTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterAuthTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterAuthTokenResponse) ProtoMessage() {}

func (x *ClusterAuthTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterAuthTokenResponse.ProtoReflect.Descriptor instead.
func (*ClusterAuthTokenResponse) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{5}
}

func (x *ClusterAuthTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *ClusterAuthTokenResponse) GetExpiry() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiry
	}
	return nil
}

// APITokenRequest is send in order to retrieve an API token valid to
// authenticate against Monoskope and authorize specific scopes.
type APITokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Scope the resulting token is issued for
	AuthorizationScopes []AuthorizationScope `protobuf:"varint,1,rep,packed,name=authorization_scopes,json=authorizationScopes,proto3,enum=gateway.AuthorizationScope" json:"authorization_scopes,omitempty"`
	// Duration for which the issued token will be valid
	Validity *durationpb.Duration `protobuf:"bytes,2,opt,name=validity,proto3" json:"validity,omitempty"`
	// Types that are assignable to User:
	//	*APITokenRequest_UserId
	//	*APITokenRequest_Username
	User isAPITokenRequest_User `protobuf_oneof:"user"`
}

func (x *APITokenRequest) Reset() {
	*x = APITokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APITokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APITokenRequest) ProtoMessage() {}

func (x *APITokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APITokenRequest.ProtoReflect.Descriptor instead.
func (*APITokenRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{6}
}

func (x *APITokenRequest) GetAuthorizationScopes() []AuthorizationScope {
	if x != nil {
		return x.AuthorizationScopes
	}
	return nil
}

func (x *APITokenRequest) GetValidity() *durationpb.Duration {
	if x != nil {
		return x.Validity
	}
	return nil
}

func (m *APITokenRequest) GetUser() isAPITokenRequest_User {
	if m != nil {
		return m.User
	}
	return nil
}

func (x *APITokenRequest) GetUserId() string {
	if x, ok := x.GetUser().(*APITokenRequest_UserId); ok {
		return x.UserId
	}
	return ""
}

func (x *APITokenRequest) GetUsername() string {
	if x, ok := x.GetUser().(*APITokenRequest_Username); ok {
		return x.Username
	}
	return ""
}

type isAPITokenRequest_User interface {
	isAPITokenRequest_User()
}

type APITokenRequest_UserId struct {
	// Unique identifier of an existing user (UUID 128-bit number)
	UserId string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3,oneof"`
}

type APITokenRequest_Username struct {
	// Name of the user the token is valid for (not necessarily a real user)
	Username string `protobuf:"bytes,4,opt,name=username,proto3,oneof"`
}

func (*APITokenRequest_UserId) isAPITokenRequest_User() {}

func (*APITokenRequest_Username) isAPITokenRequest_User() {}

// APITokenResponse is the answer to an APITokenRequest
// containing a JWT to authenticate against the m8 API.
type APITokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// JWT to authenticate against the m8 API
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// Timestamp when the token expires
	Expiry *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiry,proto3" json:"expiry,omitempty"`
}

func (x *APITokenResponse) Reset() {
	*x = APITokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APITokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APITokenResponse) ProtoMessage() {}

func (x *APITokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APITokenResponse.ProtoReflect.Descriptor instead.
func (*APITokenResponse) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{7}
}

func (x *APITokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *APITokenResponse) GetExpiry() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiry
	}
	return nil
}

// Request information that should be checked if authorized.
type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// FullMethodName is the full RPC method string, i.e.,
	// /package.service/method.
	FullMethodName string `protobuf:"bytes,1,opt,name=full_method_name,json=fullMethodName,proto3" json:"full_method_name,omitempty"`
	// AccessToken is the token from the auth header of the client request
	AccessToken string `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// The actual request to authorize
	Request []byte `protobuf:"bytes,3,opt,name=request,proto3" json:"request,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{8}
}

func (x *CheckRequest) GetFullMethodName() string {
	if x != nil {
		return x.FullMethodName
	}
	return ""
}

func (x *CheckRequest) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *CheckRequest) GetRequest() []byte {
	if x != nil {
		return x.Request
	}
	return nil
}

// Intended for gRPC and Network Authorization servers `only`.
// Status `OK` allows the request. Any other status indicates the request
// should be denied.
type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tags []*CheckResponse_CheckResponseTag `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{9}
}

func (x *CheckResponse) GetTags() []*CheckResponse_CheckResponseTag {
	if x != nil {
		return x.Tags
	}
	return nil
}

type CheckResponse_CheckResponseTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *CheckResponse_CheckResponseTag) Reset() {
	*x = CheckResponse_CheckResponseTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_messages_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse_CheckResponseTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse_CheckResponseTag) ProtoMessage() {}

func (x *CheckResponse_CheckResponseTag) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_messages_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse_CheckResponseTag.ProtoReflect.Descriptor instead.
func (*CheckResponse_CheckResponseTag) Descriptor() ([]byte, []int) {
	return file_api_gateway_messages_proto_rawDescGZIP(), []int{9, 0}
}

func (x *CheckResponse_CheckResponseTag) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CheckResponse_CheckResponseTag) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_api_gateway_messages_proto protoreflect.FileDescriptor

var file_api_gateway_messages_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x4c, 0x0a, 0x1d, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x75, 0x74, 0x68, 0x65,
	0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2b, 0x0a, 0x0c, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x88, 0x01, 0x01,
	0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55, 0x72, 0x6c, 0x22, 0x6a, 0x0a,
	0x1e, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x15, 0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x70, 0x5f,
	0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13,
	0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x70, 0x52, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x53, 0x0a, 0x15, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x1d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x8b,
	0x01, 0x0a, 0x16, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x32, 0x0a, 0x06,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79,
	0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6d, 0x0a, 0x17,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x29, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x15,
	0xfa, 0x42, 0x12, 0x72, 0x10, 0x18, 0x3c, 0x32, 0x0c, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d,
	0x39, 0x2d, 0x5d, 0x2b, 0x24, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x71, 0x0a, 0x18, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x32, 0x0a, 0x06, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x22, 0xec,
	0x01, 0x0a, 0x0f, 0x41, 0x50, 0x49, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x4e, 0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e,
	0x32, 0x1b, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x52, 0x13, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x63, 0x6f, 0x70,
	0x65, 0x73, 0x12, 0x35, 0x0a, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x12, 0x23, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0xb0, 0x01, 0x01, 0x48, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x25,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x00, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x69, 0x0a,
	0x10, 0x41, 0x50, 0x49, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x32, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x22, 0x87, 0x01, 0x0a, 0x0c, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x10, 0x66, 0x75, 0x6c,
	0x6c, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0e, 0x66, 0x75,
	0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x0c,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0b, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x88, 0x01, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x27, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x1a, 0x3a, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x54, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x4e, 0x0a,
	0x12, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x63,
	0x6f, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x07, 0x0a,
	0x03, 0x41, 0x50, 0x49, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x57, 0x52, 0x49, 0x54, 0x45, 0x5f,
	0x53, 0x43, 0x49, 0x4d, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x57, 0x52, 0x49, 0x54, 0x45, 0x5f,
	0x4b, 0x38, 0x53, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x4f, 0x52, 0x10, 0x03, 0x42, 0x36, 0x5a,
	0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6e, 0x6c,
	0x65, 0x61, 0x70, 0x2d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f,
	0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_gateway_messages_proto_rawDescOnce sync.Once
	file_api_gateway_messages_proto_rawDescData = file_api_gateway_messages_proto_rawDesc
)

func file_api_gateway_messages_proto_rawDescGZIP() []byte {
	file_api_gateway_messages_proto_rawDescOnce.Do(func() {
		file_api_gateway_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_gateway_messages_proto_rawDescData)
	})
	return file_api_gateway_messages_proto_rawDescData
}

var file_api_gateway_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_gateway_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_gateway_messages_proto_goTypes = []interface{}{
	(AuthorizationScope)(0),                // 0: gateway.AuthorizationScope
	(*UpstreamAuthenticationRequest)(nil),  // 1: gateway.UpstreamAuthenticationRequest
	(*UpstreamAuthenticationResponse)(nil), // 2: gateway.UpstreamAuthenticationResponse
	(*AuthenticationRequest)(nil),          // 3: gateway.AuthenticationRequest
	(*AuthenticationResponse)(nil),         // 4: gateway.AuthenticationResponse
	(*ClusterAuthTokenRequest)(nil),        // 5: gateway.ClusterAuthTokenRequest
	(*ClusterAuthTokenResponse)(nil),       // 6: gateway.ClusterAuthTokenResponse
	(*APITokenRequest)(nil),                // 7: gateway.APITokenRequest
	(*APITokenResponse)(nil),               // 8: gateway.APITokenResponse
	(*CheckRequest)(nil),                   // 9: gateway.CheckRequest
	(*CheckResponse)(nil),                  // 10: gateway.CheckResponse
	(*CheckResponse_CheckResponseTag)(nil), // 11: gateway.CheckResponse.CheckResponseTag
	(*timestamppb.Timestamp)(nil),          // 12: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),            // 13: google.protobuf.Duration
}
var file_api_gateway_messages_proto_depIdxs = []int32{
	12, // 0: gateway.AuthenticationResponse.expiry:type_name -> google.protobuf.Timestamp
	12, // 1: gateway.ClusterAuthTokenResponse.expiry:type_name -> google.protobuf.Timestamp
	0,  // 2: gateway.APITokenRequest.authorization_scopes:type_name -> gateway.AuthorizationScope
	13, // 3: gateway.APITokenRequest.validity:type_name -> google.protobuf.Duration
	12, // 4: gateway.APITokenResponse.expiry:type_name -> google.protobuf.Timestamp
	11, // 5: gateway.CheckResponse.tags:type_name -> gateway.CheckResponse.CheckResponseTag
	6,  // [6:6] is the sub-list for method output_type
	6,  // [6:6] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_api_gateway_messages_proto_init() }
func file_api_gateway_messages_proto_init() {
	if File_api_gateway_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_gateway_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpstreamAuthenticationRequest); i {
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
		file_api_gateway_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpstreamAuthenticationResponse); i {
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
		file_api_gateway_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticationRequest); i {
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
		file_api_gateway_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticationResponse); i {
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
		file_api_gateway_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterAuthTokenRequest); i {
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
		file_api_gateway_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterAuthTokenResponse); i {
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
		file_api_gateway_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APITokenRequest); i {
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
		file_api_gateway_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APITokenResponse); i {
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
		file_api_gateway_messages_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_api_gateway_messages_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
		file_api_gateway_messages_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse_CheckResponseTag); i {
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
	file_api_gateway_messages_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*APITokenRequest_UserId)(nil),
		(*APITokenRequest_Username)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_gateway_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_gateway_messages_proto_goTypes,
		DependencyIndexes: file_api_gateway_messages_proto_depIdxs,
		EnumInfos:         file_api_gateway_messages_proto_enumTypes,
		MessageInfos:      file_api_gateway_messages_proto_msgTypes,
	}.Build()
	File_api_gateway_messages_proto = out.File
	file_api_gateway_messages_proto_rawDesc = nil
	file_api_gateway_messages_proto_goTypes = nil
	file_api_gateway_messages_proto_depIdxs = nil
}
