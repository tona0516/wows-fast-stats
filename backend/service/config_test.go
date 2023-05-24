package service

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	DefaulfInstallPath = "install_path_test"
	GameExeName        = "WorldOfWarships.exe"
)

type ConfigServiceSuite struct {
	suite.Suite
}

func (suite *ConfigServiceSuite) SetupSuite() {
	gameExePath := filepath.Join(DefaulfInstallPath, GameExeName)
	_ = os.MkdirAll(DefaulfInstallPath, fs.ModePerm)
	_ = os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
}

func (suite *ConfigServiceSuite) TearDownSuite() {
	_ = os.RemoveAll(DefaulfInstallPath)
}

func (suite *ConfigServiceSuite) TestConfig_UpdateUser() {
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
	mockWargamingRepo.On("EncyclopediaInfo").Return(vo.WGEncyclopediaInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err := c.UpdateUser(config)

	// 結果の検証
	assert.NoError(suite.T(), err)
	mockWargamingRepo.AssertCalled(suite.T(), "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(suite.T(), "EncyclopediaInfo")
	mockConfigRepo.AssertCalled(suite.T(), "UpdateUser", config)
}

func (suite *ConfigServiceSuite) TestConfig_UpdateUser_InvalidInstallPath() {
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
	err := c.UpdateUser(config)

	// 結果の検証
	expectedErr := errors.WithStack(apperr.SrvCfg.InvalidInstallPath.WithRaw(apperr.ErrInvalidInstallPath))
	assert.EqualError(suite.T(), err, expectedErr.Error())
	mockConfigRepo.AssertNotCalled(suite.T(), "UpdateUser")
}

func (suite *ConfigServiceSuite) TestConfig_UpdateUser_InvalidAppID() {
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
	expectedErr := errors.WithStack(apperr.SrvCfg.InvalidAppID.WithRaw(apperr.ErrInvalidAppID))
	mockWargamingRepo.On("EncyclopediaInfo").Return(vo.WGEncyclopediaInfo{}, expectedErr)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err := c.UpdateUser(config)

	// 結果の検証
	assert.EqualError(suite.T(), err, expectedErr.Error())
	mockWargamingRepo.AssertCalled(suite.T(), "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(suite.T(), "EncyclopediaInfo")
	mockConfigRepo.AssertNotCalled(suite.T(), "UpdateUser")
}

func (suite *ConfigServiceSuite) TestConfig_UpdateUser_InvalidFontSize() {
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
	mockWargamingRepo.On("EncyclopediaInfo").Return(vo.WGEncyclopediaInfo{}, nil)

	// テスト実行
	c := NewConfig(mockConfigRepo, mockWargamingRepo)
	err := c.UpdateUser(config)

	// 結果の検証
	expectedErr := errors.WithStack(apperr.SrvCfg.InvalidFontSize.WithRaw(apperr.ErrInvalidFontSize))
	assert.EqualError(suite.T(), err, expectedErr.Error())
	mockWargamingRepo.AssertCalled(suite.T(), "SetAppID", config.Appid)
	mockWargamingRepo.AssertCalled(suite.T(), "EncyclopediaInfo")
	mockConfigRepo.AssertNotCalled(suite.T(), "UpdateUser")
}

//nolint:paralleltest
func TestConfigServiceSuite(t *testing.T) {
	suite.Run(t, new(ConfigServiceSuite))
}
