package service

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/vo"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	DefaulfInstallPath = "install_path_test"
	GameExeName        = "WorldOfWarships.exe"
)

var (
	errWargaming          = apperr.New(apperr.WargamingAPIError, errors.New("INVALID_APPLICATION_ID"))
	errInvalidInstallPath = apperr.New(
		apperr.ValidateInvalidInstallPath,
		errors.New("stat invalid/path/WorldOfWarships.exe: no such file or directory"),
	)
	errInvalidAppID    = apperr.New(apperr.ValidateInvalidAppID, errWargaming)
	errInvalidFontSize = apperr.New(apperr.ValidateInvalidFontSize, nil)
)

//nolint:paralleltest
func TestConfig_UpdateUser_正常系(t *testing.T) {
	gameExePath := filepath.Join(DefaulfInstallPath, GameExeName)
	err := os.MkdirAll(DefaulfInstallPath, fs.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(DefaulfInstallPath)

	// テストデータ
	config := vo.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	// モックの設定
	mockConfigRepo := &mockConfigRepo{}
	mockConfigRepo.On("UpdateUser", config).Return(nil)
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid).Return()
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err = c.UpdateUser(config)

	// アサーション
	assert.NoError(t, err)
	mockWargamingRepo.AssertCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(t, "EncycInfo")
	mockConfigRepo.AssertCalled(t, "UpdateUser", config)
}

//nolint:paralleltest
func TestConfig_UpdateUser_異常系_不正なインストールパス(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaulfInstallPath)

	// テストデータ
	config := vo.UserConfig{
		InstallPath: "invalid/path",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	// モックの設定
	mockConfigRepo := &mockConfigRepo{}

	// テスト実行
	c := NewConfig(mockConfigRepo, nil)
	err = c.UpdateUser(config)

	// アサーション
	assert.EqualError(t, err, errInvalidInstallPath.Error())
	mockConfigRepo.AssertNotCalled(t, "UpdateUser")
}

//nolint:paralleltest
func TestConfig_UpdateUser_異常系_不正なAppID(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaulfInstallPath)

	// テストデータ
	config := vo.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "invalidappid",
		FontSize:    "medium",
	}

	// モックの設定
	mockConfigRepo := &mockConfigRepo{}
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid)
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, errWargaming)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err = c.UpdateUser(config)

	// アサーション
	assert.EqualError(t, err, errInvalidAppID.Error())
	mockWargamingRepo.AssertCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(t, "EncycInfo")
	mockConfigRepo.AssertNotCalled(t, "UpdateUser")
}

//nolint:paralleltest
func TestConfig_UpdateUser_異常系_不正なフォントサイズ(t *testing.T) {
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(DefaulfInstallPath)

	// テストデータ
	config := vo.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "invalid",
	}

	// モックの設定
	mockConfigRepo := &mockConfigRepo{}
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", config.Appid)
	mockWargamingRepo.On("EncycInfo").Return(vo.WGEncycInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err = c.UpdateUser(config)

	// アサーション
	assert.EqualError(t, err, errInvalidFontSize.Error())
	mockWargamingRepo.AssertCalled(t, "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(t, "EncycInfo")
	mockConfigRepo.AssertNotCalled(t, "UpdateUser")
}

func createGameClientPath() error {
	if err := os.MkdirAll(DefaulfInstallPath, fs.ModePerm); err != nil {
		return err
	}

	gameExePath := filepath.Join(DefaulfInstallPath, GameExeName)

	return os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
}
