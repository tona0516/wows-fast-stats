// Code generated by MockGen. DO NOT EDIT.
// Source: config_v0_interface.go
//
// Generated by this command:
//
//	mockgen -source=config_v0_interface.go -destination ../mock/repository/config_v0_interface.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"
	data "wfs/backend/data"

	gomock "go.uber.org/mock/gomock"
)

// MockConfigV0Interface is a mock of ConfigV0Interface interface.
type MockConfigV0Interface struct {
	ctrl     *gomock.Controller
	recorder *MockConfigV0InterfaceMockRecorder
}

// MockConfigV0InterfaceMockRecorder is the mock recorder for MockConfigV0Interface.
type MockConfigV0InterfaceMockRecorder struct {
	mock *MockConfigV0Interface
}

// NewMockConfigV0Interface creates a new mock instance.
func NewMockConfigV0Interface(ctrl *gomock.Controller) *MockConfigV0Interface {
	mock := &MockConfigV0Interface{ctrl: ctrl}
	mock.recorder = &MockConfigV0InterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigV0Interface) EXPECT() *MockConfigV0InterfaceMockRecorder {
	return m.recorder
}

// AlertPlayers mocks base method.
func (m *MockConfigV0Interface) AlertPlayers() ([]data.AlertPlayer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AlertPlayers")
	ret0, _ := ret[0].([]data.AlertPlayer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AlertPlayers indicates an expected call of AlertPlayers.
func (mr *MockConfigV0InterfaceMockRecorder) AlertPlayers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlertPlayers", reflect.TypeOf((*MockConfigV0Interface)(nil).AlertPlayers))
}

// DeleteAlertPlayers mocks base method.
func (m *MockConfigV0Interface) DeleteAlertPlayers() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAlertPlayers")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAlertPlayers indicates an expected call of DeleteAlertPlayers.
func (mr *MockConfigV0InterfaceMockRecorder) DeleteAlertPlayers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAlertPlayers", reflect.TypeOf((*MockConfigV0Interface)(nil).DeleteAlertPlayers))
}

// DeleteUser mocks base method.
func (m *MockConfigV0Interface) DeleteUser() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockConfigV0InterfaceMockRecorder) DeleteUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockConfigV0Interface)(nil).DeleteUser))
}

// IsExistAlertPlayers mocks base method.
func (m *MockConfigV0Interface) IsExistAlertPlayers() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExistAlertPlayers")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsExistAlertPlayers indicates an expected call of IsExistAlertPlayers.
func (mr *MockConfigV0InterfaceMockRecorder) IsExistAlertPlayers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExistAlertPlayers", reflect.TypeOf((*MockConfigV0Interface)(nil).IsExistAlertPlayers))
}

// IsExistUser mocks base method.
func (m *MockConfigV0Interface) IsExistUser() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExistUser")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsExistUser indicates an expected call of IsExistUser.
func (mr *MockConfigV0InterfaceMockRecorder) IsExistUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExistUser", reflect.TypeOf((*MockConfigV0Interface)(nil).IsExistUser))
}

// User mocks base method.
func (m *MockConfigV0Interface) User() (data.UserConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User")
	ret0, _ := ret[0].(data.UserConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// User indicates an expected call of User.
func (mr *MockConfigV0InterfaceMockRecorder) User() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockConfigV0Interface)(nil).User))
}