package service

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initMocksForBattle() (*mockWargaming, *mockNumbers, *mockUnregistered, *mockLocalFile) {
	mockWargaming := &mockWargaming{}
	mockNumbers := &mockNumbers{}
	mockUnregistered := &mockUnregistered{}
	mockLocalFile := &mockLocalFile{}

	mockWargaming.On("SetAppID", mock.Anything).Return()
	accountList := domain.WGAccountList{}
	accountList.Data = []domain.WGAccountListData{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}

	mockWargaming.On("AccountList", mock.Anything).Return(accountList, nil)
	mockWargaming.On("EncycShips", mock.Anything).Return(domain.WGEncycShips{}, nil)
	mockWargaming.On("BattleArenas").Return(domain.WGBattleArenas{}, nil)
	mockWargaming.On("BattleTypes").Return(domain.WGBattleTypes{}, nil)
	mockWargaming.On("AccountInfo", mock.Anything).Return(domain.WGAccountInfo{}, nil)
	mockWargaming.On("ShipsStats", mock.Anything).Return(domain.WGShipsStats{}, nil)
	mockWargaming.On("ClansAccountInfo", mock.Anything).Return(domain.WGClansAccountInfo{}, nil)
	mockWargaming.On("ClansInfo", mock.Anything).Return(domain.WGClansInfo{}, nil)

	mockNumbers.On("ExpectedStats").Return(domain.NSExpectedStats{}, nil)

	mockUnregistered.On("Warship").Return(map[int]domain.Warship{}, nil)

	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{}, nil)
	mockLocalFile.On("SaveTempArenaInfo", mock.Anything).Return(nil)

	return mockWargaming, mockNumbers, mockUnregistered, mockLocalFile
}

func TestBattle_Battle_正常系_初回(t *testing.T) {
	t.Parallel()

	mockWargaming, mockNumbers, mockUnregistered, mockLocalFileRepo := initMocksForBattle()

	// テスト
	b := NewBattle(5, mockWargaming, mockLocalFileRepo, mockNumbers, mockUnregistered)
	_, err := b.Battle(domain.UserConfig{})

	// アサーション
	assert.NoError(t, err)

	mockWargaming.AssertCalled(t, "SetAppID", mock.Anything)
	mockWargaming.AssertCalled(t, "AccountList", mock.Anything)
	mockWargaming.AssertCalled(t, "EncycShips", mock.Anything)
	mockWargaming.AssertCalled(t, "BattleArenas")
	mockWargaming.AssertCalled(t, "BattleTypes")
	mockWargaming.AssertCalled(t, "AccountInfo", mock.Anything)
	mockWargaming.AssertNumberOfCalls(t, "ShipsStats", 2)
	mockWargaming.AssertCalled(t, "ClansAccountInfo", mock.Anything)
	mockWargaming.AssertCalled(t, "ClansInfo", mock.Anything)

	mockNumbers.AssertCalled(t, "ExpectedStats")

	mockUnregistered.AssertCalled(t, "Warship")

	mockLocalFileRepo.AssertCalled(t, "TempArenaInfo", mock.Anything)
	mockLocalFileRepo.AssertNotCalled(t, "SaveTempArenaInfo", mock.Anything)
}

func TestBattle_Battle_正常系_2回目以降(t *testing.T) {
	t.Parallel()

	mockWargaming, mockNumbers, mockUnregistered, mockLocalFile := initMocksForBattle()

	// テスト
	b := NewBattle(5, mockWargaming, mockLocalFile, mockNumbers, mockUnregistered)
	b.isFirstBattle = false
	_, err := b.Battle(domain.UserConfig{})

	// アサーション
	assert.NoError(t, err)

	mockWargaming.AssertCalled(t, "SetAppID", mock.Anything)
	mockWargaming.AssertCalled(t, "AccountList", mock.Anything)
	mockWargaming.AssertNotCalled(t, "EncycShips", mock.Anything)
	mockWargaming.AssertNotCalled(t, "BattleArenas")
	mockWargaming.AssertNotCalled(t, "BattleTypes")
	mockWargaming.AssertCalled(t, "AccountInfo", mock.Anything)
	mockWargaming.AssertNumberOfCalls(t, "ShipsStats", 2)
	mockWargaming.AssertCalled(t, "ClansAccountInfo", mock.Anything)
	mockWargaming.AssertCalled(t, "ClansInfo", mock.Anything)

	mockNumbers.AssertNotCalled(t, "ExpectedStats")

	mockUnregistered.AssertNotCalled(t, "Warship")

	mockLocalFile.AssertCalled(t, "TempArenaInfo", mock.Anything)
	mockLocalFile.AssertNotCalled(t, "SaveTempArenaInfo", mock.Anything)
}

func TestBattle_Battle_異常系(t *testing.T) {
	t.Parallel()

	mockWargaming, mockNumbers, mockUnregistered, _ := initMocksForBattle()
	mockLocalFileRepo := &mockLocalFile{}
	expectedError := failure.New(apperr.FileNotExist)
	mockLocalFileRepo.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{}, expectedError)
	mockLocalFileRepo.On("SaveTempArenaInfo", mock.Anything).Return(nil)

	// テスト
	b := NewBattle(5, mockWargaming, mockLocalFileRepo, mockNumbers, mockUnregistered)
	b.isFirstBattle = false
	_, err := b.Battle(domain.UserConfig{})

	// アサーション
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, code, apperr.FileNotExist)

	mockWargaming.AssertCalled(t, "SetAppID", mock.Anything)
	mockWargaming.AssertNotCalled(t, "AccountList", mock.Anything)
	mockWargaming.AssertNotCalled(t, "EncycShips", mock.Anything)
	mockWargaming.AssertNotCalled(t, "BattleArenas")
	mockWargaming.AssertNotCalled(t, "BattleTypes")
	mockWargaming.AssertNotCalled(t, "AccountInfo", mock.Anything)
	mockWargaming.AssertNotCalled(t, "ShipsStats", mock.Anything)
	mockWargaming.AssertNotCalled(t, "ClansAccountInfo", mock.Anything)
	mockWargaming.AssertNotCalled(t, "ClansInfo", mock.Anything)

	mockNumbers.AssertNotCalled(t, "ExpectedStats")

	mockUnregistered.AssertNotCalled(t, "Warship")

	mockLocalFileRepo.AssertCalled(t, "TempArenaInfo", mock.Anything)
	mockLocalFileRepo.AssertNotCalled(t, "SaveTempArenaInfo", mock.Anything)
}
