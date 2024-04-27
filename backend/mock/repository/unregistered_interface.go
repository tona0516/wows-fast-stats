// Code generated by MockGen. DO NOT EDIT.
// Source: unregistered_interface.go
//
// Generated by this command:
//
//	mockgen -source=unregistered_interface.go -destination ../mock/repository/unregistered_interface.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"
	data "wfs/backend/data"

	gomock "go.uber.org/mock/gomock"
)

// MockUnregisteredInterface is a mock of UnregisteredInterface interface.
type MockUnregisteredInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUnregisteredInterfaceMockRecorder
}

// MockUnregisteredInterfaceMockRecorder is the mock recorder for MockUnregisteredInterface.
type MockUnregisteredInterfaceMockRecorder struct {
	mock *MockUnregisteredInterface
}

// NewMockUnregisteredInterface creates a new mock instance.
func NewMockUnregisteredInterface(ctrl *gomock.Controller) *MockUnregisteredInterface {
	mock := &MockUnregisteredInterface{ctrl: ctrl}
	mock.recorder = &MockUnregisteredInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnregisteredInterface) EXPECT() *MockUnregisteredInterfaceMockRecorder {
	return m.recorder
}

// Warship mocks base method.
func (m *MockUnregisteredInterface) Warship() (data.Warships, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Warship")
	ret0, _ := ret[0].(data.Warships)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Warship indicates an expected call of Warship.
func (mr *MockUnregisteredInterfaceMockRecorder) Warship() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warship", reflect.TypeOf((*MockUnregisteredInterface)(nil).Warship))
}