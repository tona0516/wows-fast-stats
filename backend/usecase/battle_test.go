package usecase

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBattle_Get_正常系_初回(t *testing.T) {
	t.Parallel()

	// 準備
	mockWargaming := &mocks.WargamingInterface{}
	mockNumbers := &mocks.NumbersInterface{}
	mockUnregistered := &mocks.UnregisteredInterface{}
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}

	mockWargaming.On("SetAppID", mock.Anything).Return()
	mockWargaming.On("AccountList", mock.Anything).Return(domain.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.On("EncycShips", mock.Anything).Return(domain.WGEncycShips{
		1: domain.WGEncycShipsData{
			Tier:      1,
			Type:      "Battleship",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, 2, nil)
	mockWargaming.On("BattleArenas").Return(domain.WGBattleArenas{}, nil)
	mockWargaming.On("BattleTypes").Return(domain.WGBattleTypes{}, nil)
	mockWargaming.On("AccountInfo", mock.Anything).Return(domain.WGAccountInfo{}, nil)
	mockWargaming.On("ShipsStats", mock.Anything).Return(domain.WGShipsStats{}, nil)
	mockWargaming.On("ClansAccountInfo", mock.Anything).Return(domain.WGClansAccountInfo{}, nil)
	mockWargaming.On("ClansInfo", mock.Anything).Return(domain.WGClansInfo{}, nil)

	mockNumbers.On("ExpectedStats").Return(domain.ExpectedStats{}, nil)

	mockUnregistered.On("Warship").Return(domain.Warships{}, nil)

	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{
		Vehicles: []domain.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockStorage.On("WriteOwnIGN", mock.Anything).Return(nil)
	mockStorage.On("WriteExpectedStats", mock.Anything).Return(nil)

	// テスト
	b := NewBattle(5, mockWargaming, mockLocalFile, mockNumbers, mockUnregistered, mockStorage, nil)
	_, err := b.Get(domain.UserConfig{})

	// アサーション
	require.NoError(t, err)
	mockWargaming.AssertExpectations(t)
	mockNumbers.AssertExpectations(t)
	mockUnregistered.AssertExpectations(t)
	mockLocalFile.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestBattle_Get_正常系_2回目以降(t *testing.T) {
	t.Parallel()

	// 準備
	mockWargaming := &mocks.WargamingInterface{}
	mockNumbers := &mocks.NumbersInterface{}
	mockUnregistered := &mocks.UnregisteredInterface{}
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}

	mockWargaming.On("SetAppID", mock.Anything).Return()
	mockWargaming.On("AccountList", mock.Anything).Return(domain.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.On("AccountInfo", mock.Anything).Return(domain.WGAccountInfo{}, nil)
	mockWargaming.On("ShipsStats", mock.Anything).Return(domain.WGShipsStats{}, nil)
	mockWargaming.On("ClansAccountInfo", mock.Anything).Return(domain.WGClansAccountInfo{}, nil)
	mockWargaming.On("ClansInfo", mock.Anything).Return(domain.WGClansInfo{}, nil)

	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{
		Vehicles: []domain.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockStorage.On("WriteOwnIGN", mock.Anything).Return(nil)

	// テスト
	b := NewBattle(5, mockWargaming, mockLocalFile, mockNumbers, mockUnregistered, mockStorage, nil)
	b.isFirstBattle = false
	_, err := b.Get(domain.UserConfig{})

	// アサーション
	require.NoError(t, err)
	mockWargaming.AssertExpectations(t)
	mockNumbers.AssertExpectations(t)
	mockUnregistered.AssertExpectations(t)
	mockLocalFile.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestBattle_Get_異常系(t *testing.T) {
	t.Parallel()

	// 準備
	mockWargaming := &mocks.WargamingInterface{}
	mockNumbers := &mocks.NumbersInterface{}
	mockUnregistered := &mocks.UnregisteredInterface{}
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}

	mockWargaming.On("SetAppID", mock.Anything).Return()

	expectedError := failure.New(apperr.FileNotExist)
	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{}, expectedError)

	// テスト
	b := NewBattle(5, mockWargaming, mockLocalFile, mockNumbers, mockUnregistered, mockStorage, nil)
	b.isFirstBattle = false
	_, err := b.Get(domain.UserConfig{})

	// アサーション
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, apperr.FileNotExist, code)

	mockWargaming.AssertExpectations(t)
	mockNumbers.AssertExpectations(t)
	mockUnregistered.AssertExpectations(t)
	mockLocalFile.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}
