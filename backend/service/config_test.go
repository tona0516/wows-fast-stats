package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/infra"
	"wfs/backend/vo"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultInstallPath = "install_path_test"
	DefaultAppID       = "abc123"
)

var errWargaming = apperr.New(apperr.WargamingAPIError, errors.New("INVALID_APPLICATION_ID"))

//nolint:paralleltest
func TestConfig_UpdateRequired_正常系(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()

	// モックの設定
	mockConfigRepo := &mockConfigRepo{}
	mockConfigRepo.On("UpdateUser", config).Return(nil)
	mockConfigRepo.On("User").Return(infra.DefaultUserConfig, nil)
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid).Return()
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.ValidatedResult{}, actual)
	assert.NoError(t, err)
	mockWargamingRepo.AssertCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(t, "EncycInfo")
	mockConfigRepo.AssertCalled(t, "User")
	mockConfigRepo.AssertCalled(t, "UpdateUser", config)
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
	mockConfigRepo := &mockConfigRepo{}
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid).Return()
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.ValidatedResult{InstallPath: apperr.ErrInvalidInstallPath.Error()}, actual)
	assert.NoError(t, err)
	mockWargamingRepo.AssertCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(t, "EncycInfo")
	mockConfigRepo.AssertNotCalled(t, "UpdateUser", config)
}

//nolint:paralleltest
func TestConfig_UpdateRequired_異常系_不正なAppID(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaultInstallPath)

	// テストデータ
	config := createInputConfig()

	// モックの設定
	mockConfigRepo := &mockConfigRepo{}
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid).Return()
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, errWargaming) // Note: WG APIでエラー

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

	// アサーション
	assert.Equal(t, vo.ValidatedResult{AppID: fmt.Sprintf("%s(%s)", apperr.ErrInvalidAppID, errWargaming)}, actual)
	assert.NoError(t, err)
	mockWargamingRepo.AssertCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(t, "EncycInfo")
	mockConfigRepo.AssertNotCalled(t, "UpdateUser", config)
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
	mockConfigRepo := &mockConfigRepo{}
	// Note: requiredな値を与えてもこれらの値はUpdateUserでは含まれない
	actualWritten := infra.DefaultUserConfig
	actualWritten.FontSize = "small"
	mockConfigRepo.On("UpdateUser", actualWritten).Return(nil)
	mockConfigRepo.On("User").Return(infra.DefaultUserConfig, nil)
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid).Return()
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err = c.UpdateOptional(config)

	// アサーション
	assert.NoError(t, err)
	mockConfigRepo.AssertCalled(t, "User")
	mockConfigRepo.AssertCalled(t, "UpdateUser", actualWritten)
	mockWargamingRepo.AssertNotCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertNotCalled(t, "EncycInfo")
}

func createInputConfig() vo.UserConfig {
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
