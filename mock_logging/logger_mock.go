// Code generated by MockGen. DO NOT EDIT.
// Source: logging/logger.go

// Package mock_logging is a generated GoMock package.
package mock_logging

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLogger is a mock of Logger interface
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method
func (m *MockLogger) Debug(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Debug", message)
}

// Debug indicates an expected call of Debug
func (mr *MockLoggerMockRecorder) Debug(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), message)
}

// Info mocks base method
func (m *MockLogger) Info(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", message)
}

// Info indicates an expected call of Info
func (mr *MockLoggerMockRecorder) Info(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), message)
}

// Warn mocks base method
func (m *MockLogger) Warn(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warn", message)
}

// Warn indicates an expected call of Warn
func (mr *MockLoggerMockRecorder) Warn(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogger)(nil).Warn), message)
}

// Error mocks base method
func (m *MockLogger) Error(message string, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", message, err)
}

// Error indicates an expected call of Error
func (mr *MockLoggerMockRecorder) Error(message, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), message, err)
}
