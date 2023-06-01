package service

import (
	"changeme/backend/vo"

	"github.com/stretchr/testify/mock"
)

type mockConfigRepo struct {
	mock.Mock
}

func (m *mockConfigRepo) User() (vo.UserConfig, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.UserConfig), args.Error(1)
}

func (m *mockConfigRepo) UpdateUser(config vo.UserConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *mockConfigRepo) App() (vo.AppConfig, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.AppConfig), args.Error(1)
}

func (m *mockConfigRepo) UpdateApp(config vo.AppConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *mockConfigRepo) AlertPlayers() ([]vo.AlertPlayer, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).([]vo.AlertPlayer), args.Error(1)
}

func (m *mockConfigRepo) UpdateAlertPlayer(player vo.AlertPlayer) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *mockConfigRepo) RemoveAlertPlayer(accountID int) error {
	args := m.Called(accountID)
	return args.Error(0)
}

type mockTempArenaInfoRepo struct {
	mock.Mock
}

func (m *mockTempArenaInfoRepo) Get(installPath string) (vo.TempArenaInfo, error) {
	args := m.Called(installPath)
	//nolint:forcetypeassert
	return args.Get(0).(vo.TempArenaInfo), args.Error(1)
}

func (m *mockTempArenaInfoRepo) Save(tempArenaInfo vo.TempArenaInfo) error {
	args := m.Called(tempArenaInfo)
	return args.Error(0)
}

type mockWargamingRepo struct {
	mock.Mock
}

func (m *mockWargamingRepo) SetAppID(appID string) {
	m.Called(appID)
}

func (m *mockWargamingRepo) AccountInfo(accountIDs []int) (vo.WGAccountInfo, error) {
	args := m.Called(accountIDs)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGAccountInfo), args.Error(1)
}

func (m *mockWargamingRepo) AccountList(accountNames []string) (vo.WGAccountList, error) {
	args := m.Called(accountNames)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGAccountList), args.Error(1)
}

func (m *mockWargamingRepo) AccountListForSearch(prefix string) (vo.WGAccountList, error) {
	args := m.Called(prefix)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGAccountList), args.Error(1)
}

func (m *mockWargamingRepo) ClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error) {
	args := m.Called(accountIDs)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGClansAccountInfo), args.Error(1)
}

func (m *mockWargamingRepo) ClansInfo(clanIDs []int) (vo.WGClansInfo, error) {
	args := m.Called(clanIDs)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGClansInfo), args.Error(1)
}

func (m *mockWargamingRepo) EncycShips(pageNo int) (vo.WGEncycShips, error) {
	args := m.Called(pageNo)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGEncycShips), args.Error(1)
}

func (m *mockWargamingRepo) ShipsStats(accountID int) (vo.WGShipsStats, error) {
	args := m.Called(accountID)
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGShipsStats), args.Error(1)
}

func (m *mockWargamingRepo) EncycInfo() (vo.WGEncycInfo, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGEncycInfo), args.Error(1)
}

func (m *mockWargamingRepo) BattleArenas() (vo.WGBattleArenas, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGBattleArenas), args.Error(1)
}

func (m *mockWargamingRepo) BattleTypes() (vo.WGBattleTypes, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.WGBattleTypes), args.Error(1)
}

type mockNumbersRepo struct {
	mock.Mock
}

func (m *mockNumbersRepo) ExpectedStats() (vo.NSExpectedStats, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.NSExpectedStats), args.Error(1)
}

type mockUnregisteredRepo struct {
	mock.Mock
}

func (m *mockUnregisteredRepo) Warship() (map[int]vo.Warship, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(map[int]vo.Warship), args.Error(1)
}

type mockScreenshotRepo struct {
	mock.Mock
}

func (m *mockScreenshotRepo) Save(path string, base64Data string) error {
	args := m.Called(path, base64Data)
	return args.Error(0)
}
