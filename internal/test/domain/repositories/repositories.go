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
// Source: github.com/finleap-connect/monoskope/pkg/domain/repositories (interfaces: UserRepository,ClusterRepository,ClusterAccessRepository)

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	context "context"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	projections0 "github.com/finleap-connect/monoskope/pkg/domain/projections"
	eventsourcing "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	reflect "reflect"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockUserRepository) All(arg0 context.Context) ([]*projections0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*projections0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockUserRepositoryMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockUserRepository)(nil).All), arg0)
}

// AllWith mocks base method.
func (m *MockUserRepository) AllWith(arg0 context.Context, arg1 bool) ([]*projections0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWith", arg0, arg1)
	ret0, _ := ret[0].([]*projections0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWith indicates an expected call of AllWith.
func (mr *MockUserRepositoryMockRecorder) AllWith(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWith", reflect.TypeOf((*MockUserRepository)(nil).AllWith), arg0, arg1)
}

// ByEmail mocks base method.
func (m *MockUserRepository) ByEmail(arg0 context.Context, arg1 string) (*projections0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByEmail", arg0, arg1)
	ret0, _ := ret[0].(*projections0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByEmail indicates an expected call of ByEmail.
func (mr *MockUserRepositoryMockRecorder) ByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByEmail", reflect.TypeOf((*MockUserRepository)(nil).ByEmail), arg0, arg1)
}

// ByEmailIncludingDeleted mocks base method.
func (m *MockUserRepository) ByEmailIncludingDeleted(arg0 context.Context, arg1 string) ([]*projections0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByEmailIncludingDeleted", arg0, arg1)
	ret0, _ := ret[0].([]*projections0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByEmailIncludingDeleted indicates an expected call of ByEmailIncludingDeleted.
func (mr *MockUserRepositoryMockRecorder) ByEmailIncludingDeleted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByEmailIncludingDeleted", reflect.TypeOf((*MockUserRepository)(nil).ByEmailIncludingDeleted), arg0, arg1)
}

// ById mocks base method.
func (m *MockUserRepository) ById(arg0 context.Context, arg1 uuid.UUID) (*projections0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ById", arg0, arg1)
	ret0, _ := ret[0].(*projections0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ById indicates an expected call of ById.
func (mr *MockUserRepositoryMockRecorder) ById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ById", reflect.TypeOf((*MockUserRepository)(nil).ById), arg0, arg1)
}

// ByUserId mocks base method.
func (m *MockUserRepository) ByUserId(arg0 context.Context, arg1 uuid.UUID) (*projections0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByUserId", arg0, arg1)
	ret0, _ := ret[0].(*projections0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByUserId indicates an expected call of ByUserId.
func (mr *MockUserRepositoryMockRecorder) ByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByUserId", reflect.TypeOf((*MockUserRepository)(nil).ByUserId), arg0, arg1)
}

// DeregisterObserver mocks base method.
func (m *MockUserRepository) DeregisterObserver(arg0 eventsourcing.RepositoryObserver[*projections0.User]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeregisterObserver", arg0)
}

// DeregisterObserver indicates an expected call of DeregisterObserver.
func (mr *MockUserRepositoryMockRecorder) DeregisterObserver(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterObserver", reflect.TypeOf((*MockUserRepository)(nil).DeregisterObserver), arg0)
}

// GetCount mocks base method.
func (m *MockUserRepository) GetCount(arg0 context.Context, arg1 bool) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockUserRepositoryMockRecorder) GetCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockUserRepository)(nil).GetCount), arg0, arg1)
}

// RegisterObserver mocks base method.
func (m *MockUserRepository) RegisterObserver(arg0 eventsourcing.RepositoryObserver[*projections0.User]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterObserver", arg0)
}

// RegisterObserver indicates an expected call of RegisterObserver.
func (mr *MockUserRepositoryMockRecorder) RegisterObserver(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterObserver", reflect.TypeOf((*MockUserRepository)(nil).RegisterObserver), arg0)
}

// Upsert mocks base method.
func (m *MockUserRepository) Upsert(arg0 context.Context, arg1 *projections0.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockUserRepositoryMockRecorder) Upsert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockUserRepository)(nil).Upsert), arg0, arg1)
}

// MockClusterRepository is a mock of ClusterRepository interface.
type MockClusterRepository struct {
	ctrl     *gomock.Controller
	recorder *MockClusterRepositoryMockRecorder
}

// MockClusterRepositoryMockRecorder is the mock recorder for MockClusterRepository.
type MockClusterRepositoryMockRecorder struct {
	mock *MockClusterRepository
}

// NewMockClusterRepository creates a new mock instance.
func NewMockClusterRepository(ctrl *gomock.Controller) *MockClusterRepository {
	mock := &MockClusterRepository{ctrl: ctrl}
	mock.recorder = &MockClusterRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterRepository) EXPECT() *MockClusterRepositoryMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockClusterRepository) All(arg0 context.Context) ([]*projections0.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*projections0.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockClusterRepositoryMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockClusterRepository)(nil).All), arg0)
}

