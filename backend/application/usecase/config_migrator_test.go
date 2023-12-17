package usecase

import (
	"testing"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfigMigrator_Migrate_正常系(t *testing.T) {
	t.Parallel()

	// mockling
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
	mockStorage.On("WriteUserConfig", mock.Anything).Return(nil)
	mockStorage.On("IsExistAlertPlayers").Return(false)
	mockStorage.On("WriteAlertPlayers", mock.Anything).Return(nil)
	mockStorage.On("WriteDataVersion", mock.Anything).Return(nil)

	// test
	cm := NewConfigMigrator(mockLocalFile, mockStorage)
	err := cm.Execute()

	// assertion
	require.NoError(t, err)
	mockLocalFile.AssertCalled(t, "User")
	mockLocalFile.AssertCalled(t, "AlertPlayers")
	mockStorage.AssertCalled(t, "ReadDataVersion")
	mockStorage.AssertCalled(t, "WriteUserConfig", expectedUserConfig)
	mockStorage.AssertCalled(t, "WriteAlertPlayers", expectedAlertPlayers)
	mockStorage.AssertCalled(t, "WriteDataVersion", uint(1))
}

func TestConfigMigrator_Migrate_正常系_マイグレ不要_バージョン1以上(t *testing.T) {
	t.Parallel()

	// mocking
	mockLocalFile := &mocks.LocalFileInterface{}
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadDataVersion").Return(uint(1), nil)

	// test
	cm := NewConfigMigrator(mockLocalFile, mockStorage)
	err := cm.Execute()

	// assertion
	require.NoError(t, err)
	mockStorage.AssertNotCalled(t, "WriteUserConfig")
	mockStorage.AssertNotCalled(t, "WriteAlertPlayers")
	mockStorage.AssertNotCalled(t, "WriteDataVersion")
}

func TestConfigMigrator_Migrate_正常系_マイグレ不要_すでにストレージに存在(t *testing.T) {
	t.Parallel()

	// mocking
	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("IsExistUser").Return(true)
	mockLocalFile.On("IsExistAlertPlayers").Return(true)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadDataVersion").Return(uint(0), nil)
	mockStorage.On("IsExistUserConfig").Return(true)
	mockStorage.On("IsExistAlertPlayers").Return(true)
	mockStorage.On("WriteDataVersion", mock.Anything).Return(nil)

	// test
	cm := NewConfigMigrator(mockLocalFile, mockStorage)
	err := cm.Execute()

	// assertion
	require.NoError(t, err)
	mockStorage.AssertNotCalled(t, "WriteUserConfig")
	mockStorage.AssertNotCalled(t, "WriteAlertPlayers")
	mockStorage.AssertCalled(t, "WriteDataVersion", uint(1))
}
