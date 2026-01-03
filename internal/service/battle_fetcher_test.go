package service

import (
	"context"
	"errors"
	"testing"
	"wfs/internal/data"
	"wfs/internal/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBattle_Get_正常系_初回(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	// 準備
	mockWargaming := mock.NewMockWargamingApiClient(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any()).Return(data.WGAccountList{
		WGResponseCommon: data.WGResponseCommon[[]data.WGAccountListData]{
			Status: "",
			Error:  data.WGError{},
			Data: []data.WGAccountListData{
				{NickName: "player_1", AccountID: 1},
				{NickName: "player_2", AccountID: 2},
			},
		},
	}, nil)
	mockWargaming.EXPECT().EncycShips(gomock.Any()).Return(data.WGEncycShips{
		WGResponseCommon: data.WGResponseCommon[map[int]data.WGEncycShipsData]{
			Status: "",
			Error:  data.WGError{},
			Data: map[int]data.WGEncycShipsData{
				1: {
					Tier:      1,
					Type:      "Battleship",
					Name:      "ship_1",
					Nation:    "japan",
					IsPremium: false,
				},
			},
		},
		Meta: struct {
			PageTotal int "json:\"page_total\""
			Page      int "json:\"page\""
		}{
			PageTotal: 2,
			Page:      1,
		},
	}, nil).Times(2)
	mockWargaming.EXPECT().BattleArenas().Return(data.WGBattleArenas{}, nil)
	mockWargaming.EXPECT().BattleTypes().Return(data.WGBattleTypes{}, nil)
	mockWargaming.EXPECT().AccountInfo(gomock.Any()).Return(data.WGAccountInfo{}, nil)
	mockWargaming.EXPECT().ShipsStats(gomock.Any()).Return(data.WGShipsStats{}, nil).AnyTimes()
	mockWargaming.EXPECT().ShipsBadges(gomock.Any()).Return(data.WGShipsBadges{}, nil).AnyTimes()
	mockWargaming.EXPECT().ClansAccountInfo(gomock.Any()).Return(data.WGClansAccountInfo{}, nil)
	mockWargaming.EXPECT().ClansInfo(gomock.Any()).Return(data.WGClansInfo{}, nil)

	mockUnofficialWargaming := mock.NewMockClanApiClient(ctrl)
	mockUnofficialWargaming.EXPECT().ClanAutoComplete(gomock.Any()).Return(data.ClanAutocomplete{
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

	mockNumbers := mock.NewMockNumbersApiClient(ctrl)
	mockNumbers.EXPECT().ExpectedStats().Return(data.NSExpectedStats{}, nil)

	mockPersistence := mock.NewMockLocalStorage(ctrl)
	mockPersistence.EXPECT().SetOwnIGN(gomock.Any()).Return(nil)
	mockPersistence.EXPECT().SetExpectedStats(gomock.Any()).Return(nil)

	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	// テスト
	b := NewBattleFetcher(
		context.TODO(),
		mockPersistence,
		mockWargaming,
		mockUnofficialWargaming,
		mockNumbers,
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
	mockWargaming := mock.NewMockWargamingApiClient(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any()).Return(data.WGAccountList{
		WGResponseCommon: data.WGResponseCommon[[]data.WGAccountListData]{
			Status: "",
			Error:  data.WGError{},
			Data: []data.WGAccountListData{
				{NickName: "player_1", AccountID: 1},
				{NickName: "player_2", AccountID: 2},
			},
		},
	}, nil)
	mockWargaming.EXPECT().AccountInfo(gomock.Any()).Return(data.WGAccountInfo{}, nil)
	mockWargaming.EXPECT().ShipsStats(gomock.Any()).Return(data.WGShipsStats{}, nil).AnyTimes()
	mockWargaming.EXPECT().ShipsBadges(gomock.Any()).Return(data.WGShipsBadges{}, nil).AnyTimes()
	mockWargaming.EXPECT().ClansAccountInfo(gomock.Any()).Return(data.WGClansAccountInfo{}, nil)
	mockWargaming.EXPECT().ClansInfo(gomock.Any()).Return(data.WGClansInfo{}, nil)

	mockUnofficialWargaming := mock.NewMockClanApiClient(ctrl)
	mockUnofficialWargaming.EXPECT().ClanAutoComplete(gomock.Any()).Return(data.ClanAutocomplete{
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

	mockNumbers := mock.NewMockNumbersApiClient(ctrl)

	mockPersistence := mock.NewMockLocalStorage(ctrl)
	mockPersistence.EXPECT().SetOwnIGN(gomock.Any()).Return(nil)

	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	// テスト
	b := NewBattleFetcher(
		context.TODO(),
		mockPersistence,
		mockWargaming,
		mockUnofficialWargaming,
		mockNumbers,
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
	mockWargaming := mock.NewMockWargamingApiClient(ctrl)
	mockWargaming.EXPECT().AccountList(gomock.Any()).Return(data.WGAccountList{}, errors.New("hoge"))

	mockUnofficialWargaming := mock.NewMockClanApiClient(ctrl)

	mockNumbers := mock.NewMockNumbersApiClient(ctrl)

	mockPersistence := mock.NewMockLocalStorage(ctrl)
	mockPersistence.EXPECT().SetOwnIGN(gomock.Any()).Return(nil)

	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	// テスト
	b := NewBattleFetcher(
		context.TODO(),
		mockPersistence,
		mockWargaming,
		mockUnofficialWargaming,
		mockNumbers,
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
