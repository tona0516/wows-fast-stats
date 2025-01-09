// Code generated by MockGen. DO NOT EDIT.
// Source: local_file_interface.go
//
// Generated by this command:
//
//	mockgen -source=local_file_interface.go -destination ../mock/repository/local_file_interface.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"
	data "wfs/backend/data"

	gomock "go.uber.org/mock/gomock"
)

// MockLocalFileInterface is a mock of LocalFileInterface interface.
type MockLocalFileInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLocalFileInterfaceMockRecorder
	isgomock struct{}
}

// MockLocalFileInterfaceMockRecorder is the mock recorder for MockLocalFileInterface.
type MockLocalFileInterfaceMockRecorder struct {
	mock *MockLocalFileInterface
}

// NewMockLocalFileInterface creates a new mock instance.
func NewMockLocalFileInterface(ctrl *gomock.Controller) *MockLocalFileInterface {
	mock := &MockLocalFileInterface{ctrl: ctrl}
	mock.recorder = &MockLocalFileInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocalFileInterface) EXPECT() *MockLocalFileInterfaceMockRecorder {
	return m.recorder
}

// SaveScreenshot mocks base method.
func (m *MockLocalFileInterface) SaveScreenshot(path, base64Data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveScreenshot", path, base64Data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveScreenshot indicates an expected call of SaveScreenshot.
func (mr *MockLocalFileInterfaceMockRecorder) SaveScreenshot(path, base64Data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveScreenshot", reflect.TypeOf((*MockLocalFileInterface)(nil).SaveScreenshot), path, base64Data)
}

// SaveTempArenaInfo mocks base method.
func (m *MockLocalFileInterface) SaveTempArenaInfo(tempArenaInfo data.TempArenaInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTempArenaInfo", tempArenaInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTempArenaInfo indicates an expected call of SaveTempArenaInfo.
func (mr *MockLocalFileInterfaceMockRecorder) SaveTempArenaInfo(tempArenaInfo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTempArenaInfo", reflect.TypeOf((*MockLocalFileInterface)(nil).SaveTempArenaInfo), tempArenaInfo)
}

// TempArenaInfo mocks base method.
func (m *MockLocalFileInterface) TempArenaInfo(installPath string) (data.TempArenaInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TempArenaInfo", installPath)
	ret0, _ := ret[0].(data.TempArenaInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TempArenaInfo indicates an expected call of TempArenaInfo.
func (mr *MockLocalFileInterfaceMockRecorder) TempArenaInfo(installPath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TempArenaInfo", reflect.TypeOf((*MockLocalFileInterface)(nil).TempArenaInfo), installPath)
}
