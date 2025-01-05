package service

import (
	"context"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/mock/repository"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//nolint:gochecknoglobals
var testUserConfig = data.UserConfigV2{
	InstallPath:       "test_install_path",
	SaveTempArenaInfo: false,
}

func TestBattle_Get_正常系_初回(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	// 準備
	mockWargaming := repository.NewMockWargamingInterface(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any()).Return(data.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.EXPECT().EncycShips(gomock.Any()).Return(data.WGEncycShips{
		1: data.WGEncycShipsData{
			Tier:      1,
			Type:      "Battleship",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, 2, nil).Times(2)
	mockWargaming.EXPECT().BattleArenas().Return(data.WGBattleArenas{}, nil)
	mockWargaming.EXPECT().BattleTypes().Return(data.WGBattleTypes{}, nil)
	mockWargaming.EXPECT().AccountInfo(gomock.Any()).Return(data.WGAccountInfo{}, nil)
	mockWargaming.EXPECT().ShipsStats(gomock.Any()).Return(data.WGShipsStats{}, nil).AnyTimes()
	mockWargaming.EXPECT().ClansAccountInfo(gomock.Any()).Return(data.WGClansAccountInfo{}, nil)
	mockWargaming.EXPECT().ClansInfo(gomock.Any()).Return(data.WGClansInfo{}, nil)

	mockUnofficialWargaming := repository.NewMockUnofficialWargamingInterface(ctrl)
	mockUnofficialWargaming.EXPECT().ClansAutoComplete(gomock.Any()).Return(data.UWGClansAutocomplete{
		SearchAutocompleteResult: []struct {
			HexColor string `json:"hex_color"`
			Tag      string `json:"tag"`
			ID       int    `json:"id"`
		}{
			{
				HexColor: "#114514",
				Tag:      "TEST",
				ID:       1919810,
			},
		},
	}, nil).AnyTimes()

	mockNumbers := repository.NewMockNumbersInterface(ctrl)
	mockNumbers.EXPECT().ExpectedStats().Return(data.ExpectedStats{}, nil)

	mockUnregistered := repository.NewMockUnregisteredInterface(ctrl)
	mockUnregistered.EXPECT().Warship().Return(data.Warships{}, nil)

	mockLocalFile := repository.NewMockLocalFileInterface(ctrl)
	mockLocalFile.EXPECT().TempArenaInfo(gomock.Any()).Return(data.TempArenaInfo{
		Vehicles: []data.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)
	mockStorage.EXPECT().WriteExpectedStats(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// テスト
	b := NewBattle(
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
	assert.NoError(t, err)
}

func TestBattle_Get_正常系_2回目以降(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	// 準備
	mockWargaming := repository.NewMockWargamingInterface(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any()).Return(data.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.EXPECT().AccountInfo(gomock.Any()).Return(data.WGAccountInfo{}, nil)
	mockWargaming.EXPECT().ShipsStats(gomock.Any()).Return(data.WGShipsStats{}, nil).AnyTimes()
	mockWargaming.EXPECT().ClansAccountInfo(gomock.Any()).Return(data.WGClansAccountInfo{}, nil)
	mockWargaming.EXPECT().ClansInfo(gomock.Any()).Return(data.WGClansInfo{}, nil)

	mockUnofficialWargaming := repository.NewMockUnofficialWargamingInterface(ctrl)
	mockUnofficialWargaming.EXPECT().ClansAutoComplete(gomock.Any()).Return(data.UWGClansAutocomplete{
		SearchAutocompleteResult: []struct {
			HexColor string `json:"hex_color"`
			Tag      string `json:"tag"`
			ID       int    `json:"id"`
		}{
			{
				HexColor: "#114514",
				Tag:      "TEST",
				ID:       1919810,
			},
		},
	}, nil).AnyTimes()

	mockNumbers := repository.NewMockNumbersInterface(ctrl)

	mockUnregistered := repository.NewMockUnregisteredInterface(ctrl)

	mockLocalFile := repository.NewMockLocalFileInterface(ctrl)
	mockLocalFile.EXPECT().TempArenaInfo(gomock.Any()).Return(data.TempArenaInfo{
		Vehicles: []data.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
		PlayerName: "player_1",
	}, nil)

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// テスト
	b := NewBattle(
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
	assert.NoError(t, err)
}

func TestBattle_Get_異常系(t *testing.T) {
	t.Parallel()

	// 準備
	ctrl := gomock.NewController(t)

	// 準備
	mockLocalFile := repository.NewMockLocalFileInterface(ctrl)
	expectedError := failure.New(apperr.FileNotExist)
	mockLocalFile.EXPECT().TempArenaInfo(gomock.Any()).Return(data.TempArenaInfo{}, expectedError)

	// テスト
	b := NewBattle(
		nil,
		nil,
		mockLocalFile,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	b.isFirstBattle = false
	_, err := b.Get(context.TODO(), testUserConfig)

	// アサーション
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, apperr.FileNotExist, code)
}
