// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vmware-tanzu/octant/pkg/action (interfaces: Alerter)

// Package fake is a generated GoMock package.
package fake

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	action "github.com/vmware-tanzu/octant/pkg/action"
)

// MockAlerter is a mock of Alerter interface
type MockAlerter struct {
	ctrl     *gomock.Controller
	recorder *MockAlerterMockRecorder
}

// MockAlerterMockRecorder is the mock recorder for MockAlerter
type MockAlerterMockRecorder struct {
	mock *MockAlerter
}

// NewMockAlerter creates a new mock instance
func NewMockAlerter(ctrl *gomock.Controller) *MockAlerter {
	mock := &MockAlerter{ctrl: ctrl}
	mock.recorder = &MockAlerterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAlerter) EXPECT() *MockAlerterMockRecorder {
	return m.recorder
}

// SendAlert mocks base method
func (m *MockAlerter) SendAlert(arg0 action.Alert) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendAlert", arg0)
}

// SendAlert indicates an expected call of SendAlert
func (mr *MockAlerterMockRecorder) SendAlert(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAlert", reflect.TypeOf((*MockAlerter)(nil).SendAlert), arg0)
}
