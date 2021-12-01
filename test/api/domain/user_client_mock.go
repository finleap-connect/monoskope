// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/finleap-connect/monoskope/pkg/api/domain (interfaces: UserClient)

// Package domain is a generated GoMock package.
package domain

import (
	context "context"
	reflect "reflect"

	domain "github.com/finleap-connect/monoskope/pkg/api/domain"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
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