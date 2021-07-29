// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vmware-tanzu/octant/pkg/plugin/service (interfaces: Dashboard)

// Package fake is a generated GoMock package.
package fake

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	action "github.com/vmware-tanzu/octant/pkg/action"
	event "github.com/vmware-tanzu/octant/pkg/event"
	api "github.com/vmware-tanzu/octant/pkg/plugin/api"
	store "github.com/vmware-tanzu/octant/pkg/store"
)

// MockDashboard is a mock of Dashboard interface
type MockDashboard struct {
	ctrl     *gomock.Controller
	recorder *MockDashboardMockRecorder
}

// MockDashboardMockRecorder is the mock recorder for MockDashboard
type MockDashboardMockRecorder struct {
	mock *MockDashboard
}

// NewMockDashboard creates a new mock instance
func NewMockDashboard(ctrl *gomock.Controller) *MockDashboard {
	mock := &MockDashboard{ctrl: ctrl}
	mock.recorder = &MockDashboardMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDashboard) EXPECT() *MockDashboardMockRecorder {
	return m.recorder
}

// CancelPortForward mocks base method
func (m *MockDashboard) CancelPortForward(arg0 context.Context, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CancelPortForward", arg0, arg1)
}

// CancelPortForward indicates an expected call of CancelPortForward
func (mr *MockDashboardMockRecorder) CancelPortForward(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelPortForward", reflect.TypeOf((*MockDashboard)(nil).CancelPortForward), arg0, arg1)
}

// Close mocks base method
func (m *MockDashboard) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockDashboardMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDashboard)(nil).Close))
}

// Create mocks base method
func (m *MockDashboard) Create(arg0 context.Context, arg1 *unstructured.Unstructured) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockDashboardMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDashboard)(nil).Create), arg0, arg1)
}

// CreateLink mocks base method
func (m *MockDashboard) CreateLink(arg0 context.Context, arg1 store.Key) (api.LinkResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLink", arg0, arg1)
	ret0, _ := ret[0].(api.LinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLink indicates an expected call of CreateLink
func (mr *MockDashboardMockRecorder) CreateLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLink", reflect.TypeOf((*MockDashboard)(nil).CreateLink), arg0, arg1)
}

// Delete mocks base method
func (m *MockDashboard) Delete(arg0 context.Context, arg1 store.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockDashboardMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDashboard)(nil).Delete), arg0, arg1)
}

// ForceFrontendUpdate mocks base method
func (m *MockDashboard) ForceFrontendUpdate(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForceFrontendUpdate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForceFrontendUpdate indicates an expected call of ForceFrontendUpdate
func (mr *MockDashboardMockRecorder) ForceFrontendUpdate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForceFrontendUpdate", reflect.TypeOf((*MockDashboard)(nil).ForceFrontendUpdate), arg0)
}

// Get mocks base method
func (m *MockDashboard) Get(arg0 context.Context, arg1 store.Key) (*unstructured.Unstructured, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*unstructured.Unstructured)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockDashboardMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDashboard)(nil).Get), arg0, arg1)
}

// List mocks base method
func (m *MockDashboard) List(arg0 context.Context, arg1 store.Key) (*unstructured.UnstructuredList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*unstructured.UnstructuredList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockDashboardMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDashboard)(nil).List), arg0, arg1)
}

// ListNamespaces mocks base method
func (m *MockDashboard) ListNamespaces(arg0 context.Context) (api.NamespacesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNamespaces", arg0)
	ret0, _ := ret[0].(api.NamespacesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNamespaces indicates an expected call of ListNamespaces
func (mr *MockDashboardMockRecorder) ListNamespaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNamespaces", reflect.TypeOf((*MockDashboard)(nil).ListNamespaces), arg0)
}

// PortForward mocks base method
func (m *MockDashboard) PortForward(arg0 context.Context, arg1 api.PortForwardRequest) (api.PortForwardResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PortForward", arg0, arg1)
	ret0, _ := ret[0].(api.PortForwardResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PortForward indicates an expected call of PortForward
func (mr *MockDashboardMockRecorder) PortForward(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PortForward", reflect.TypeOf((*MockDashboard)(nil).PortForward), arg0, arg1)
}

// SendAlert mocks base method
func (m *MockDashboard) SendAlert(arg0 context.Context, arg1 string, arg2 action.Alert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAlert", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAlert indicates an expected call of SendAlert
func (mr *MockDashboardMockRecorder) SendAlert(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAlert", reflect.TypeOf((*MockDashboard)(nil).SendAlert), arg0, arg1, arg2)
}

// SendEvent mocks base method
func (m *MockDashboard) SendEvent(arg0 context.Context, arg1 string, arg2 event.EventType, arg3 action.Payload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEvent", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEvent indicates an expected call of SendEvent
func (mr *MockDashboardMockRecorder) SendEvent(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEvent", reflect.TypeOf((*MockDashboard)(nil).SendEvent), arg0, arg1, arg2, arg3)
}

// Update mocks base method
func (m *MockDashboard) Update(arg0 context.Context, arg1 *unstructured.Unstructured) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDashboardMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDashboard)(nil).Update), arg0, arg1)
}
