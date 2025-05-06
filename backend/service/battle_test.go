package service

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock/repository"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//nolint:gochecknoglobals
var testUserConfig = model.UserConfigV2{
	InstallPath:       "test_install_path",
	SaveTempArenaInfo: false,
}

func TestBattle_Get_正常系(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	// 準備
	mockLocalFile := repository.NewMockLocalFile(ctrl)
	mockLocalFile.EXPECT().ReadTempArenaInfo(testUserConfig.InstallPath).Return(model.TempArenaInfo{
		Vehicles: []model.Vehicle{
			{ShipID: 1, Name: "player_1", Relation: 0},
			{ShipID: 2, Name: "player_2", Relation: 2},
		},
		PlayerName: "player_1",
	}, nil)

	mockWarshipStore := repository.NewMockWarshipFetcher(ctrl)
	mockWarshipStore.EXPECT().Fetch().Return(model.Warships{
		1: {
			Tier:      1,
			Type:      "bb",
			Name:      "ship_1",
			Nation:    "japan",
			IsPremium: false,
		},
	}, nil)

	mockClanFetcher := repository.NewMockClanFetcher(ctrl)
	mockClanFetcher.EXPECT().Fetch(gomock.Any()).Return(model.Clans{
		1: {
			ID:       1919810,
			Tag:      "TEST",
			HexColor: "#114514",
		},
	}, nil)

	mockRawStatFetcher := repository.NewMockRawStatFetcher(ctrl)
	mockRawStatFetcher.EXPECT().Fetch(gomock.Any()).Return(model.RawStats{}, nil)

	mockBattleMetaFetcher := repository.NewMockBattleMetaFetcher(ctrl)
	mockBattleMetaFetcher.EXPECT().Fetch().Return(model.BattleMeta{}, nil)

	mockAccountFetcher := repository.NewMockAccountFetcher(ctrl)
	mockAccountFetcher.EXPECT().FetchByNames(gomock.Any()).Return(model.Accounts{
		"player_1": 1,
		"player_2": 2,
	}, nil)

	mockLogger := repository.NewMockLogger(ctrl)
	mockLogger.EXPECT().SetOwnIGN(gomock.Any()).Return()

	// テスト
	b := NewBattle(
		mockLocalFile,
		mockWarshipStore,
		mockClanFetcher,
		mockRawStatFetcher,
		mockBattleMetaFetcher,
		mockAccountFetcher,
		mockLogger,
	)
	_, err := b.Get(t.Context(), testUserConfig)

	// アサーション
	assert.NoError(t, err)
}

func TestBattle_Get_異常系(t *testing.T) {
	t.Parallel()

	// 準備
	ctrl := gomock.NewController(t)

	// 準備
	expectedError := failure.New(apperr.FileNotExist)
	mockLocalFile := repository.NewMockLocalFile(ctrl)
	mockLocalFile.EXPECT().ReadTempArenaInfo(testUserConfig.InstallPath).Return(model.TempArenaInfo{}, expectedError)

	// テスト
	b := NewBattle(
		mockLocalFile,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	_, err := b.Get(t.Context(), testUserConfig)

	// アサーション
	assert.True(t, failure.Is(err, apperr.FileNotExist))
}
