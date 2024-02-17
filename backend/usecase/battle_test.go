package usecase

import (
	"context"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var testUserConfig = model.UserConfigV2{
	Appid:             "test_appid",
	InstallPath:       "test_install_path",
	SaveTempArenaInfo: false,
}

func TestBattle_Get_正常系_初回(t *testing.T) {
	t.Parallel()

	// 準備
	mockWargaming := &mocks.WargamingInterface{}
	mockUnofficialWargaming := &mocks.UnofficialWargamingInterface{}
	mockNumbers := &mocks.NumbersInterface{}
	mockUnregistered := &mocks.UnregisteredInterface{}
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}
	mockLogger := &mocks.LoggerInterface{}

	mockWargaming.On("AccountList", mock.Anything, mock.Anything).Return(model.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.On("EncycShips", mock.Anything, mock.Anything).Return(model.WGEncycShips{
		1: model.WGEncycShipsData{
			Tier:      1,
			Type:      "Battleship",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, 2, nil)
	mockWargaming.On("BattleArenas", mock.Anything).Return(model.WGBattleArenas{}, nil)
	mockWargaming.On("BattleTypes", mock.Anything).Return(model.WGBattleTypes{}, nil)
	mockWargaming.On("AccountInfo", mock.Anything, mock.Anything).Return(model.WGAccountInfo{}, nil)
	mockWargaming.On("ShipsStats", mock.Anything, mock.Anything).Return(model.WGShipsStats{}, nil)
	mockWargaming.On("ClansAccountInfo", mock.Anything, mock.Anything).Return(model.WGClansAccountInfo{}, nil)
	mockWargaming.On("ClansInfo", mock.Anything, mock.Anything).Return(model.WGClansInfo{
		1: model.WGClansInfoData{Tag: "TEST"},
	}, nil)
	mockUnofficialWargaming.On("ClansAutoComplete", mock.Anything).Return(model.UWGClansAutocomplete{}, nil)

	mockNumbers.On("ExpectedStats").Return(model.ExpectedStats{}, nil)

	mockUnregistered.On("Warship").Return(model.Warships{}, nil)

	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockStorage.On("WriteOwnIGN", mock.Anything).Return(nil)
	mockStorage.On("WriteExpectedStats", mock.Anything).Return(nil)

	mockLogger.On("SetOwnIGN", mock.Anything).Return()

	// テスト
	b := NewBattle(
		5,
		mockWargaming,
		mockUnofficialWargaming,
		mockLocalFile,
		mockNumbers,
		mockUnregistered,
		mockStorage,
		mockLogger,
		nil,
	)
	_, err := b.Get(context.TODO(), testUserConfig)

	// アサーション
	require.NoError(t, err)
	mockWargaming.AssertExpectations(t)
	mockUnofficialWargaming.AssertExpectations(t)
	mockNumbers.AssertExpectations(t)
	mockUnregistered.AssertExpectations(t)
	mockLocalFile.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestBattle_Get_正常系_2回目以降(t *testing.T) {
	t.Parallel()

	// 準備
	mockWargaming := &mocks.WargamingInterface{}
	mockUnofficialWargaming := &mocks.UnofficialWargamingInterface{}
	mockNumbers := &mocks.NumbersInterface{}
	mockUnregistered := &mocks.UnregisteredInterface{}
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}
	mockLogger := &mocks.LoggerInterface{}

	mockWargaming.On("AccountList", mock.Anything, mock.Anything).Return(model.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.On("AccountInfo", mock.Anything, mock.Anything).Return(model.WGAccountInfo{}, nil)
	mockWargaming.On("ShipsStats", mock.Anything, mock.Anything).Return(model.WGShipsStats{}, nil)
	mockWargaming.On("ClansAccountInfo", mock.Anything, mock.Anything).Return(model.WGClansAccountInfo{}, nil)
	mockWargaming.On("ClansInfo", mock.Anything, mock.Anything).Return(model.WGClansInfo{
		1: model.WGClansInfoData{Tag: "TEST"},
	}, nil)

	mockUnofficialWargaming.On("ClansAutoComplete", mock.Anything).Return(model.UWGClansAutocomplete{}, nil)

	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockStorage.On("WriteOwnIGN", mock.Anything).Return(nil)

	mockLogger.On("SetOwnIGN", mock.Anything).Return()

	// テスト
	b := NewBattle(
		5,
		mockWargaming,
		mockUnofficialWargaming,
		mockLocalFile,
		mockNumbers,
		mockUnregistered,
		mockStorage,
		mockLogger,
		nil,
	)
	b.isFirstBattle = false
	_, err := b.Get(context.TODO(), testUserConfig)

	// アサーション
	require.NoError(t, err)
	mockWargaming.AssertExpectations(t)
	mockUnofficialWargaming.AssertExpectations(t)
	mockNumbers.AssertExpectations(t)
	mockUnregistered.AssertExpectations(t)
	mockLocalFile.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestBattle_Get_異常系(t *testing.T) {
	t.Parallel()

	// 準備
	mockWargaming := &mocks.WargamingInterface{}
	mockUnofficialWargaming := &mocks.UnofficialWargamingInterface{}
	mockNumbers := &mocks.NumbersInterface{}
	mockUnregistered := &mocks.UnregisteredInterface{}
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}
	mockLogger := &mocks.LoggerInterface{}

	expectedError := failure.New(apperr.FileNotExist)
	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(model.TempArenaInfo{}, expectedError)

	// テスト
	b := NewBattle(
		5,
		mockWargaming,
		mockUnofficialWargaming,
		mockLocalFile,
		mockNumbers,
		mockUnregistered,
		mockStorage,
		mockLogger,
		nil,
	)
	b.isFirstBattle = false
	_, err := b.Get(context.TODO(), testUserConfig)

	// アサーション
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, apperr.FileNotExist, code)

	mockWargaming.AssertExpectations(t)
	mockUnofficialWargaming.AssertExpectations(t)
	mockNumbers.AssertExpectations(t)
	mockUnregistered.AssertExpectations(t)
	mockLocalFile.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}
