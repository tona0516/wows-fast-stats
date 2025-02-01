// Code generated by MockGen. DO NOT EDIT.
// Source: discord.go
//
// Generated by this command:
//
//	mockgen -source=discord.go -destination ../mock/webapi/discord.go -package webapi
//

// Package webapi is a generated GoMock package.
package webapi

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDiscord is a mock of Discord interface.
type MockDiscord struct {
	ctrl     *gomock.Controller
	recorder *MockDiscordMockRecorder
	isgomock struct{}
}

// MockDiscordMockRecorder is the mock recorder for MockDiscord.
type MockDiscordMockRecorder struct {
	mock *MockDiscord
}

// NewMockDiscord creates a new mock instance.
func NewMockDiscord(ctrl *gomock.Controller) *MockDiscord {
	mock := &MockDiscord{ctrl: ctrl}
	mock.recorder = &MockDiscordMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDiscord) EXPECT() *MockDiscordMockRecorder {
	return m.recorder
}

// Comment mocks base method.
func (m *MockDiscord) Comment(message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Comment", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Comment indicates an expected call of Comment.
func (mr *MockDiscordMockRecorder) Comment(message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Comment", reflect.TypeOf((*MockDiscord)(nil).Comment), message)
}
