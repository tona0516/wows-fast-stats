//nolint:paralleltest
package service

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
