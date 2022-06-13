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

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/finleap-connect/monoskope/pkg/api/domain (interfaces: UserClient,User_GetAllClient)

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	domain "github.com/finleap-connect/monoskope/pkg/api/domain"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// MockUserClient is a mock of UserClient interface.
type MockUserClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserClientMockRecorder
}

// MockUserClientMockRecorder is the mock recorder for MockUserClient.
type MockUserClientMockRecorder struct {
	mock *MockUserClient
}

// NewMockUserClient creates a new mock instance.
func NewMockUserClient(ctrl *gomock.Controller) *MockUserClient {
	mock := &MockUserClient{ctrl: ctrl}
	mock.recorder = &MockUserClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserClient) EXPECT() *MockUserClientMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockUserClient) GetAll(arg0 context.Context, arg1 *domain.GetAllRequest, arg2 ...grpc.CallOption) (domain.User_GetAllClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAll", varargs...)
	ret0, _ := ret[0].(domain.User_GetAllClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUserClientMockRecorder) GetAll(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUserClient)(nil).GetAll), varargs...)
}

// GetByEmail mocks base method.
func (m *MockUserClient) GetByEmail(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (*projections.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByEmail", varargs...)
	ret0, _ := ret[0].(*projections.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUserClientMockRecorder) GetByEmail(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserClient)(nil).GetByEmail), varargs...)
}

// GetById mocks base method.
func (m *MockUserClient) GetById(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (*projections.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetById", varargs...)
	ret0, _ := ret[0].(*projections.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUserClientMockRecorder) GetById(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUserClient)(nil).GetById), varargs...)
}

// GetCount mocks base method.
func (m *MockUserClient) GetCount(arg0 context.Context, arg1 *domain.GetCountRequest, arg2 ...grpc.CallOption) (*domain.GetCountResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCount", varargs...)
	ret0, _ := ret[0].(*domain.GetCountResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockUserClientMockRecorder) GetCount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockUserClient)(nil).GetCount), varargs...)
}

// GetRoleBindingsById mocks base method.
func (m *MockUserClient) GetRoleBindingsById(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (domain.User_GetRoleBindingsByIdClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRoleBindingsById", varargs...)
	ret0, _ := ret[0].(domain.User_GetRoleBindingsByIdClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleBindingsById indicates an expected call of GetRoleBindingsById.
func (mr *MockUserClientMockRecorder) GetRoleBindingsById(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleBindingsById", reflect.TypeOf((*MockUserClient)(nil).GetRoleBindingsById), varargs...)
}

// MockUser_GetAllClient is a mock of User_GetAllClient interface.
type MockUser_GetAllClient struct {
	ctrl     *gomock.Controller
	recorder *MockUser_GetAllClientMockRecorder
}

// MockUser_GetAllClientMockRecorder is the mock recorder for MockUser_GetAllClient.
type MockUser_GetAllClientMockRecorder struct {
	mock *MockUser_GetAllClient
}

// NewMockUser_GetAllClient creates a new mock instance.
func NewMockUser_GetAllClient(ctrl *gomock.Controller) *MockUser_GetAllClient {
	mock := &MockUser_GetAllClient{ctrl: ctrl}
	mock.recorder = &MockUser_GetAllClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser_GetAllClient) EXPECT() *MockUser_GetAllClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockUser_GetAllClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockUser_GetAllClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockUser_GetAllClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockUser_GetAllClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockUser_GetAllClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockUser_GetAllClient)(nil).Context))
}

// Header mocks base method.
func (m *MockUser_GetAllClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockUser_GetAllClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockUser_GetAllClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockUser_GetAllClient) Recv() (*projections.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*projections.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockUser_GetAllClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockUser_GetAllClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockUser_GetAllClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockUser_GetAllClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockUser_GetAllClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockUser_GetAllClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockUser_GetAllClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockUser_GetAllClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockUser_GetAllClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockUser_GetAllClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockUser_GetAllClient)(nil).Trailer))
}