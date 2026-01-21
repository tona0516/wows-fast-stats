package usecase

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

func TestBadgeService_FetchAll(t *testing.T) {
	t.Parallel()

	t.Run("正常系_複数のアカウントIDに対してバッジを取得", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)

		accountID1 := core.AccountID(123)
		accountID2 := core.AccountID(456)
		accountIDs := []core.AccountID{accountID1, accountID2}

		// shipID1 のバッジデータ
		badge1 := core.WGShipsBadgesData{
			ShipID: 1,
		}
		// shipID2 のバッジデータ
		badge2 := core.WGShipsBadgesData{
			ShipID: 2,
		}
		// shipID3 のバッジデータ
		badge3 := core.WGShipsBadgesData{
			ShipID: 3,
		}

		resp1 := core.WGShipsBadges{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID][]core.WGShipsBadgesData]{
				Data: map[core.AccountID][]core.WGShipsBadgesData{
					accountID1: {badge1, badge2},
				},
			},
		}
		resp2 := core.WGShipsBadges{
			WGResponseCommon: core.WGResponseCommon[map[core.AccountID][]core.WGShipsBadgesData]{
				Data: map[core.AccountID][]core.WGShipsBadgesData{
					accountID2: {badge3},
				},
			},
		}

		mockWargamingClient.EXPECT().
			ShipsBadges(gomock.Any(), accountID1).
			Return(resp1, nil).
			Times(1)
		mockWargamingClient.EXPECT().
			ShipsBadges(gomock.Any(), accountID2).
			Return(resp2, nil).
			Times(1)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewBadgeService(injector)
		require.NoError(t, err)

		result, err := svc.fetchAll(context.Background(), accountIDs)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		// アカウント1のバッジ確認
		badges1, ok := result[accountID1]
		assert.True(t, ok)
		assert.Len(t, badges1, 2)
		assert.Equal(t, core.ShipID(1), badges1[1].ShipID)
		assert.Equal(t, core.ShipID(2), badges1[2].ShipID)

		// アカウント2のバッジ確認
		badges2, ok := result[accountID2]
		assert.True(t, ok)
		assert.Len(t, badges2, 1)
		assert.Equal(t, core.ShipID(3), badges2[3].ShipID)
	})

	t.Run("正常系_アカウントIDが空の場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewBadgeService(injector)
		require.NoError(t, err)

		result, err := svc.fetchAll(context.Background(), []core.AccountID{})

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
			ShipsBadges(gomock.Any(), accountID1).
			Return(core.WGShipsBadges{}, expectedErr)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewBadgeService(injector)
		require.NoError(t, err)

		result, err := svc.fetchAll(context.Background(), accountIDs)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("異常系_複数アカウントのうち1つでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockWargamingClient := mock.NewMockWargamingClient(ctrl)

		accountID1 := core.AccountID(123)
		accountID2 := core.AccountID(456)
		accountIDs := []core.AccountID{accountID1, accountID2}

		mockWargamingClient.EXPECT().
			ShipsBadges(gomock.Any(), gomock.Any()).
			Return(core.WGShipsBadges{}, failure.New(core.ErrWGAPI)).
			AnyTimes()

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
			return mockWargamingClient, nil
		})
		svc, err := NewBadgeService(injector)
		require.NoError(t, err)

		result, err := svc.fetchAll(context.Background(), accountIDs)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
