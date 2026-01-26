package service

import (
	"context"
	"testing"
	"wfs/backend/core"
	"wfs/backend/mock"

	"github.com/abadojack/whatlanggo"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestClanFetcher_FetchAll(t *testing.T) {
	t.Parallel()

	t.Run("正常系_全データが正常に取得できる場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)
		mockClanClient := mock.NewMockClanClient(ctrl)

		accountID1 := core.AccountID(123)
		accountID2 := core.AccountID(456)
		accountID3 := core.AccountID(789)
		accountIDs := []core.AccountID{accountID1, accountID2, accountID3}

		clanID1 := core.ClanID(1001)
		clanID2 := core.ClanID(1002)
		clanID3 := core.ClanID(1003)

		// ClansAccountInfo の期待値
		clansAccountInfo := core.WGClansAccountInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID]core.WGClansAccountInfoData]{
				Data: map[core.AccountID]core.WGClansAccountInfoData{
					accountID1: {ClanID: clanID1},
					accountID2: {ClanID: clanID2},
					accountID3: {ClanID: clanID3},
				},
			},
		}

		// ClansInfo の期待値 (言語検出用に日本語の説明を含める)
		clansInfo := core.WGClansInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.ClanID]core.WGClansInfoData]{
				Data: map[core.ClanID]core.WGClansInfoData{
					clanID1: {Tag: "TAG1", Description: "こんにちは"},
					clanID2: {Tag: "TAG2", Description: "你好"},
					clanID3: {Tag: "TAG3", Description: "안녕하세요"},
				},
			},
		}

		// ClanAutocomplete の期待値
		autocomplete1 := core.ClanAutocomplete{
			SearchAutocompleteResult: []struct {
				HexColor string      `json:"hex_color"`
				ID       core.ClanID `json:"id"`
			}{
				{HexColor: "#FF0000", ID: clanID1},
			},
		}
		autocomplete2 := core.ClanAutocomplete{
			SearchAutocompleteResult: []struct {
				HexColor string      `json:"hex_color"`
				ID       core.ClanID `json:"id"`
			}{
				{HexColor: "#00FF00", ID: clanID2},
			},
		}
		autocomplete3 := core.ClanAutocomplete{
			SearchAutocompleteResult: []struct {
				HexColor string      `json:"hex_color"`
				ID       core.ClanID `json:"id"`
			}{
				{HexColor: "#0000FF", ID: clanID3},
			},
		}

		// モックの設定
		gomock.InOrder(
			mockWargamingClient.EXPECT().
				ClansAccountInfo(gomock.Any(), accountIDs).
				Return(clansAccountInfo, nil),
			mockWargamingClient.EXPECT().
				ClansInfo(gomock.Any(), gomock.Any()).
				Return(clansInfo, nil),
		)

		// fetchClanColor で使用されるモック
		mockClanClient.EXPECT().
			ClanAutoComplete(gomock.Any(), "TAG1").
			Return(autocomplete1, nil)
		mockClanClient.EXPECT().
			ClanAutoComplete(gomock.Any(), "TAG2").
			Return(autocomplete2, nil)
		mockClanClient.EXPECT().
			ClanAutoComplete(gomock.Any(), "TAG3").
			Return(autocomplete3, nil)

		// サービスの初期化
		injector := do.New()
		do.Provide(injector, func(i do.Injector) (interface{}, error) {
			return mockWargamingClient, nil
		})
		do.ProvideValue(injector, mockClanClient)

		service := &ClanFetcher{
			wargamingClient: mockWargamingClient,
			clanClient:      mockClanClient,
		}

		// テスト実行
		result, err := service.FetchAll(context.Background(), accountIDs)

		// アサーション
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 3, len(result))

		// accountID1 のデータ確認
		clan1, ok := result[accountID1]
		assert.True(t, ok)
		assert.Equal(t, clanID1, clan1.ID)
		assert.Equal(t, "TAG1", clan1.Tag)
		assert.Equal(t, "#FF0000", clan1.ColorCode)
		assert.Equal(t, whatlanggo.Jpn.Iso6391(), clan1.Language)

		// accountID2 のデータ確認
		clan2, ok := result[accountID2]
		assert.True(t, ok)
		assert.Equal(t, clanID2, clan2.ID)
		assert.Equal(t, "TAG2", clan2.Tag)
		assert.Equal(t, "#00FF00", clan2.ColorCode)
		assert.Equal(t, whatlanggo.Cmn.Iso6391(), clan2.Language)

		// accountID3 のデータ確認
		clan3, ok := result[accountID3]
		assert.True(t, ok)
		assert.Equal(t, clanID3, clan3.ID)
		assert.Equal(t, "TAG3", clan3.Tag)
		assert.Equal(t, "#0000FF", clan3.ColorCode)
		assert.Equal(t, whatlanggo.Kor.Iso6391(), clan3.Language)
	})

	t.Run("異常系_ClansAccountInfo でエラーが発生した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)
		mockClanClient := mock.NewMockClanClient(ctrl)

		accountID := core.AccountID(123)
		accountIDs := []core.AccountID{accountID}

		mockWargamingClient.EXPECT().
			ClansAccountInfo(gomock.Any(), accountIDs).
			Return(core.WGClansAccountInfo{}, assert.AnError)

		service := &ClanFetcher{
			wargamingClient: mockWargamingClient,
			clanClient:      mockClanClient,
		}

		// テスト実行
		result, err := service.FetchAll(context.Background(), accountIDs)

		// アサーション
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("異常系_ClansInfo でエラーが発生した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)
		mockClanClient := mock.NewMockClanClient(ctrl)

		accountID := core.AccountID(123)
		accountIDs := []core.AccountID{accountID}
		clanID := core.ClanID(1001)

		clansAccountInfo := core.WGClansAccountInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID]core.WGClansAccountInfoData]{
				Data: map[core.AccountID]core.WGClansAccountInfoData{
					accountID: {ClanID: clanID},
				},
			},
		}

		gomock.InOrder(
			mockWargamingClient.EXPECT().
				ClansAccountInfo(gomock.Any(), accountIDs).
				Return(clansAccountInfo, nil),
			mockWargamingClient.EXPECT().
				ClansInfo(gomock.Any(), gomock.Any()).
				Return(core.WGClansInfo{}, assert.AnError),
		)

		service := &ClanFetcher{
			wargamingClient: mockWargamingClient,
			clanClient:      mockClanClient,
		}

		// テスト実行
		result, err := service.FetchAll(context.Background(), accountIDs)

		// アサーション
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("異常系_fetchClanColor でエラーが発生した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)
		mockClanClient := mock.NewMockClanClient(ctrl)

		accountID := core.AccountID(123)
		accountIDs := []core.AccountID{accountID}
		clanID := core.ClanID(1001)

		clansAccountInfo := core.WGClansAccountInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID]core.WGClansAccountInfoData]{
				Data: map[core.AccountID]core.WGClansAccountInfoData{
					accountID: {ClanID: clanID},
				},
			},
		}

		clansInfo := core.WGClansInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.ClanID]core.WGClansInfoData]{
				Data: map[core.ClanID]core.WGClansInfoData{
					clanID: {Tag: "TAG1", Description: "Clan 1 Description"},
				},
			},
		}

		gomock.InOrder(
			mockWargamingClient.EXPECT().
				ClansAccountInfo(gomock.Any(), accountIDs).
				Return(clansAccountInfo, nil),
			mockWargamingClient.EXPECT().
				ClansInfo(gomock.Any(), gomock.Any()).
				Return(clansInfo, nil),
		)

		mockClanClient.EXPECT().
			ClanAutoComplete(gomock.Any(), "TAG1").
			Return(core.ClanAutocomplete{}, assert.AnError)

		service := &ClanFetcher{
			wargamingClient: mockWargamingClient,
			clanClient:      mockClanClient,
		}

		// テスト実行
		result, err := service.FetchAll(context.Background(), accountIDs)

		// アサーション
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("正常系_クランに属していないアカウントを含む場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)
		mockClanClient := mock.NewMockClanClient(ctrl)

		accountID1 := core.AccountID(123)
		accountID2 := core.AccountID(456)
		accountIDs := []core.AccountID{accountID1, accountID2}

		clanID1 := core.ClanID(1001)

		// accountID2 はクランに属していない (ClanID = 0)
		clansAccountInfo := core.WGClansAccountInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID]core.WGClansAccountInfoData]{
				Data: map[core.AccountID]core.WGClansAccountInfoData{
					accountID1: {ClanID: clanID1},
					accountID2: {ClanID: 0},
				},
			},
		}

		clansInfo := core.WGClansInfo{
			WGResponseCommon: core.WGResponseCommon[map[core.ClanID]core.WGClansInfoData]{
				Data: map[core.ClanID]core.WGClansInfoData{
					clanID1: {Tag: "TAG1", Description: "Clan 1 Description"},
				},
			},
		}

		autocomplete1 := core.ClanAutocomplete{
			SearchAutocompleteResult: []struct {
				HexColor string      `json:"hex_color"`
				ID       core.ClanID `json:"id"`
			}{
				{HexColor: "#FF0000", ID: clanID1},
			},
		}

		gomock.InOrder(
			mockWargamingClient.EXPECT().
				ClansAccountInfo(gomock.Any(), accountIDs).
				Return(clansAccountInfo, nil),
			mockWargamingClient.EXPECT().
				ClansInfo(gomock.Any(), gomock.Any()).
				Return(clansInfo, nil),
		)

		mockClanClient.EXPECT().
			ClanAutoComplete(gomock.Any(), "TAG1").
			Return(autocomplete1, nil)

		service := &ClanFetcher{
			wargamingClient: mockWargamingClient,
			clanClient:      mockClanClient,
		}

		// テスト実行
		result, err := service.FetchAll(context.Background(), accountIDs)

		// アサーション
		require.NoError(t, err)
		assert.NotNil(t, result)

		// accountID1 のみ結果に含まれることを確認
		clan1, ok := result[accountID1]
		assert.True(t, ok)
		assert.Equal(t, clanID1, clan1.ID)

		// accountID2 も結果に含まれることを確認 (ClanID = 0)
		clan2, ok := result[accountID2]
		assert.True(t, ok)
		assert.Equal(t, core.ClanID(0), clan2.ID)
	})
}
