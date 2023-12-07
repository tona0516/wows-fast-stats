// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	domain "wfs/backend/domain"

	mock "github.com/stretchr/testify/mock"
)

// WargamingInterface is an autogenerated mock type for the WargamingInterface type
type WargamingInterface struct {
	mock.Mock
}

// AccountInfo provides a mock function with given fields: accountIDs
func (_m *WargamingInterface) AccountInfo(accountIDs []int) (domain.WGAccountInfo, error) {
	ret := _m.Called(accountIDs)

	if len(ret) == 0 {
		panic("no return value specified for AccountInfo")
	}

	var r0 domain.WGAccountInfo
	var r1 error
	if rf, ok := ret.Get(0).(func([]int) (domain.WGAccountInfo, error)); ok {
		return rf(accountIDs)
	}
	if rf, ok := ret.Get(0).(func([]int) domain.WGAccountInfo); ok {
		r0 = rf(accountIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGAccountInfo)
		}
	}

	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(accountIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccountList provides a mock function with given fields: accountNames
func (_m *WargamingInterface) AccountList(accountNames []string) (domain.WGAccountList, error) {
	ret := _m.Called(accountNames)

	if len(ret) == 0 {
		panic("no return value specified for AccountList")
	}

	var r0 domain.WGAccountList
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) (domain.WGAccountList, error)); ok {
		return rf(accountNames)
	}
	if rf, ok := ret.Get(0).(func([]string) domain.WGAccountList); ok {
		r0 = rf(accountNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGAccountList)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(accountNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccountListForSearch provides a mock function with given fields: prefix
func (_m *WargamingInterface) AccountListForSearch(prefix string) (domain.WGAccountList, error) {
	ret := _m.Called(prefix)

	if len(ret) == 0 {
		panic("no return value specified for AccountListForSearch")
	}

	var r0 domain.WGAccountList
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.WGAccountList, error)); ok {
		return rf(prefix)
	}
	if rf, ok := ret.Get(0).(func(string) domain.WGAccountList); ok {
		r0 = rf(prefix)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGAccountList)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(prefix)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BattleArenas provides a mock function with given fields:
func (_m *WargamingInterface) BattleArenas() (domain.WGBattleArenas, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BattleArenas")
	}

	var r0 domain.WGBattleArenas
	var r1 error
	if rf, ok := ret.Get(0).(func() (domain.WGBattleArenas, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() domain.WGBattleArenas); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGBattleArenas)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BattleTypes provides a mock function with given fields:
func (_m *WargamingInterface) BattleTypes() (domain.WGBattleTypes, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BattleTypes")
	}

	var r0 domain.WGBattleTypes
	var r1 error
	if rf, ok := ret.Get(0).(func() (domain.WGBattleTypes, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() domain.WGBattleTypes); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGBattleTypes)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClansAccountInfo provides a mock function with given fields: accountIDs
func (_m *WargamingInterface) ClansAccountInfo(accountIDs []int) (domain.WGClansAccountInfo, error) {
	ret := _m.Called(accountIDs)

	if len(ret) == 0 {
		panic("no return value specified for ClansAccountInfo")
	}

	var r0 domain.WGClansAccountInfo
	var r1 error
	if rf, ok := ret.Get(0).(func([]int) (domain.WGClansAccountInfo, error)); ok {
		return rf(accountIDs)
	}
	if rf, ok := ret.Get(0).(func([]int) domain.WGClansAccountInfo); ok {
		r0 = rf(accountIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGClansAccountInfo)
		}
	}

	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(accountIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClansInfo provides a mock function with given fields: clanIDs
func (_m *WargamingInterface) ClansInfo(clanIDs []int) (domain.WGClansInfo, error) {
	ret := _m.Called(clanIDs)

	if len(ret) == 0 {
		panic("no return value specified for ClansInfo")
	}

	var r0 domain.WGClansInfo
	var r1 error
	if rf, ok := ret.Get(0).(func([]int) (domain.WGClansInfo, error)); ok {
		return rf(clanIDs)
	}
	if rf, ok := ret.Get(0).(func([]int) domain.WGClansInfo); ok {
		r0 = rf(clanIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGClansInfo)
		}
	}

	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(clanIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncycShips provides a mock function with given fields: pageNo
func (_m *WargamingInterface) EncycShips(pageNo int) (domain.WGEncycShips, int, error) {
	ret := _m.Called(pageNo)

	if len(ret) == 0 {
		panic("no return value specified for EncycShips")
	}

	var r0 domain.WGEncycShips
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(int) (domain.WGEncycShips, int, error)); ok {
		return rf(pageNo)
	}
	if rf, ok := ret.Get(0).(func(int) domain.WGEncycShips); ok {
		r0 = rf(pageNo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGEncycShips)
		}
	}

	if rf, ok := ret.Get(1).(func(int) int); ok {
		r1 = rf(pageNo)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(int) error); ok {
		r2 = rf(pageNo)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SetAppID provides a mock function with given fields: appid
func (_m *WargamingInterface) SetAppID(appid string) {
	_m.Called(appid)
}

// ShipsStats provides a mock function with given fields: accountID
func (_m *WargamingInterface) ShipsStats(accountID int) (domain.WGShipsStats, error) {
	ret := _m.Called(accountID)

	if len(ret) == 0 {
		panic("no return value specified for ShipsStats")
	}

	var r0 domain.WGShipsStats
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (domain.WGShipsStats, error)); ok {
		return rf(accountID)
	}
	if rf, ok := ret.Get(0).(func(int) domain.WGShipsStats); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.WGShipsStats)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Test provides a mock function with given fields: appid
func (_m *WargamingInterface) Test(appid string) (bool, error) {
	ret := _m.Called(appid)

	if len(ret) == 0 {
		panic("no return value specified for Test")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(appid)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(appid)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(appid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWargamingInterface creates a new instance of WargamingInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWargamingInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *WargamingInterface {
	mock := &WargamingInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}