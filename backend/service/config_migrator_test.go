package service

import (
	"testing"
	"wfs/backend/domain/mock/repository"
	"wfs/backend/domain/model"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

		ctrl := gomock.NewController(t)
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockStorage := repository.NewMockStorage(ctrl)

		// toV1()
		mockUserConfig.EXPECT().IsExistV0().Return(true)
		mockUserConfig.EXPECT().IsExistV1().Return(false)
		mockUserConfig.EXPECT().GetV0().Return(savedUserConfig, nil)
		mockUserConfig.EXPECT().DeleteV0().Return(nil)
		mockUserConfig.EXPECT().SaveV1(savedUserConfig).Return(nil)

		mockAlertPlayer.EXPECT().IsExistV0().Return(true)
		mockAlertPlayer.EXPECT().IsExistV1().Return(false)
		mockAlertPlayer.EXPECT().GetV0().Return(savedAlertPlayers, nil)
		mockAlertPlayer.EXPECT().DeleteV0().Return(nil)
		mockAlertPlayer.EXPECT().SaveV1(savedAlertPlayers).Return(nil)

		mockStorage.EXPECT().DataVersion().Return(uint(0), nil)
		mockStorage.EXPECT().WriteDataVersion(uint(1)).Return(nil)

		// toV2()
		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)
		mockUserConfig.EXPECT().GetV2().Return(model.UserConfigV2{}, nil)
		mockUserConfig.EXPECT().GetV1().Return(savedUserConfig, nil)
		mockUserConfig.EXPECT().SaveV2(model.FromUserConfigV1(savedUserConfig)).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, mockUserConfig, mockAlertPlayer)
		err := cm.ExecuteIfNeeded()

		// アサーション
		assert.NoError(t, err)
	})
}

func TestConfigMigrator_toV1(t *testing.T) {
	t.Parallel()

	t.Run("正常系_マイグレ不要_バージョン1以上", func(t *testing.T) {
		t.Parallel()

		// 準備
		ctrl := gomock.NewController(t)
		mockStorage := repository.NewMockStorage(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, nil, nil)
		err := cm.toV1()

		// アサーション
		assert.NoError(t, err)
	})

	t.Run("正常系_マイグレ不要_すでにストレージに存在", func(t *testing.T) {
		t.Parallel()

		// 準備
		ctrl := gomock.NewController(t)
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockStorage := repository.NewMockStorage(ctrl)

		mockUserConfig.EXPECT().IsExistV0().Return(true)
		mockUserConfig.EXPECT().IsExistV1().Return(true)
		mockAlertPlayer.EXPECT().IsExistV0().Return(true)
		mockAlertPlayer.EXPECT().IsExistV1().Return(true)

		mockStorage.EXPECT().DataVersion().Return(uint(0), nil)
		mockStorage.EXPECT().WriteDataVersion(uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, mockUserConfig, mockAlertPlayer)
		err := cm.toV1()

		// アサーション
		assert.NoError(t, err)
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

		ctrl := gomock.NewController(t)
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockStorage := repository.NewMockStorage(ctrl)

		// toV1()
		mockUserConfig.EXPECT().IsExistV0().Return(true)
		mockUserConfig.EXPECT().IsExistV1().Return(false)
		mockUserConfig.EXPECT().GetV0().Return(savedUserConfig, nil)
		mockUserConfig.EXPECT().DeleteV0().Return(nil)
		mockUserConfig.EXPECT().SaveV1(savedUserConfig).Return(nil)

		mockAlertPlayer.EXPECT().IsExistV0().Return(true)
		mockAlertPlayer.EXPECT().IsExistV1().Return(false)
		mockAlertPlayer.EXPECT().GetV0().Return(savedAlertPlayers, nil)
		mockAlertPlayer.EXPECT().DeleteV0().Return(nil)
		mockAlertPlayer.EXPECT().SaveV1(savedAlertPlayers).Return(nil)

		mockStorage.EXPECT().DataVersion().Return(uint(0), nil)
		mockStorage.EXPECT().WriteDataVersion(uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, mockUserConfig, mockAlertPlayer)
		err := cm.toV1()

		// アサーション
		assert.NoError(t, err)
	})
}

func TestConfigMigrator_toV2(t *testing.T) {
	t.Parallel()

	t.Run("正常系_マイグレ不要_バージョン2以上", func(t *testing.T) {
		t.Parallel()

		// 準備
		ctrl := gomock.NewController(t)
		mockStorage := repository.NewMockStorage(ctrl)
		mockStorage.EXPECT().DataVersion().Return(uint(2), nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, nil, nil)
		err := cm.toV2()

		// アサーション
		assert.NoError(t, err)
	})
	t.Run("正常系_UserConfigV2のVersionのみ更新", func(t *testing.T) {
		t.Parallel()

		// 準備
		v2 := model.UserConfigV2{
			InstallPath: "test_install_path",
			Display: model.UCDisplay{
				Ship: model.UCDisplayShip{
					PR: true,
				},
			},
		}

		expect := model.UserConfigV2{
			Version:     2,
			InstallPath: "test_install_path",
			Display: model.UCDisplay{
				Ship: model.UCDisplayShip{
					PR: true,
				},
			},
		}

		ctrl := gomock.NewController(t)
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockStorage := repository.NewMockStorage(ctrl)

		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)
		mockUserConfig.EXPECT().GetV2().Return(v2, nil)
		mockUserConfig.EXPECT().SaveV2(expect).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, mockUserConfig, nil)
		err := cm.toV2()

		// アサーション
		assert.NoError(t, err)
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

		ctrl := gomock.NewController(t)
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockStorage := repository.NewMockStorage(ctrl)

		mockStorage.EXPECT().DataVersion().Return(uint(1), nil)
		mockUserConfig.EXPECT().GetV2().Return(model.UserConfigV2{}, nil)
		mockUserConfig.EXPECT().GetV1().Return(v1, nil)
		mockUserConfig.EXPECT().SaveV2(expected).Return(nil)
		mockStorage.EXPECT().WriteDataVersion(uint(2)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockStorage, mockUserConfig, nil)
		err := cm.toV2()

		// アサーション
		assert.NoError(t, err)
	})
}
