package usecase

import (
	"testing"
	"wfs/backend/domain/mock_repository"
	"wfs/backend/domain/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConfigMigrator_Migrate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

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
		mockConfigV0 := mock_repository.NewMockConfigV0Interface(ctrl)
		mockConfigV0.EXPECT().IsExistUser().Return(true)
		mockConfigV0.EXPECT().User().Return(savedUserConfig, nil)
		mockConfigV0.EXPECT().DeleteUser().Return(nil)
		mockConfigV0.EXPECT().IsExistAlertPlayers().Return(true).MaxTimes(1)
		mockConfigV0.EXPECT().AlertPlayers().Return(savedAlertPlayers, nil)
		mockConfigV0.EXPECT().DeleteAlertPlayers().Return(nil)

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(0), nil).MaxTimes(1)
		mockStorage.EXPECT().IsExistUserConfig().Return(false)
		mockStorage.EXPECT().WriteUserConfig(savedUserConfig).Return(nil)
		mockStorage.EXPECT().IsExistAlertPlayers().Return(false).MaxTimes(1)
		mockStorage.EXPECT().WriteAlertPlayers(savedAlertPlayers).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(1)).Return(nil).MaxTimes(1)

		// toV2()
		mockStorage.EXPECT().DataVersion().Return(uint(1), nil).MaxTimes(1)
		mockStorage.EXPECT().UserConfigV2().Return(model.UserConfigV2{}, nil)
		mockStorage.EXPECT().UserConfig().Return(savedUserConfig, nil)
		mockStorage.EXPECT().WriteUserConfigV2(model.FromUserConfigV1(savedUserConfig)).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(2)).Return(nil).MaxTimes(1)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.ExecuteIfNeeded()

		// アサーション
		require.NoError(t, err)
	})
}

func TestConfigMigrator_toV1(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_マイグレ不要_バージョン1以上", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
	})

	t.Run("正常系_マイグレ不要_すでにストレージに存在", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockConfigV0 := mock_repository.NewMockConfigV0Interface(ctrl)
		mockConfigV0.EXPECT().IsExistUser().Return(true)
		mockConfigV0.EXPECT().IsExistAlertPlayers().Return(true)

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(0), nil)
		mockStorage.EXPECT().IsExistUserConfig().Return(true)
		mockStorage.EXPECT().IsExistAlertPlayers().Return(true)
		mockStorage.EXPECT().WriteDataVersion(uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
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
		mockConfigV0 := mock_repository.NewMockConfigV0Interface(ctrl)
		mockConfigV0.EXPECT().IsExistUser().Return(true)
		mockConfigV0.EXPECT().User().Return(savedUserConfig, nil)
		mockConfigV0.EXPECT().DeleteUser().Return(nil)
		mockConfigV0.EXPECT().IsExistAlertPlayers().Return(true)
		mockConfigV0.EXPECT().AlertPlayers().Return(savedAlertPlayers, nil)
		mockConfigV0.EXPECT().DeleteAlertPlayers().Return(nil)

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(0), nil)
		mockStorage.EXPECT().IsExistUserConfig().Return(false)
		mockStorage.EXPECT().WriteUserConfig(savedUserConfig).Return(nil)
		mockStorage.EXPECT().IsExistAlertPlayers().Return(false)
		mockStorage.EXPECT().WriteAlertPlayers(savedAlertPlayers).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
	})
}

func TestConfigMigrator_toV2(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_マイグレ不要_バージョン2以上", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(2), nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
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

		expect := model.UserConfigV2{
			Version:     2,
			Appid:       "test_appid",
			InstallPath: "test_install_path",
			Display: model.UCDisplay{
				Ship: model.UCDisplayShip{
					PR: true,
				},
			},
		}

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)
		mockStorage.EXPECT().UserConfigV2().Return(v2, nil)
		mockStorage.EXPECT().WriteUserConfigV2(expect).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
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

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)
		mockStorage.EXPECT().UserConfigV2().Return(model.UserConfigV2{}, nil)
		mockStorage.EXPECT().UserConfig().Return(v1, nil)
		mockStorage.EXPECT().WriteUserConfigV2(expected).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(nil, mockStorage, nil)
		err := cm.toV2()

		// アサーション
		require.NoError(t, err)
	})
}
