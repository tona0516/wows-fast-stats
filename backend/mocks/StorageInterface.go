// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	domain "wfs/backend/domain"

	mock "github.com/stretchr/testify/mock"
)

// StorageInterface is an autogenerated mock type for the StorageInterface type
type StorageInterface struct {
	mock.Mock
}

// IsExistAlertPlayers provides a mock function with given fields:
func (_m *StorageInterface) IsExistAlertPlayers() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsExistAlertPlayers")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsExistUserConfig provides a mock function with given fields:
func (_m *StorageInterface) IsExistUserConfig() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsExistUserConfig")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ReadAlertPlayers provides a mock function with given fields:
func (_m *StorageInterface) ReadAlertPlayers() ([]domain.AlertPlayer, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadAlertPlayers")
	}

	var r0 []domain.AlertPlayer
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.AlertPlayer, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.AlertPlayer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.AlertPlayer)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadDataVersion provides a mock function with given fields:
func (_m *StorageInterface) ReadDataVersion() (uint, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadDataVersion")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadNSExpectedStats provides a mock function with given fields:
func (_m *StorageInterface) ReadNSExpectedStats() (domain.NSExpectedStats, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadNSExpectedStats")
	}

	var r0 domain.NSExpectedStats
	var r1 error
	if rf, ok := ret.Get(0).(func() (domain.NSExpectedStats, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() domain.NSExpectedStats); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.NSExpectedStats)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadUserConfig provides a mock function with given fields:
func (_m *StorageInterface) ReadUserConfig() (domain.UserConfig, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadUserConfig")
	}

	var r0 domain.UserConfig
	var r1 error
	if rf, ok := ret.Get(0).(func() (domain.UserConfig, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() domain.UserConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.UserConfig)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteAlertPlayers provides a mock function with given fields: players
func (_m *StorageInterface) WriteAlertPlayers(players []domain.AlertPlayer) error {
	ret := _m.Called(players)

	if len(ret) == 0 {
		panic("no return value specified for WriteAlertPlayers")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]domain.AlertPlayer) error); ok {
		r0 = rf(players)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteDataVersion provides a mock function with given fields: version
func (_m *StorageInterface) WriteDataVersion(version uint) error {
	ret := _m.Called(version)

	if len(ret) == 0 {
		panic("no return value specified for WriteDataVersion")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(version)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteNSExpectedStats provides a mock function with given fields: nsExpectedStats
func (_m *StorageInterface) WriteNSExpectedStats(nsExpectedStats domain.NSExpectedStats) error {
	ret := _m.Called(nsExpectedStats)

	if len(ret) == 0 {
		panic("no return value specified for WriteNSExpectedStats")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.NSExpectedStats) error); ok {
		r0 = rf(nsExpectedStats)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteUserConfig provides a mock function with given fields: config
func (_m *StorageInterface) WriteUserConfig(config domain.UserConfig) error {
	ret := _m.Called(config)

	if len(ret) == 0 {
		panic("no return value specified for WriteUserConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.UserConfig) error); ok {
		r0 = rf(config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorageInterface creates a new instance of StorageInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageInterface {
	mock := &StorageInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}