// AllWith mocks base method.
func (m *MockClusterRepository) AllWith(arg0 context.Context, arg1 bool) ([]*projections0.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWith", arg0, arg1)
	ret0, _ := ret[0].([]*projections0.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWith indicates an expected call of AllWith.
func (mr *MockClusterRepositoryMockRecorder) AllWith(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWith", reflect.TypeOf((*MockClusterRepository)(nil).AllWith), arg0, arg1)
}

// ByClusterName mocks base method.
func (m *MockClusterRepository) ByClusterName(arg0 context.Context, arg1 string) (*projections0.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByClusterName", arg0, arg1)
	ret0, _ := ret[0].(*projections0.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByClusterName indicates an expected call of ByClusterName.
func (mr *MockClusterRepositoryMockRecorder) ByClusterName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByClusterName", reflect.TypeOf((*MockClusterRepository)(nil).ByClusterName), arg0, arg1)
}

// ById mocks base method.
func (m *MockClusterRepository) ById(arg0 context.Context, arg1 uuid.UUID) (*projections0.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ById", arg0, arg1)
	ret0, _ := ret[0].(*projections0.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ById indicates an expected call of ById.
func (mr *MockClusterRepositoryMockRecorder) ById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ById", reflect.TypeOf((*MockClusterRepository)(nil).ById), arg0, arg1)
}

// DeregisterObserver mocks base method.
func (m *MockClusterRepository) DeregisterObserver(arg0 eventsourcing.RepositoryObserver[*projections0.Cluster]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeregisterObserver", arg0)
}

// DeregisterObserver indicates an expected call of DeregisterObserver.
func (mr *MockClusterRepositoryMockRecorder) DeregisterObserver(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterObserver", reflect.TypeOf((*MockClusterRepository)(nil).DeregisterObserver), arg0)
}

// RegisterObserver mocks base method.
func (m *MockClusterRepository) RegisterObserver(arg0 eventsourcing.RepositoryObserver[*projections0.Cluster]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterObserver", arg0)
}

// RegisterObserver indicates an expected call of RegisterObserver.
func (mr *MockClusterRepositoryMockRecorder) RegisterObserver(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterObserver", reflect.TypeOf((*MockClusterRepository)(nil).RegisterObserver), arg0)
}

// Upsert mocks base method.
func (m *MockClusterRepository) Upsert(arg0 context.Context, arg1 *projections0.Cluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockClusterRepositoryMockRecorder) Upsert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockClusterRepository)(nil).Upsert), arg0, arg1)
}

// MockClusterAccessRepository is a mock of ClusterAccessRepository interface.
type MockClusterAccessRepository struct {
	ctrl     *gomock.Controller
	recorder *MockClusterAccessRepositoryMockRecorder
}

// MockClusterAccessRepositoryMockRecorder is the mock recorder for MockClusterAccessRepository.
type MockClusterAccessRepositoryMockRecorder struct {
	mock *MockClusterAccessRepository
}

// NewMockClusterAccessRepository creates a new mock instance.
func NewMockClusterAccessRepository(ctrl *gomock.Controller) *MockClusterAccessRepository {
	mock := &MockClusterAccessRepository{ctrl: ctrl}
	mock.recorder = &MockClusterAccessRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterAccessRepository) EXPECT() *MockClusterAccessRepositoryMockRecorder {
	return m.recorder
}

// GetClustersAccessibleByUserId mocks base method.
func (m *MockClusterAccessRepository) GetClustersAccessibleByUserId(arg0 context.Context, arg1 uuid.UUID) ([]*projections.ClusterAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClustersAccessibleByUserId", arg0, arg1)
	ret0, _ := ret[0].([]*projections.ClusterAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClustersAccessibleByUserId indicates an expected call of GetClustersAccessibleByUserId.
func (mr *MockClusterAccessRepositoryMockRecorder) GetClustersAccessibleByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClustersAccessibleByUserId", reflect.TypeOf((*MockClusterAccessRepository)(nil).GetClustersAccessibleByUserId), arg0, arg1)
}

// GetClustersAccessibleByUserIdV2 mocks base method.
func (m *MockClusterAccessRepository) GetClustersAccessibleByUserIdV2(arg0 context.Context, arg1 uuid.UUID) ([]*projections.ClusterAccessV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClustersAccessibleByUserIdV2", arg0, arg1)
	ret0, _ := ret[0].([]*projections.ClusterAccessV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClustersAccessibleByUserIdV2 indicates an expected call of GetClustersAccessibleByUserIdV2.
func (mr *MockClusterAccessRepositoryMockRecorder) GetClustersAccessibleByUserIdV2(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClustersAccessibleByUserIdV2", reflect.TypeOf((*MockClusterAccessRepository)(nil).GetClustersAccessibleByUserIdV2), arg0, arg1)
}
