package service

import (
	"context"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"
	domainRepository "wfs/backend/domain/mock/repository"
	"wfs/backend/domain/model"
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

	mockWarshipFetcher := domainRepository.NewMockWarshipFetcherInterface(ctrl)
	mockWarshipFetcher.EXPECT().Fetch().Return(model.Warships{
		1: {
			Tier:      1,
			Type:      "bb",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, nil)

	mockClanFetcher := domainRepository.NewMockClanFetcherInterface(ctrl)
	mockClanFetcher.EXPECT().Fetch(gomock.Any()).Return(model.Clans{
		1: {
			ID:       1919810,
			Tag:      "TEST",
			HexColor: "#114514",
		},
	}, nil)

	mockTAIFetcher := domainRepository.NewMockTAIFetcherInterface(ctrl)
	mockTAIFetcher.EXPECT().Get(gomock.Any()).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
	}, nil)

	mockRawStatFetcher := domainRepository.NewMockRawStatFetcherInterface(ctrl)
	mockRawStatFetcher.EXPECT().Fetch(gomock.Any()).Return(model.RawStats{}, nil)

	mockBattleMetaFetcher := domainRepository.NewMockBattleMetaFetcherInterface(ctrl)
	mockBattleMetaFetcher.EXPECT().Fetch().Return(model.BattleMeta{}, nil)

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// テスト
	b := NewBattle(
		mockWargaming,
		nil,
		mockWarshipFetcher,
		mockClanFetcher,
		mockTAIFetcher,
		mockRawStatFetcher,
		mockBattleMetaFetcher,
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

	mockWarshipFetcher := domainRepository.NewMockWarshipFetcherInterface(ctrl)
	mockWarshipFetcher.EXPECT().Fetch().Return(model.Warships{
		1: {
			Tier:      1,
			Type:      "Battleship",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, nil)

	mockClanFetcher := domainRepository.NewMockClanFetcherInterface(ctrl)
	mockClanFetcher.EXPECT().Fetch(gomock.Any()).Return(model.Clans{
		1: {
			ID:       1919810,
			Tag:      "TEST",
			HexColor: "#114514",
		},
	}, nil)

	mockTAIFetcher := domainRepository.NewMockTAIFetcherInterface(ctrl)
	mockTAIFetcher.EXPECT().Get(gomock.Any()).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
		PlayerName: "player_1",
	}, nil)

	mockRawStatFetcher := domainRepository.NewMockRawStatFetcherInterface(ctrl)
	mockRawStatFetcher.EXPECT().Fetch(gomock.Any()).Return(model.RawStats{}, nil)

	mockBattleMetaFetcher := domainRepository.NewMockBattleMetaFetcherInterface(ctrl)
	mockBattleMetaFetcher.EXPECT().Fetch().Return(model.BattleMeta{}, nil)

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockStorage.EXPECT().WriteOwnIGN(gomock.Any()).Return(nil)

	mockLogger := repository.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// テスト
	b := NewBattle(
		mockWargaming,
		nil,
		mockWarshipFetcher,
		mockClanFetcher,
		mockTAIFetcher,
		mockRawStatFetcher,
		mockBattleMetaFetcher,
		mockStorage,
		mockLogger,
		nil,
	)
	_, err := b.Get(context.TODO(), testUserConfig)

	// アサーション
	assert.NoError(t, err)
}

func TestBattle_Get_異常系(t *testing.T) {
	t.Parallel()

	// 準備
	ctrl := gomock.NewController(t)

	// 準備
	expectedError := failure.New(apperr.FileNotExist)
	mockTAIFetcher := domainRepository.NewMockTAIFetcherInterface(ctrl)
	mockTAIFetcher.EXPECT().Get(gomock.Any()).Return(model.TempArenaInfo{}, expectedError)

	// テスト
	b := NewBattle(
		nil,
		nil,
		nil,
		nil,
		mockTAIFetcher,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	_, err := b.Get(context.TODO(), testUserConfig)

	// アサーション
	assert.True(t, failure.Is(err, apperr.FileNotExist))
}
