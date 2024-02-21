package usecase

import (
	"context"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock_repository"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//nolint:gochecknoglobals
var testUserConfig = model.UserConfigV2{
	Appid:             "test_appid",
	InstallPath:       "test_install_path",
	SaveTempArenaInfo: false,
}

func TestBattle_Get_正常系_初回(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	// 準備
	mockWargaming := mock_repository.NewMockWargamingInterface(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any(), gomock.Any()).Return(model.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.EXPECT().EncycShips(gomock.Any(), gomock.Any()).Return(model.WGEncycShips{
		1: model.WGEncycShipsData{
			Tier:      1,
			Type:      "Battleship",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, 2, nil).Times(2)
	mockWargaming.EXPECT().BattleArenas(gomock.Any()).Return(model.WGBattleArenas{}, nil)
	mockWargaming.EXPECT().BattleTypes(gomock.Any()).Return(model.WGBattleTypes{}, nil)
	mockWargaming.EXPECT().AccountInfo(gomock.Any(), gomock.Any()).Return(model.WGAccountInfo{}, nil)
	mockWargaming.EXPECT().ShipsStats(gomock.Any(), gomock.Any()).Return(model.WGShipsStats{}, nil).AnyTimes()
	mockWargaming.EXPECT().ClansAccountInfo(gomock.Any(), gomock.Any()).Return(model.WGClansAccountInfo{}, nil)
	mockWargaming.EXPECT().ClansInfo(gomock.Any(), gomock.Any()).Return(model.WGClansInfo{}, nil)

	mockUnofficialWargaming := mock_repository.NewMockUnofficialWargamingInterface(ctrl)
	mockUnofficialWargaming.EXPECT().ClansAutoComplete(gomock.Any()).Return(model.UWGClansAutocomplete{
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

	mockNumbers := mock_repository.NewMockNumbersInterface(ctrl)
	mockNumbers.EXPECT().ExpectedStats().Return(model.ExpectedStats{}, nil)

	mockUnregistered := mock_repository.NewMockUnregisteredInterface(ctrl)
	mockUnregistered.EXPECT().Warship().Return(model.Warships{}, nil)

	mockLocalFile := mock_repository.NewMockLocalFileInterface(ctrl)
	mockLocalFile.EXPECT().TempArenaInfo(gomock.Any()).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockStorage := mock_repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)
	mockStorage.EXPECT().WriteExpectedStats(gomock.Any()).Return(nil)

	mockLogger := mock_repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

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
}

func TestBattle_Get_正常系_2回目以降(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	// 準備
	mockWargaming := mock_repository.NewMockWargamingInterface(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any(), gomock.Any()).Return(model.WGAccountList{
		{NickName: "player_1", AccountID: 1},
		{NickName: "player_2", AccountID: 2},
	}, nil)
	mockWargaming.EXPECT().AccountInfo(gomock.Any(), gomock.Any()).Return(model.WGAccountInfo{}, nil)
	mockWargaming.EXPECT().ShipsStats(gomock.Any(), gomock.Any()).Return(model.WGShipsStats{}, nil).AnyTimes()
	mockWargaming.EXPECT().ClansAccountInfo(gomock.Any(), gomock.Any()).Return(model.WGClansAccountInfo{}, nil)
	mockWargaming.EXPECT().ClansInfo(gomock.Any(), gomock.Any()).Return(model.WGClansInfo{}, nil)

	mockUnofficialWargaming := mock_repository.NewMockUnofficialWargamingInterface(ctrl)
	mockUnofficialWargaming.EXPECT().ClansAutoComplete(gomock.Any()).Return(model.UWGClansAutocomplete{
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

	mockNumbers := mock_repository.NewMockNumbersInterface(ctrl)

	mockUnregistered := mock_repository.NewMockUnregisteredInterface(ctrl)

	mockLocalFile := mock_repository.NewMockLocalFileInterface(ctrl)
	mockLocalFile.EXPECT().TempArenaInfo(gomock.Any()).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
		PlayerName: "player_1",
	}, nil)

	mockStorage := mock_repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)

	mockLogger := mock_repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

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
}

func TestBattle_Get_異常系(t *testing.T) {
	t.Parallel()

	// 準備
	ctrl := gomock.NewController(t)

	// 準備
	mockLocalFile := mock_repository.NewMockLocalFileInterface(ctrl)
	expectedError := failure.New(apperr.FileNotExist)
	mockLocalFile.EXPECT().TempArenaInfo(gomock.Any()).Return(model.TempArenaInfo{}, expectedError)

	// テスト
	b := NewBattle(
		5,
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
