package service

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/vo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initMocksForBattle() (*mockWargamingRepo, *mockNumbersRepo, *mockUnregisteredRepo, *mockTempArenaInfoRepo) {
	mockWargamingRepo := &mockWargamingRepo{}
	mockNumbersRepo := &mockNumbersRepo{}
	mockUnregisteredRepo := &mockUnregisteredRepo{}
	mockTempArenaInfoRepo := &mockTempArenaInfoRepo{}

	mockWargamingRepo.On("SetAppID", mock.Anything).Return()
	accountList := vo.WGAccountList{
		Data: []vo.WGAccountListData{
			{NickName: "player_1", AccountID: 1},
			{NickName: "player_2", AccountID: 2},
		},
	}
	mockWargamingRepo.On("AccountList", mock.Anything).Return(accountList, nil)
	mockWargamingRepo.On("EncycShips", mock.Anything).Return(vo.WGEncycShips{}, nil)
	mockWargamingRepo.On("BattleArenas").Return(vo.WGBattleArenas{}, nil)
	mockWargamingRepo.On("BattleTypes").Return(vo.WGBattleTypes{}, nil)
	mockWargamingRepo.On("AccountInfo", mock.Anything).Return(vo.WGAccountInfo{}, nil)
	mockWargamingRepo.On("ShipsStats", mock.Anything).Return(vo.WGShipsStats{}, nil)
	mockWargamingRepo.On("ClansAccountInfo", mock.Anything).Return(vo.WGClansAccountInfo{}, nil)
	mockWargamingRepo.On("ClansInfo", mock.Anything).Return(vo.WGClansInfo{}, nil)

	mockNumbersRepo.On("ExpectedStats").Return(vo.NSExpectedStats{}, nil)

	mockUnregisteredRepo.On("Warship").Return(map[int]vo.Warship{}, nil)

	mockTempArenaInfoRepo.On("Get", mock.Anything).Return(vo.TempArenaInfo{}, nil)
	mockTempArenaInfoRepo.On("Save", mock.Anything).Return(nil)

	return mockWargamingRepo, mockNumbersRepo, mockUnregisteredRepo, mockTempArenaInfoRepo
}

func TestBattle_Battle_正常系_初回(t *testing.T) {
	t.Parallel()

	mockWargamingRepo, mockNumbersRepo, mockUnregisteredRepo, mockTempArenaInfoRepo := initMocksForBattle()

	// テスト
	b := NewBattle(5, mockWargamingRepo, mockTempArenaInfoRepo, mockNumbersRepo, mockUnregisteredRepo)
	_, err := b.Battle(vo.UserConfig{})

	// アサーション
	assert.NoError(t, err)

	mockWargamingRepo.AssertCalled(t, "SetAppID", mock.Anything)
	mockWargamingRepo.AssertCalled(t, "AccountList", mock.Anything)
	mockWargamingRepo.AssertCalled(t, "EncycShips", mock.Anything)
	mockWargamingRepo.AssertCalled(t, "BattleArenas")
	mockWargamingRepo.AssertCalled(t, "BattleTypes")
	mockWargamingRepo.AssertCalled(t, "AccountInfo", mock.Anything)
	mockWargamingRepo.AssertNumberOfCalls(t, "ShipsStats", 2)
	mockWargamingRepo.AssertCalled(t, "ClansAccountInfo", mock.Anything)
	mockWargamingRepo.AssertCalled(t, "ClansInfo", mock.Anything)

	mockNumbersRepo.AssertCalled(t, "ExpectedStats")

	mockUnregisteredRepo.AssertCalled(t, "Warship")

	mockTempArenaInfoRepo.AssertCalled(t, "Get", mock.Anything)
	mockTempArenaInfoRepo.AssertNotCalled(t, "Save", mock.Anything)
}

func TestBattle_Battle_正常系_2回目以降(t *testing.T) {
	t.Parallel()

	mockWargamingRepo, mockNumbersRepo, mockUnregisteredRepo, mockTempArenaInfoRepo := initMocksForBattle()

	// テスト
	b := NewBattle(5, mockWargamingRepo, mockTempArenaInfoRepo, mockNumbersRepo, mockUnregisteredRepo)
	b.isFirstBattle = false
	_, err := b.Battle(vo.UserConfig{})

	// アサーション
	assert.NoError(t, err)

	mockWargamingRepo.AssertCalled(t, "SetAppID", mock.Anything)
	mockWargamingRepo.AssertCalled(t, "AccountList", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "EncycShips", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "BattleArenas")
	mockWargamingRepo.AssertNotCalled(t, "BattleTypes")
	mockWargamingRepo.AssertCalled(t, "AccountInfo", mock.Anything)
	mockWargamingRepo.AssertNumberOfCalls(t, "ShipsStats", 2)
	mockWargamingRepo.AssertCalled(t, "ClansAccountInfo", mock.Anything)
	mockWargamingRepo.AssertCalled(t, "ClansInfo", mock.Anything)

	mockNumbersRepo.AssertNotCalled(t, "ExpectedStats")

	mockUnregisteredRepo.AssertNotCalled(t, "Warship")

	mockTempArenaInfoRepo.AssertCalled(t, "Get", mock.Anything)
	mockTempArenaInfoRepo.AssertNotCalled(t, "Save", mock.Anything)
}

func TestBattle_Battle_異常系(t *testing.T) {
	t.Parallel()

	mockWargamingRepo, mockNumbersRepo, mockUnregisteredRepo, _ := initMocksForBattle()
	mockTempArenaInfoRepo := &mockTempArenaInfoRepo{}
	expectedError := apperr.New(apperr.ReadFile, nil)
	mockTempArenaInfoRepo.On("Get", mock.Anything).Return(vo.TempArenaInfo{}, expectedError)
	mockTempArenaInfoRepo.On("Save", mock.Anything).Return(nil)

	// テスト
	b := NewBattle(5, mockWargamingRepo, mockTempArenaInfoRepo, mockNumbersRepo, mockUnregisteredRepo)
	b.isFirstBattle = false
	_, err := b.Battle(vo.UserConfig{})

	// アサーション
	assert.EqualError(t, err, expectedError.Error())

	mockWargamingRepo.AssertCalled(t, "SetAppID", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "AccountList", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "EncycShips", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "BattleArenas")
	mockWargamingRepo.AssertNotCalled(t, "BattleTypes")
	mockWargamingRepo.AssertNotCalled(t, "AccountInfo", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "ShipsStats", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "ClansAccountInfo", mock.Anything)
	mockWargamingRepo.AssertNotCalled(t, "ClansInfo", mock.Anything)

	mockNumbersRepo.AssertNotCalled(t, "ExpectedStats")

	mockUnregisteredRepo.AssertNotCalled(t, "Warship")

	mockTempArenaInfoRepo.AssertCalled(t, "Get", mock.Anything)
	mockTempArenaInfoRepo.AssertNotCalled(t, "Save", mock.Anything)
}
