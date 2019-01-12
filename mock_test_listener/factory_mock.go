// Code generated by MockGen. DO NOT EDIT.
// Source: test_listener/factory.go

// Package mock_test_listener is a generated GoMock package.
package mock_test_listener

import (
	gomock "github.com/golang/mock/gomock"
	devices "github.com/ybonjour/atr/devices"
	test_listener "github.com/ybonjour/atr/test_listener"
	reflect "reflect"
)

// MockFactory is a mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// ForDevice mocks base method
func (m *MockFactory) ForDevice(device devices.Device) []test_listener.TestListener {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForDevice", device)
	ret0, _ := ret[0].([]test_listener.TestListener)
	return ret0
}

// ForDevice indicates an expected call of ForDevice
func (mr *MockFactoryMockRecorder) ForDevice(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForDevice", reflect.TypeOf((*MockFactory)(nil).ForDevice), device)
}
