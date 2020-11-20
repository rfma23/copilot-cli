// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/ecs/ecs.go

// Package mocks is a generated GoMock package.
package mocks

import (
	ecs "github.com/aws/copilot-cli/internal/pkg/aws/ecs"
	resourcegroups "github.com/aws/copilot-cli/internal/pkg/aws/resourcegroups"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockresourceGetter is a mock of resourceGetter interface
type MockresourceGetter struct {
	ctrl     *gomock.Controller
	recorder *MockresourceGetterMockRecorder
}

// MockresourceGetterMockRecorder is the mock recorder for MockresourceGetter
type MockresourceGetterMockRecorder struct {
	mock *MockresourceGetter
}

// NewMockresourceGetter creates a new mock instance
func NewMockresourceGetter(ctrl *gomock.Controller) *MockresourceGetter {
	mock := &MockresourceGetter{ctrl: ctrl}
	mock.recorder = &MockresourceGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockresourceGetter) EXPECT() *MockresourceGetterMockRecorder {
	return m.recorder
}

// GetResourcesByTags mocks base method
func (m *MockresourceGetter) GetResourcesByTags(resourceType string, tags map[string]string) ([]*resourcegroups.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourcesByTags", resourceType, tags)
	ret0, _ := ret[0].([]*resourcegroups.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourcesByTags indicates an expected call of GetResourcesByTags
func (mr *MockresourceGetterMockRecorder) GetResourcesByTags(resourceType, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourcesByTags", reflect.TypeOf((*MockresourceGetter)(nil).GetResourcesByTags), resourceType, tags)
}

// MockRunningTasksInFamilyGetter is a mock of RunningTasksInFamilyGetter interface
type MockRunningTasksInFamilyGetter struct {
	ctrl     *gomock.Controller
	recorder *MockRunningTasksInFamilyGetterMockRecorder
}

// MockRunningTasksInFamilyGetterMockRecorder is the mock recorder for MockRunningTasksInFamilyGetter
type MockRunningTasksInFamilyGetterMockRecorder struct {
	mock *MockRunningTasksInFamilyGetter
}

// NewMockRunningTasksInFamilyGetter creates a new mock instance
func NewMockRunningTasksInFamilyGetter(ctrl *gomock.Controller) *MockRunningTasksInFamilyGetter {
	mock := &MockRunningTasksInFamilyGetter{ctrl: ctrl}
	mock.recorder = &MockRunningTasksInFamilyGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRunningTasksInFamilyGetter) EXPECT() *MockRunningTasksInFamilyGetterMockRecorder {
	return m.recorder
}

// RunningTasksInFamily mocks base method
func (m *MockRunningTasksInFamilyGetter) RunningTasksInFamily(cluster, family string) ([]*ecs.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunningTasksInFamily", cluster, family)
	ret0, _ := ret[0].([]*ecs.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunningTasksInFamily indicates an expected call of RunningTasksInFamily
func (mr *MockRunningTasksInFamilyGetterMockRecorder) RunningTasksInFamily(cluster, family interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunningTasksInFamily", reflect.TypeOf((*MockRunningTasksInFamilyGetter)(nil).RunningTasksInFamily), cluster, family)
}
