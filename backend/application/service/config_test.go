package service

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/infra"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	DefaultInstallPath = "install_path_test"
	DefaultAppID       = "abc123"
)

var errWargaming = failure.New(apperr.WGAPIError)

//nolint:paralleltest
func TestConfig_UpdateRequired_正常系(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()

	// モックの設定
	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("UpdateUser", config).Return(nil)
	mockLocalFile.On("User").Return(infra.DefaultUserConfig, nil)
	mockWargaming := &mockWargaming{}
	mockWargaming.On("Test", mock.Anything).Return(true, nil)

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.RequiredConfigError{Valid: true}, actual)
	assert.NoError(t, err)
	mockWargaming.AssertCalled(t, "Test", config.Appid)
	mockLocalFile.AssertCalled(t, "User")
	mockLocalFile.AssertCalled(t, "UpdateUser", config)
}

//nolint:paralleltest
func TestConfig_UpdateRequired_異常系_不正なインストールパス(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := infra.DefaultUserConfig
	config.InstallPath = "invalid/path" // Note: 不正なパス
	config.Appid = "abc123"

	// モックの設定
	mockLocalFile := &mockLocalFile{}
	mockWargaming := &mockWargaming{}
	mockWargaming.On("Test", mock.Anything).Return(true, nil)

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.RequiredConfigError{InstallPath: apperr.InvalidInstallPath.ErrorCode()}, actual)
	assert.NoError(t, err)
	mockWargaming.AssertCalled(t, "Test", config.Appid)
	mockLocalFile.AssertNotCalled(t, "UpdateUser", config)
}

//nolint:paralleltest
func TestConfig_UpdateRequired_異常系_不正なAppID(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()

	// モックの設定
	mockLocalFile := &mockLocalFile{}
	mockWargaming := &mockWargaming{}
	mockWargaming.On("Test", mock.Anything).Return(false, errWargaming) // Note: WG APIでエラー

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.RequiredConfigError{AppID: apperr.InvalidAppID.ErrorCode()}, actual)
	assert.NoError(t, err)
	mockWargaming.AssertCalled(t, "Test", config.Appid)
	mockLocalFile.AssertNotCalled(t, "UpdateUser", config)
}

//nolint:paralleltest
func TestConfig_UpdateOptional_正常系(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()
	config.FontSize = "small"

	// モックの設定
	mockLocalFile := &mockLocalFile{}
	// Note: requiredな値を与えてもこれらの値はUpdateUserでは含まれない
	actualWritten := infra.DefaultUserConfig
	actualWritten.FontSize = "small"
	mockLocalFile.On("UpdateUser", actualWritten).Return(nil)
	mockLocalFile.On("User").Return(infra.DefaultUserConfig, nil)
	mockWargaming := &mockWargaming{}
	mockWargaming.On("Test", mock.Anything).Return(true, nil)

	// テスト実行
	c := NewConfig(mockLocalFile, mockWargaming)
	err = c.UpdateOptional(config)

	// アサーション
	assert.NoError(t, err)
	mockLocalFile.AssertCalled(t, "User")
	mockLocalFile.AssertCalled(t, "UpdateUser", actualWritten)
	mockWargaming.AssertNotCalled(t, "Test", config.Appid)
}

func createInputConfig() domain.UserConfig {
	config := infra.DefaultUserConfig
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
