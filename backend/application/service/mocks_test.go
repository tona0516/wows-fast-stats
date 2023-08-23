package service

import (
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/stretchr/testify/mock"
)

type mockLocalFile struct {
	mock.Mock
}

func (m *mockLocalFile) User() (domain.UserConfig, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(domain.UserConfig), args.Error(1)
}

func (m *mockLocalFile) UpdateUser(config domain.UserConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *mockLocalFile) App() (vo.AppConfig, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(vo.AppConfig), args.Error(1)
}

func (m *mockLocalFile) UpdateApp(config vo.AppConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *mockLocalFile) AlertPlayers() ([]domain.AlertPlayer, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).([]domain.AlertPlayer), args.Error(1)
}

func (m *mockLocalFile) UpdateAlertPlayer(player domain.AlertPlayer) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *mockLocalFile) RemoveAlertPlayer(accountID int) error {
	args := m.Called(accountID)
	return args.Error(0)
}

func (m *mockLocalFile) SaveScreenshot(path string, base64Data string) error {
	args := m.Called(path, base64Data)
	return args.Error(0)
}

func (m *mockLocalFile) TempArenaInfo(installPath string) (domain.TempArenaInfo, error) {
	args := m.Called(installPath)
	//nolint:forcetypeassert
	return args.Get(0).(domain.TempArenaInfo), args.Error(1)
}

func (m *mockLocalFile) SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error {
	args := m.Called(tempArenaInfo)
	return args.Error(0)
}

func (m *mockLocalFile) CachedNSExpectedStats() (domain.NSExpectedStats, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(domain.NSExpectedStats), args.Error(1)
}

func (m *mockLocalFile) SaveNSExpectedStats(expectedStats domain.NSExpectedStats) error {
	args := m.Called(expectedStats)
	return args.Error(0)
}

type mockWargaming struct {
	mock.Mock
}

func (m *mockWargaming) SetAppID(appID string) {
	m.Called(appID)
}

func (m *mockWargaming) AccountInfo(accountIDs []int) (domain.WGAccountInfo, error) {
	args := m.Called(accountIDs)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGAccountInfo), args.Error(1)
}

func (m *mockWargaming) AccountList(accountNames []string) (domain.WGAccountList, error) {
	args := m.Called(accountNames)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGAccountList), args.Error(1)
}

func (m *mockWargaming) AccountListForSearch(prefix string) (domain.WGAccountList, error) {
	args := m.Called(prefix)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGAccountList), args.Error(1)
}

func (m *mockWargaming) ClansAccountInfo(accountIDs []int) (domain.WGClansAccountInfo, error) {
	args := m.Called(accountIDs)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGClansAccountInfo), args.Error(1)
}

func (m *mockWargaming) ClansInfo(clanIDs []int) (domain.WGClansInfo, error) {
	args := m.Called(clanIDs)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGClansInfo), args.Error(1)
}

func (m *mockWargaming) EncycShips(pageNo int) (domain.WGEncycShips, error) {
	args := m.Called(pageNo)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGEncycShips), args.Error(1)
}

func (m *mockWargaming) ShipsStats(accountID int) (domain.WGShipsStats, error) {
	args := m.Called(accountID)
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGShipsStats), args.Error(1)
}

func (m *mockWargaming) BattleArenas() (domain.WGBattleArenas, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGBattleArenas), args.Error(1)
}

func (m *mockWargaming) BattleTypes() (domain.WGBattleTypes, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(domain.WGBattleTypes), args.Error(1)
}

func (m *mockWargaming) Test(appid string) (bool, error) {
	args := m.Called(appid)
	//nolint:forcetypeassert
	return args.Get(0).(bool), args.Error(1)
}

type mockNumbers struct {
	mock.Mock
}

func (m *mockNumbers) ExpectedStats() (domain.NSExpectedStats, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(domain.NSExpectedStats), args.Error(1)
}

type mockUnregistered struct {
	mock.Mock
}

func (m *mockUnregistered) Warship() (map[int]domain.Warship, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(map[int]domain.Warship), args.Error(1)
}

type mockGithub struct {
	mock.Mock
}

func (m *mockGithub) LatestRelease() (domain.GHLatestRelease, error) {
	args := m.Called()
	//nolint:forcetypeassert
	return args.Get(0).(domain.GHLatestRelease), args.Error(1)
}
