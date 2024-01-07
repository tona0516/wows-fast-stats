package usecase

import (
	"testing"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/stretchr/testify/require"
)

func TestConfigMigrator_Migrate(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		expectedUserConfig := domain.UserConfig{
			InstallPath: "a",
			Appid:       "a",
			FontSize:    "large",
		}
		expectedAlertPlayers := []domain.AlertPlayer{
			{
				AccountID: 1,
				Name:      "a",
				Pattern:   "bi-check-circle-fill",
				Message:   "a",
			},
		}
		mockConfigV0 := &mocks.ConfigV0Interface{}
		mockConfigV0.On("IsExistUser").Return(true)
		mockConfigV0.On("User").Return(expectedUserConfig, nil)
		mockConfigV0.On("DeleteUser").Return(nil)
		mockConfigV0.On("IsExistAlertPlayers").Return(true)
		mockConfigV0.On("AlertPlayers").Return(expectedAlertPlayers, nil)
		mockConfigV0.On("DeleteAlertPlayers").Return(nil)
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadDataVersion").Return(uint(0), nil)
		mockStorage.On("IsExistUserConfig").Return(false)
		mockStorage.On("WriteUserConfig", expectedUserConfig).Return(nil)
		mockStorage.On("IsExistAlertPlayers").Return(false)
		mockStorage.On("WriteAlertPlayers", expectedAlertPlayers).Return(nil)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockConfigV0, mockStorage, nil)
		err := cm.Execute()

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
		mockStorage.On("ReadDataVersion").Return(uint(1), nil)

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
		mockStorage.On("ReadDataVersion").Return(uint(0), nil)
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
