// Code generated by MockGen. DO NOT EDIT.
// Source: numbers_interface.go
//
// Generated by this command:
//
//	mockgen -source=numbers_interface.go -destination ../mock/repository/numbers_interface.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"
	data "wfs/backend/data"

	gomock "go.uber.org/mock/gomock"
)

// MockNumbersInterface is a mock of NumbersInterface interface.
type MockNumbersInterface struct {
	ctrl     *gomock.Controller
	recorder *MockNumbersInterfaceMockRecorder
	isgomock struct{}
}

// MockNumbersInterfaceMockRecorder is the mock recorder for MockNumbersInterface.
type MockNumbersInterfaceMockRecorder struct {
	mock *MockNumbersInterface
}

// NewMockNumbersInterface creates a new mock instance.
func NewMockNumbersInterface(ctrl *gomock.Controller) *MockNumbersInterface {
	mock := &MockNumbersInterface{ctrl: ctrl}
	mock.recorder = &MockNumbersInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNumbersInterface) EXPECT() *MockNumbersInterfaceMockRecorder {
	return m.recorder
}

// ExpectedStats mocks base method.
func (m *MockNumbersInterface) ExpectedStats() (data.ExpectedStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpectedStats")
	ret0, _ := ret[0].(data.ExpectedStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExpectedStats indicates an expected call of ExpectedStats.
func (mr *MockNumbersInterfaceMockRecorder) ExpectedStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpectedStats", reflect.TypeOf((*MockNumbersInterface)(nil).ExpectedStats))
}
