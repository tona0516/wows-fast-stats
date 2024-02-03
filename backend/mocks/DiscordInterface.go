// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DiscordInterface is an autogenerated mock type for the DiscordInterface type
type DiscordInterface struct {
	mock.Mock
}

// Comment provides a mock function with given fields: message
func (_m *DiscordInterface) Comment(message string) error {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for Comment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDiscordInterface creates a new instance of DiscordInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDiscordInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DiscordInterface {
	mock := &DiscordInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
