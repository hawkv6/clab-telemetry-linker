// Code generated by MockGen. DO NOT EDIT.
// Source: helpers.go
//
// Generated by this command:
//
//	mockgen -source=helpers.go -destination=helpers_mock.go -package=helpers
//

// Package helpers is a generated GoMock package.
package helpers

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHelper is a mock of Helper interface.
type MockHelper struct {
	ctrl     *gomock.Controller
	recorder *MockHelperMockRecorder
}

// MockHelperMockRecorder is the mock recorder for MockHelper.
type MockHelperMockRecorder struct {
	mock *MockHelper
}

// NewMockHelper creates a new mock instance.
func NewMockHelper(ctrl *gomock.Controller) *MockHelper {
	mock := &MockHelper{ctrl: ctrl}
	mock.recorder = &MockHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelper) EXPECT() *MockHelperMockRecorder {
	return m.recorder
}

// GetDefaultClabName mocks base method.
func (m *MockHelper) GetDefaultClabName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultClabName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDefaultClabName indicates an expected call of GetDefaultClabName.
func (mr *MockHelperMockRecorder) GetDefaultClabName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultClabName", reflect.TypeOf((*MockHelper)(nil).GetDefaultClabName))
}

// GetDefaultClabNameKey mocks base method.
func (m *MockHelper) GetDefaultClabNameKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultClabNameKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDefaultClabNameKey indicates an expected call of GetDefaultClabNameKey.
func (mr *MockHelperMockRecorder) GetDefaultClabNameKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultClabNameKey", reflect.TypeOf((*MockHelper)(nil).GetDefaultClabNameKey))
}

// GetDefaultImpairmentsPrefix mocks base method.
func (m *MockHelper) GetDefaultImpairmentsPrefix(node, interface_ string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultImpairmentsPrefix", node, interface_)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDefaultImpairmentsPrefix indicates an expected call of GetDefaultImpairmentsPrefix.
func (mr *MockHelperMockRecorder) GetDefaultImpairmentsPrefix(node, interface_ any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultImpairmentsPrefix", reflect.TypeOf((*MockHelper)(nil).GetDefaultImpairmentsPrefix), node, interface_)
}

// GetUserHome mocks base method.
func (m *MockHelper) GetUserHome() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserHome")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserHome indicates an expected call of GetUserHome.
func (mr *MockHelperMockRecorder) GetUserHome() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserHome", reflect.TypeOf((*MockHelper)(nil).GetUserHome))
}

// IsRoot mocks base method.
func (m *MockHelper) IsRoot() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRoot")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRoot indicates an expected call of IsRoot.
func (mr *MockHelperMockRecorder) IsRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRoot", reflect.TypeOf((*MockHelper)(nil).IsRoot))
}
