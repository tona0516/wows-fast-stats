//nolint:paralleltest
package usecase

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	DefaultInstallPath = "install_path_test"
	DefaultAppID       = "abc123"
)

var errWargaming = failure.New(apperr.WGAPIError)

func TestConfig_UpdateRequired_正常系(t *testing.T) {
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()

	// モックの設定
	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("UpdateUser", config).Return(nil)
	mockWargaming := &mocks.WargamingInterface{}
	mockWargaming.On("Test", mock.Anything).Return(true, nil)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(domain.DefaultUserConfig, nil)
	mockStorage.On("WriteUserConfig", mock.Anything).Return(nil)

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming, mockStorage)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.RequiredConfigError{Valid: true}, actual)
	require.NoError(t, err)
	mockWargaming.AssertCalled(t, "Test", config.Appid)
	mockStorage.AssertCalled(t, "ReadUserConfig")
	mockStorage.AssertCalled(t, "WriteUserConfig", config)
}

func TestConfig_UpdateRequired_異常系_不正なインストールパス(t *testing.T) {
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := domain.DefaultUserConfig
	config.InstallPath = "invalid/path" // Note: 不正なパス
	config.Appid = "abc123"

	// モックの設定
	mockLocalFile := &mocks.LocalFileInterface{}
	mockWargaming := &mocks.WargamingInterface{}
	mockWargaming.On("Test", mock.Anything).Return(true, nil)
	mockStorage := &mocks.StorageInterface{}

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming, mockStorage)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.RequiredConfigError{InstallPath: apperr.InvalidInstallPath.ErrorCode()}, actual)
	require.NoError(t, err)
	mockWargaming.AssertCalled(t, "Test", config.Appid)
	mockStorage.AssertNotCalled(t, "WriteUserConfig", mock.Anything)
}

func TestConfig_UpdateRequired_異常系_不正なAppID(t *testing.T) {
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()

	// モックの設定
	mockLocalFile := &mocks.LocalFileInterface{}
	mockWargaming := &mocks.WargamingInterface{}
	mockWargaming.On("Test", mock.Anything).Return(false, errWargaming) // Note: WG APIでエラー
	mockStorage := &mocks.StorageInterface{}

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming, mockStorage)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.RequiredConfigError{AppID: apperr.InvalidAppID.ErrorCode()}, actual)
	require.NoError(t, err)
	mockWargaming.AssertCalled(t, "Test", config.Appid)
	mockStorage.AssertNotCalled(t, "WriteUserConfig", mock.Anything)
}

func TestConfig_UpdateOptional_正常系(t *testing.T) {
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()
	config.FontSize = "small"

	// モックの設定
	mockLocalFile := &mocks.LocalFileInterface{}
	// Note: requiredな値を与えてもこれらの値はUpdateUserでは含まれない
	actualWritten := domain.DefaultUserConfig
	actualWritten.FontSize = "small"
	mockWargaming := &mocks.WargamingInterface{}
	mockWargaming.On("Test", mock.Anything).Return(true, nil)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(domain.DefaultUserConfig, nil)
	mockStorage.On("WriteUserConfig", actualWritten).Return(nil)

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming, mockStorage)
	err = c.UpdateOptional(config)

	// アサーション
	require.NoError(t, err)

	mockWargaming.AssertNotCalled(t, "Test", config.Appid)
	mockStorage.AssertCalled(t, "ReadUserConfig")
	mockStorage.AssertCalled(t, "WriteUserConfig", actualWritten)
}

func TestConfig_AlertPlayers_正常系(t *testing.T) {
	expected := []domain.AlertPlayer{
		{AccountID: 1, Name: "Player1"},
		{AccountID: 2, Name: "Player2"},
	}
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadAlertPlayers").Return(expected, nil)

	config := NewConfig(nil, nil, mockStorage)
	actual, err := config.AlertPlayers()

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	mockStorage.AssertExpectations(t)
}

func TestConfig_UpdateAlertPlayer_正常系_追加(t *testing.T) {
	existingPlayers := []domain.AlertPlayer{
		{AccountID: 1, Name: "Player1"},
		{AccountID: 2, Name: "Player2"},
	}
	newPlayer := domain.AlertPlayer{AccountID: 3, Name: "Player3"}
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadAlertPlayers").Return(existingPlayers, nil)
	mockStorage.On("WriteAlertPlayers", append(existingPlayers, newPlayer)).Return(nil)

	config := NewConfig(nil, nil, mockStorage)
	err := config.UpdateAlertPlayer(newPlayer)

	require.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestConfig_UpdateAlertPlayer_正常系_更新(t *testing.T) {
	existingPlayers := []domain.AlertPlayer{
		{AccountID: 1, Name: "Player1"},
		{AccountID: 2, Name: "Player2"},
	}
	playerToUpdate := domain.AlertPlayer{AccountID: 1, Name: "UpdatedPlayer"}

	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadAlertPlayers").Return(existingPlayers, nil)
	mockStorage.On("WriteAlertPlayers", []domain.AlertPlayer{
		{AccountID: 1, Name: "UpdatedPlayer"},
		{AccountID: 2, Name: "Player2"},
	}).Return(nil)

	config := NewConfig(nil, nil, mockStorage)
	err := config.UpdateAlertPlayer(playerToUpdate)

	require.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestConfig_RemoveAlertPlayer_正常系(t *testing.T) {
	accountIDToRemove := 1
	existingPlayers := []domain.AlertPlayer{
		{AccountID: 1, Name: "Player1"},
		{AccountID: 2, Name: "Player2"},
	}

	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadAlertPlayers").Return(existingPlayers, nil)
	mockStorage.On("WriteAlertPlayers", mock.Anything).Return(nil)

	config := NewConfig(nil, nil, mockStorage)
	err := config.RemoveAlertPlayer(accountIDToRemove)

	// Assert
	require.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func createInputConfig() domain.UserConfig {
	config := domain.DefaultUserConfig
	config.InstallPath = "install_path_test"
	config.Appid = "abc123"

	return config
}

func createGameClientPath() error {
	if err := os.MkdirAll(DefaultInstallPath, fs.ModePerm); err != nil {
		return err
	}

	gameExePath := filepath.Join(DefaultInstallPath, GameExeName)

	return os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
}
