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
		mockLocalFile := &mocks.LocalFileInterface{}
		mockLocalFile.On("IsExistUser").Return(true)
		mockLocalFile.On("User").Return(expectedUserConfig, nil)
		mockLocalFile.On("DeleteUser").Return(nil)
		mockLocalFile.On("IsExistAlertPlayers").Return(true)
		mockLocalFile.On("AlertPlayers").Return(expectedAlertPlayers, nil)
		mockLocalFile.On("DeleteAlertPlayers").Return(nil)
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadDataVersion").Return(uint(0), nil)
		mockStorage.On("IsExistUserConfig").Return(false)
		mockStorage.On("WriteUserConfig", expectedUserConfig).Return(nil)
		mockStorage.On("IsExistAlertPlayers").Return(false)
		mockStorage.On("WriteAlertPlayers", expectedAlertPlayers).Return(nil)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockLocalFile, mockStorage)
		err := cm.Execute()

		// アサーション
		require.NoError(t, err)
		mockLocalFile.AssertExpectations(t)
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
		cm := NewConfigMigrator(nil, mockStorage)
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
		mockLocalFile := &mocks.LocalFileInterface{}
		mockLocalFile.On("IsExistUser").Return(true)
		mockLocalFile.On("IsExistAlertPlayers").Return(true)
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadDataVersion").Return(uint(0), nil)
		mockStorage.On("IsExistUserConfig").Return(true)
		mockStorage.On("IsExistAlertPlayers").Return(true)
		mockStorage.On("WriteDataVersion", uint(1)).Return(nil)

		// テスト
		cm := NewConfigMigrator(mockLocalFile, mockStorage)
		err := cm.toV1()

		// アサーション
		require.NoError(t, err)
		mockLocalFile.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
		mockStorage.AssertNotCalled(t, "WriteUserConfig")
		mockStorage.AssertNotCalled(t, "WriteAlertPlayers")
	})
}
