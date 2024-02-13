package usecase

import (
	"testing"
	"wfs/backend/domain/model"
	"wfs/backend/mocks"

	"github.com/stretchr/testify/require"
)

func TestConfigMigrator_Migrate(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		savedUserConfig := model.UserConfig{
			InstallPath: "a",
			Appid:       "a",
			FontSize:    "large",
			Displays: model.Displays{
				Ship: model.Ship{PR: true},
			},
		}
		savedAlertPlayers := []model.AlertPlayer{
			{
				AccountID: 1,
				Name:      "a",
				Pattern:   "bi-check-circle-fill",
				Message:   "a",
			},
		}
		// toV1()
		mockConfigV0 := &mocks.ConfigV0Interface{}
		mockConfigV0.On("IsExistUser").Return(true)
		mockConfigV0.On("User").Return(savedUserConfig, nil)
		mockConfigV0.On("DeleteUser").Return(nil)
		mockConfigV0.On("IsExistAlertPlayers").Return(true).Once()
		mockConfigV0.On("AlertPlayers").Return(savedAlertPlayers, nil)
		mockConfigV0.On("DeleteAlertPlayers").Return(nil)
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(0), nil).Once()
		mockStorage.On("IsExistUserConfig").Return(false)
		mockStorage.On("WriteUserConfig", savedUserConfig).Return(nil)
		mockStorage.On("IsExistAlertPlayers").Return(false).Once()
		mockStorage.On("WriteAlertPlayers", savedAlertPlayers).Return(nil)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil).Once()
		// toV2()
		mockStorage.On("DataVersion").Return(uint(1), nil).Once()
		mockStorage.On("UserConfigV2").Return(model.UserConfigV2{}, nil)
		mockStorage.On("UserConfig").Return(savedUserConfig, nil)
		mockStorage.On("WriteUserConfigV2", model.FromUserConfigV1(savedUserConfig)).Return(nil)
		mockStorage.On("WriteDataVersion", uint(2)).Return(nil).Once()

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
	t.Run("正常系_成功", func(t *testing.T) {
		t.Parallel()

		// 準備
		savedUserConfig := model.UserConfig{
			InstallPath: "a",
			Appid:       "a",
			FontSize:    "large",
			Displays: model.Displays{
				Ship: model.Ship{PR: true},
			},
		}
		savedAlertPlayers := []model.AlertPlayer{
			{
				AccountID: 1,
				Name:      "a",
				Pattern:   "bi-check-circle-fill",
				Message:   "a",
			},
		}
		// toV1()
		mockConfigV0 := &mocks.ConfigV0Interface{}
		mockConfigV0.On("IsExistUser").Return(true)
		mockConfigV0.On("User").Return(savedUserConfig, nil)
		mockConfigV0.On("DeleteUser").Return(nil)
		mockConfigV0.On("IsExistAlertPlayers").Return(true)
		mockConfigV0.On("AlertPlayers").Return(savedAlertPlayers, nil)
		mockConfigV0.On("DeleteAlertPlayers").Return(nil)
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(0), nil)
		mockStorage.On("IsExistUserConfig").Return(false)
		mockStorage.On("WriteUserConfig", savedUserConfig).Return(nil)
		mockStorage.On("IsExistAlertPlayers").Return(false)
		mockStorage.On("WriteAlertPlayers", savedAlertPlayers).Return(nil)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
		mockConfigV0.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
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
		mockStorage.AssertNotCalled(t, "WriteUserConfigV2")
		mockStorage.AssertNotCalled(t, "WriteDataVersion")
	})

	t.Run("正常系_UserConfigV2のVersionのみ更新", func(t *testing.T) {
		t.Parallel()

		// 準備
		v2 := model.UserConfigV2{
			Appid:       "test_appid",
			InstallPath: "test_install_path",
			Display: model.UCDisplay{
				Ship: model.UCDisplayShip{
					PR: true,
				},
			},
		}

		expeect := model.UserConfigV2{
			Version:     2,
			Appid:       "test_appid",
			InstallPath: "test_install_path",
			Display: model.UCDisplay{
				Ship: model.UCDisplayShip{
					PR: true,
				},
			},
		}

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(1), nil)
		mockStorage.On("UserConfigV2").Return(v2, nil)
		mockStorage.On("WriteUserConfigV2", expeect).Return(nil)
		mockStorage.On("WriteDataVersion", uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("正常系_成功", func(t *testing.T) {
		t.Parallel()

		// 準備
		v1 := model.UserConfig{
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

		expected := model.FromUserConfigV1(v1)

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("DataVersion").Return(uint(1), nil)
		mockStorage.On("UserConfigV2").Return(model.UserConfigV2{}, nil)
		mockStorage.On("UserConfig").Return(v1, nil)
		mockStorage.On("WriteUserConfigV2", expected).Return(nil)
		mockStorage.On("WriteDataVersion", uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})
}
