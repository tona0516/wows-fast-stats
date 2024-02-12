package usecase

import (
	"testing"
	"wfs/backend/domain/model"
	"wfs/backend/mocks"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
)

func TestConfigMigrator_Migrate(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		savedUserConfig := model.UserConfigV1{
			InstallPath: "a",
			Appid:       "a",
			FontSize:    "large",
		}
		savedAlertPlayers := []model.AlertPlayer{
			{
				AccountID: 1,
				Name:      "a",
				Pattern:   "bi-check-circle-fill",
				Message:   "a",
			},
		}
		mockConfigV0 := &mocks.ConfigV0Interface{}
		mockConfigV0.On("IsExistUser").Return(true)
		mockConfigV0.On("UserV1").Return(savedUserConfig, nil)
		mockConfigV0.On("DeleteUser").Return(nil)
		mockConfigV0.On("IsExistAlertPlayers").Return(true)
		mockConfigV0.On("AlertPlayers").Return(savedAlertPlayers, nil)
		mockConfigV0.On("DeleteAlertPlayers").Return(nil)
		mockStorage := &mocks.StorageInterface{}
		// toV1()
		mockStorage.On("DataVersion").Return(uint(0), nil)
		mockStorage.On("IsExistUserConfig").Return(false)
		mockStorage.On("WriteUserConfigV1", savedUserConfig).Return(nil)
		mockStorage.On("IsExistAlertPlayers").Return(false)
		mockStorage.On("WriteAlertPlayers", savedAlertPlayers).Return(nil)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil).Once()
		// toV2()
		mockStorage.On("DataVersion").Return(uint(1), nil)
		mockStorage.On("UserConfigV1").Return(savedUserConfig, nil)
		mockStorage.On("WriteUserConfig", model.FromUserConfigV1(savedUserConfig)).Return(nil)
		mockStorage.On("WriteDataVersion", uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.ExecuteIfNeeded()

		// アサーション
		require.NoError(t, err)
		mockConfigV0.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}

func TestConfigMigrator_toV1(t *testing.T) {
	t.Parallel()
	t.Run("正常系_マイグレ不要_バージョン1以上", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(1), nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
		mockStorage.AssertNotCalled(t, "WriteUserConfig")
		mockStorage.AssertNotCalled(t, "WriteAlertPlayers")
		mockStorage.AssertNotCalled(t, "WriteDataVersion")
	})
	t.Run("正常系_マイグレ不要_すでにストレージに存在", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockConfigV0 := &mocks.ConfigV0Interface{}
		mockConfigV0.On("IsExistUser").Return(true)
		mockConfigV0.On("IsExistAlertPlayers").Return(true)
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(0), nil)
		mockStorage.On("IsExistUserConfig").Return(true)
		mockStorage.On("IsExistAlertPlayers").Return(true)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
		mockConfigV0.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
		mockStorage.AssertNotCalled(t, "WriteUserConfig")
		mockStorage.AssertNotCalled(t, "WriteAlertPlayers")
	})
}

func TestConfigMigrator_toV2(t *testing.T) {
	t.Parallel()
	t.Run("正常系_マイグレ不要_バージョン2以上", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(2), nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
		mockStorage.AssertNotCalled(t, "WriteUserConfig")
		mockStorage.AssertNotCalled(t, "WriteDataVersion")
	})
	t.Run("正常系_マイグレ不要_v1が存在しない", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(1), nil)
		mockStorage.On("UserConfigV1").Return(model.UserConfigV1{}, badger.ErrKeyNotFound)
		mockStorage.On("WriteDataVersion", uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
		mockStorage.AssertNotCalled(t, "WriteUserConfig")
	})

	t.Run("正常系_マイグレ成功", func(t *testing.T) {
		t.Parallel()

		// 準備
		v1 := model.UserConfigV1{
			InstallPath: "test_install_path",
			Appid:       "test_appid",
			FontSize:    "test_font_size",
			Displays: model.Displays{
				Ship: model.Ship{
					PR: true,
				},
				Overall: model.Overall{
					PR: true,
				},
			},
			CustomColor: model.CustomColor{
				Skill: model.SkillColor{
					Text: model.SkillColorCode{
						Bad: "#000001",
					},
					Background: model.SkillColorCode{
						Bad: "#000002",
					},
				},
				Tier: model.TierColor{
					Own: model.TierColorCode{
						Low: "#000003",
					},
					Other: model.TierColorCode{
						Low: "#000004",
					},
				},
				ShipType: model.ShipTypeColor{
					Own: model.ShipTypeColorCode{
						SS: "#000005",
					},
					Other: model.ShipTypeColorCode{
						SS: "#000006",
					},
				},
				PlayerName: model.StatsPatternPvPSolo,
			},
			CustomDigit: model.CustomDigit{
				KdRate: 2,
			},
			TeamAverage: model.TeamAverage{
				MinShipBattles:    1,
				MinOverallBattles: 10,
			},
			SaveScreenshot:    true,
			SaveTempArenaInfo: true,
			SendReport:        false,
			NotifyUpdatable:   false,
			StatsPattern:      "pvp_solo",
		}

		v2 := model.FromUserConfigV1(v1)

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(1), nil)
		mockStorage.On("UserConfigV1").Return(v1, nil)
		mockStorage.On("WriteUserConfig", v2).Return(nil)
		mockStorage.On("WriteDataVersion", uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})
}
