package service

import (
	"context"
	"errors"
	"testing"
	"wfs/backend/data"
	"wfs/backend/mock/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

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

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)
	mockStorage.EXPECT().WriteExpectedStats(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	// テスト
	b := NewBattleFetcher(
		context.TODO(),
		mockWargaming,
		mockUnofficialWargaming,
		mockNumbers,
		mockUnregistered,
		mockStorage,
		mockLogger,
		emitFunc,
	)
	b.Invoke(data.TempArenaInfo{
		Vehicles: []data.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	})

	assert.Equal(t, 3, len(events))
	assert.Equal(t, EventFetchOthers, events[0])
	assert.Equal(t, EventFetchPlayers, events[1])
	assert.Equal(t, EventFetchDone, events[2])
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

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	// テスト
	b := NewBattleFetcher(
		context.TODO(),
		mockWargaming,
		mockUnofficialWargaming,
		mockNumbers,
		mockUnregistered,
		mockStorage,
		mockLogger,
		emitFunc,
	)
	b.isFirstBattle = false
	b.Invoke(data.TempArenaInfo{
		Vehicles: []data.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
		PlayerName: "player_1",
	})

	assert.Equal(t, 2, len(events))
	assert.Equal(t, EventFetchPlayers, events[0])
	assert.Equal(t, EventFetchDone, events[1])
}

func TestBattle_Get_異常系_アカウントリスト取得失敗(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	// 準備
	mockWargaming := repository.NewMockWargamingInterface(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any()).Return(nil, errors.New("hoge"))

	mockUnofficialWargaming := repository.NewMockUnofficialWargamingInterface(ctrl)

	mockNumbers := repository.NewMockNumbersInterface(ctrl)

	mockUnregistered := repository.NewMockUnregisteredInterface(ctrl)

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	// テスト
	b := NewBattleFetcher(
		context.TODO(),
		mockWargaming,
		mockUnofficialWargaming,
		mockNumbers,
		mockUnregistered,
		mockStorage,
		mockLogger,
		emitFunc,
	)
	b.isFirstBattle = false
	b.Invoke(data.TempArenaInfo{
		Vehicles: []data.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
		PlayerName: "player_1",
	})

	assert.Equal(t, 1, len(events))
	assert.Equal(t, EventErr, events[0])
}
