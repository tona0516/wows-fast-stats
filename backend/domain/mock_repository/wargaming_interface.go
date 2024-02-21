// Code generated by MockGen. DO NOT EDIT.
// Source: wargaming_interface.go
//
// Generated by this command:
//
//	mockgen -source=wargaming_interface.go -destination ../mock_repository/wargaming_interface.go -package mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	model "wfs/backend/domain/model"

	gomock "go.uber.org/mock/gomock"
)

// MockWargamingInterface is a mock of WargamingInterface interface.
type MockWargamingInterface struct {
	ctrl     *gomock.Controller
	recorder *MockWargamingInterfaceMockRecorder
}

// MockWargamingInterfaceMockRecorder is the mock recorder for MockWargamingInterface.
type MockWargamingInterfaceMockRecorder struct {
	mock *MockWargamingInterface
}

// NewMockWargamingInterface creates a new mock instance.
func NewMockWargamingInterface(ctrl *gomock.Controller) *MockWargamingInterface {
	mock := &MockWargamingInterface{ctrl: ctrl}
	mock.recorder = &MockWargamingInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWargamingInterface) EXPECT() *MockWargamingInterfaceMockRecorder {
	return m.recorder
}

// AccountInfo mocks base method.
func (m *MockWargamingInterface) AccountInfo(appID string, accountIDs []int) (model.WGAccountInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountInfo", appID, accountIDs)
	ret0, _ := ret[0].(model.WGAccountInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountInfo indicates an expected call of AccountInfo.
func (mr *MockWargamingInterfaceMockRecorder) AccountInfo(appID, accountIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountInfo", reflect.TypeOf((*MockWargamingInterface)(nil).AccountInfo), appID, accountIDs)
}

// AccountList mocks base method.
func (m *MockWargamingInterface) AccountList(appID string, accountNames []string) (model.WGAccountList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountList", appID, accountNames)
	ret0, _ := ret[0].(model.WGAccountList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountList indicates an expected call of AccountList.
func (mr *MockWargamingInterfaceMockRecorder) AccountList(appID, accountNames any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountList", reflect.TypeOf((*MockWargamingInterface)(nil).AccountList), appID, accountNames)
}

// AccountListForSearch mocks base method.
func (m *MockWargamingInterface) AccountListForSearch(appID, prefix string) (model.WGAccountList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountListForSearch", appID, prefix)
	ret0, _ := ret[0].(model.WGAccountList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountListForSearch indicates an expected call of AccountListForSearch.
func (mr *MockWargamingInterfaceMockRecorder) AccountListForSearch(appID, prefix any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountListForSearch", reflect.TypeOf((*MockWargamingInterface)(nil).AccountListForSearch), appID, prefix)
}

// BattleArenas mocks base method.
func (m *MockWargamingInterface) BattleArenas(appID string) (model.WGBattleArenas, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BattleArenas", appID)
	ret0, _ := ret[0].(model.WGBattleArenas)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BattleArenas indicates an expected call of BattleArenas.
func (mr *MockWargamingInterfaceMockRecorder) BattleArenas(appID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BattleArenas", reflect.TypeOf((*MockWargamingInterface)(nil).BattleArenas), appID)
}

// BattleTypes mocks base method.
func (m *MockWargamingInterface) BattleTypes(appID string) (model.WGBattleTypes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BattleTypes", appID)
	ret0, _ := ret[0].(model.WGBattleTypes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BattleTypes indicates an expected call of BattleTypes.
func (mr *MockWargamingInterfaceMockRecorder) BattleTypes(appID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BattleTypes", reflect.TypeOf((*MockWargamingInterface)(nil).BattleTypes), appID)
}

// ClansAccountInfo mocks base method.
func (m *MockWargamingInterface) ClansAccountInfo(appID string, accountIDs []int) (model.WGClansAccountInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClansAccountInfo", appID, accountIDs)
	ret0, _ := ret[0].(model.WGClansAccountInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClansAccountInfo indicates an expected call of ClansAccountInfo.
func (mr *MockWargamingInterfaceMockRecorder) ClansAccountInfo(appID, accountIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClansAccountInfo", reflect.TypeOf((*MockWargamingInterface)(nil).ClansAccountInfo), appID, accountIDs)
}

// ClansInfo mocks base method.
func (m *MockWargamingInterface) ClansInfo(appID string, clanIDs []int) (model.WGClansInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClansInfo", appID, clanIDs)
	ret0, _ := ret[0].(model.WGClansInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClansInfo indicates an expected call of ClansInfo.
func (mr *MockWargamingInterfaceMockRecorder) ClansInfo(appID, clanIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClansInfo", reflect.TypeOf((*MockWargamingInterface)(nil).ClansInfo), appID, clanIDs)
}

// EncycShips mocks base method.
func (m *MockWargamingInterface) EncycShips(appID string, pageNo int) (model.WGEncycShips, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncycShips", appID, pageNo)
	ret0, _ := ret[0].(model.WGEncycShips)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EncycShips indicates an expected call of EncycShips.
func (mr *MockWargamingInterfaceMockRecorder) EncycShips(appID, pageNo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncycShips", reflect.TypeOf((*MockWargamingInterface)(nil).EncycShips), appID, pageNo)
}

// ShipsStats mocks base method.
func (m *MockWargamingInterface) ShipsStats(appID string, accountID int) (model.WGShipsStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShipsStats", appID, accountID)
	ret0, _ := ret[0].(model.WGShipsStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShipsStats indicates an expected call of ShipsStats.
func (mr *MockWargamingInterfaceMockRecorder) ShipsStats(appID, accountID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShipsStats", reflect.TypeOf((*MockWargamingInterface)(nil).ShipsStats), appID, accountID)
}

// Test mocks base method.
func (m *MockWargamingInterface) Test(appID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Test", appID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Test indicates an expected call of Test.
func (mr *MockWargamingInterfaceMockRecorder) Test(appID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Test", reflect.TypeOf((*MockWargamingInterface)(nil).Test), appID)
}
