// Code generated by MockGen. DO NOT EDIT.
// Source: adb/adb.go

// Package mock_adb is a generated GoMock package.
package mock_adb

import (
	gomock "github.com/golang/mock/gomock"
	command "github.com/ybonjour/atr/command"
	reflect "reflect"
)

// MockAdb is a mock of Adb interface
type MockAdb struct {
	ctrl     *gomock.Controller
	recorder *MockAdbMockRecorder
}

// MockAdbMockRecorder is the mock recorder for MockAdb
type MockAdbMockRecorder struct {
	mock *MockAdb
}

// NewMockAdb creates a new mock instance
func NewMockAdb(ctrl *gomock.Controller) *MockAdb {
	mock := &MockAdb{ctrl: ctrl}
	mock.recorder = &MockAdbMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdb) EXPECT() *MockAdbMockRecorder {
	return m.recorder
}

// ConnectedDevices mocks base method
func (m *MockAdb) ConnectedDevices() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectedDevices")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectedDevices indicates an expected call of ConnectedDevices
func (mr *MockAdbMockRecorder) ConnectedDevices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectedDevices", reflect.TypeOf((*MockAdb)(nil).ConnectedDevices))
}

// Install mocks base method
func (m *MockAdb) Install(apkPath, deviceSerial string) command.ExecutionResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install", apkPath, deviceSerial)
	ret0, _ := ret[0].(command.ExecutionResult)
	return ret0
}

// Install indicates an expected call of Install
func (mr *MockAdbMockRecorder) Install(apkPath, deviceSerial interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockAdb)(nil).Install), apkPath, deviceSerial)
}

// Uninstall mocks base method
func (m *MockAdb) Uninstall(packageName, deviceSerial string) command.ExecutionResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uninstall", packageName, deviceSerial)
	ret0, _ := ret[0].(command.ExecutionResult)
	return ret0
}

// Uninstall indicates an expected call of Uninstall
func (mr *MockAdbMockRecorder) Uninstall(packageName, deviceSerial interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uninstall", reflect.TypeOf((*MockAdb)(nil).Uninstall), packageName, deviceSerial)
}

// ExecuteTest mocks base method
func (m *MockAdb) ExecuteTest(packageName, testRunner, test, deviceSerial string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteTest", packageName, testRunner, test, deviceSerial)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteTest indicates an expected call of ExecuteTest
func (mr *MockAdbMockRecorder) ExecuteTest(packageName, testRunner, test, deviceSerial interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteTest", reflect.TypeOf((*MockAdb)(nil).ExecuteTest), packageName, testRunner, test, deviceSerial)
}

// ClearLogcat mocks base method
func (m *MockAdb) ClearLogcat(deviceSerial string) command.ExecutionResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearLogcat", deviceSerial)
	ret0, _ := ret[0].(command.ExecutionResult)
	return ret0
}

// ClearLogcat indicates an expected call of ClearLogcat
func (mr *MockAdbMockRecorder) ClearLogcat(deviceSerial interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearLogcat", reflect.TypeOf((*MockAdb)(nil).ClearLogcat), deviceSerial)
}

// GetLogcat mocks base method
func (m *MockAdb) GetLogcat(deviceSerial string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogcat", deviceSerial)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogcat indicates an expected call of GetLogcat
func (mr *MockAdbMockRecorder) GetLogcat(deviceSerial interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogcat", reflect.TypeOf((*MockAdb)(nil).GetLogcat), deviceSerial)
}

// RecordScreen mocks base method
func (m *MockAdb) RecordScreen(deviceSerial, filePath string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordScreen", deviceSerial, filePath)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RecordScreen indicates an expected call of RecordScreen
func (mr *MockAdbMockRecorder) RecordScreen(deviceSerial, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordScreen", reflect.TypeOf((*MockAdb)(nil).RecordScreen), deviceSerial, filePath)
}

// PullFile mocks base method
func (m *MockAdb) PullFile(deviceSerial, filePathOnDevice, filePathLocal string) command.ExecutionResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullFile", deviceSerial, filePathOnDevice, filePathLocal)
	ret0, _ := ret[0].(command.ExecutionResult)
	return ret0
}

// PullFile indicates an expected call of PullFile
func (mr *MockAdbMockRecorder) PullFile(deviceSerial, filePathOnDevice, filePathLocal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullFile", reflect.TypeOf((*MockAdb)(nil).PullFile), deviceSerial, filePathOnDevice, filePathLocal)
}

// RemoveFile mocks base method
func (m *MockAdb) RemoveFile(deviceSerial, filePathOnDevice string) command.ExecutionResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFile", deviceSerial, filePathOnDevice)
	ret0, _ := ret[0].(command.ExecutionResult)
	return ret0
}

// RemoveFile indicates an expected call of RemoveFile
func (mr *MockAdbMockRecorder) RemoveFile(deviceSerial, filePathOnDevice interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFile", reflect.TypeOf((*MockAdb)(nil).RemoveFile), deviceSerial, filePathOnDevice)
}
