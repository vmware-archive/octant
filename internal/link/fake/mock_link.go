// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vmware-tanzu/octant/internal/link (interfaces: Interface,Config)

// Package fake is a generated GoMock package.
package fake

import (
	url "net/url"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"

	component "github.com/vmware-tanzu/octant/pkg/view/component"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// ForGVK mocks base method.
func (m *MockInterface) ForGVK(arg0, arg1, arg2, arg3, arg4 string) (*component.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForGVK", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*component.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForGVK indicates an expected call of ForGVK.
func (mr *MockInterfaceMockRecorder) ForGVK(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForGVK", reflect.TypeOf((*MockInterface)(nil).ForGVK), arg0, arg1, arg2, arg3, arg4)
}

// ForObject mocks base method.
func (m *MockInterface) ForObject(arg0 runtime.Object, arg1 string) (*component.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForObject", arg0, arg1)
	ret0, _ := ret[0].(*component.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForObject indicates an expected call of ForObject.
func (mr *MockInterfaceMockRecorder) ForObject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForObject", reflect.TypeOf((*MockInterface)(nil).ForObject), arg0, arg1)
}

// ForObjectWithQuery mocks base method.
func (m *MockInterface) ForObjectWithQuery(arg0 runtime.Object, arg1 string, arg2 url.Values) (*component.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForObjectWithQuery", arg0, arg1, arg2)
	ret0, _ := ret[0].(*component.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForObjectWithQuery indicates an expected call of ForObjectWithQuery.
func (mr *MockInterfaceMockRecorder) ForObjectWithQuery(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForObjectWithQuery", reflect.TypeOf((*MockInterface)(nil).ForObjectWithQuery), arg0, arg1, arg2)
}

// ForOwner mocks base method.
func (m *MockInterface) ForOwner(arg0 runtime.Object, arg1 *v1.OwnerReference) (*component.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForOwner", arg0, arg1)
	ret0, _ := ret[0].(*component.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForOwner indicates an expected call of ForOwner.
func (mr *MockInterfaceMockRecorder) ForOwner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForOwner", reflect.TypeOf((*MockInterface)(nil).ForOwner), arg0, arg1)
}

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// ObjectPath mocks base method.
func (m *MockConfig) ObjectPath(arg0, arg1, arg2, arg3 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectPath", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ObjectPath indicates an expected call of ObjectPath.
func (mr *MockConfigMockRecorder) ObjectPath(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectPath", reflect.TypeOf((*MockConfig)(nil).ObjectPath), arg0, arg1, arg2, arg3)
}
