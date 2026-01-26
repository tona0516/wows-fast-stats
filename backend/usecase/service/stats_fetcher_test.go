package service

import (
	"context"
	"testing"
	"wfs/backend/adapter"
	"wfs/backend/core"
	"wfs/backend/mock"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStatsFetcher_FetchAll(t *testing.T) {
	t.Parallel()

	t.Run("正常系_複数のアカウントIDに対して船統計を取得", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)

		accountID1 := core.AccountID(123)
		accountID2 := core.AccountID(456)
		accountIDs := []core.AccountID{accountID1, accountID2}

		// 船統計データの構築
		shipStats1 := core.WGShipsStatsData{
			ShipID: 1,
			Pvp: core.WGShipStatsValues{
				Battles: 100,
				Wins:    50,
			},
		}
		shipStats2 := core.WGShipsStatsData{
			ShipID: 2,
			Pvp: core.WGShipStatsValues{
				Battles: 50,
				Wins:    25,
			},
		}
		shipStats3 := core.WGShipsStatsData{
			ShipID: 3,
			Pvp: core.WGShipStatsValues{
				Battles: 75,
				Wins:    40,
			},
		}

		resp1 := core.WGShipsStats{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID][]core.WGShipsStatsData]{
				Data: map[core.AccountID][]core.WGShipsStatsData{
					accountID1: {shipStats1, shipStats2},
				},
			},
		}
		resp2 := core.WGShipsStats{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID][]core.WGShipsStatsData]{
				Data: map[core.AccountID][]core.WGShipsStatsData{
					accountID2: {shipStats3},
				},
			},
		}

		mockWargamingClient.EXPECT().
			ShipsStats(gomock.Any(), accountID1).
			Return(resp1, nil).
			Times(1)
		mockWargamingClient.EXPECT().
			ShipsStats(gomock.Any(), accountID2).
			Return(resp2, nil).
			Times(1)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewStatsFetcher(injector)
		require.NoError(t, err)

		result, err := svc.FetchAll(context.Background(), accountIDs)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("正常系_アカウントIDが空の場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewStatsFetcher(injector)
		require.NoError(t, err)

		result, err := svc.FetchAll(context.Background(), []core.AccountID{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("異常系_APIからエラーが返される", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)

		accountID1 := core.AccountID(123)
		accountIDs := []core.AccountID{accountID1}

		expectedErr := failure.New(core.ErrWGAPI)
		mockWargamingClient.EXPECT().
			ShipsStats(gomock.Any(), accountID1).
			Return(core.WGShipsStats{}, expectedErr)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewStatsFetcher(injector)
		require.NoError(t, err)

		result, err := svc.FetchAll(context.Background(), accountIDs)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
