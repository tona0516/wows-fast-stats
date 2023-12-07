package service

import (
	"testing"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfigMigrator_Migrate_正常系(t *testing.T) {
	t.Parallel()

	// モックの設定
	expectedUserConfig := domain.UserConfig{
		FontSize: "large",
		Displays: domain.Displays{
			Ship:    domain.Ship{PR: true},
			Overall: domain.Overall{PR: false},
		},
	}
	expectedAlertPlayers := []domain.AlertPlayer{
		{
			AccountID: 100,
			Name:      "test",
			Pattern:   "bi-check-circle-fill",
			Message:   "hello",
		},
		{
			AccountID: 200,
			Name:      "hoge",
			Pattern:   "bi-check-circle-fill",
			Message:   "memo",
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
	mockStorage.On("IsExistUserConfig", mock.Anything).Return(false)
	mockStorage.On("WriteUserConfig", mock.Anything).Return(nil)
	mockStorage.On("IsExistAlertPlayers", mock.Anything).Return(false)
	mockStorage.On("WriteAlertPlayers", mock.Anything).Return(nil)
	mockStorage.On("WriteDataVersion", mock.Anything).Return(nil)

	// テスト実行
	cm := NewConfigMigrator(mockLocalFile, mockStorage)
	err := cm.Execute()

	// アサーション
	require.NoError(t, err)
	mockLocalFile.AssertCalled(t, "User")
	mockLocalFile.AssertCalled(t, "AlertPlayers")
	mockStorage.AssertCalled(t, "ReadDataVersion")
	mockStorage.AssertCalled(t, "WriteUserConfig", expectedUserConfig)
	mockStorage.AssertCalled(t, "WriteAlertPlayers", expectedAlertPlayers)
	mockStorage.AssertCalled(t, "WriteDataVersion", uint(1))
}
