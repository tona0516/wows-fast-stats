package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"changeme/backend/vo"
)

//nolint:paralleltest
func TestConfig_User(t *testing.T) {
	// テスト用のユーザー設定
	userConfig := vo.UserConfig{
		FontSize: "large",
		Displays: vo.Displays{
			Basic: vo.Basic{
				IsInAvg:    false,
				PlayerName: true,
				ShipInfo:   false,
			},
		},
	}

	// テスト用のディレクトリと設定ファイルの作成
	err := os.Mkdir(ConfigDirName, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(ConfigDirName)

	// 書き込み：正常系
	config := &Config{}
	err = config.UpdateUser(userConfig)
	assert.NoError(t, err)

	// 取得：正常系
	loadedConfig, err := config.User()
	assert.NoError(t, err)
	assert.Equal(t, userConfig, loadedConfig)

	// 取得：異常系 存在しない場合
	err = os.Remove(filepath.Join(ConfigDirName, ConfigUserName))
	assert.NoError(t, err)
	_, err = config.User()
	assert.EqualError(t, err, "I400 Read open config/user.json: no such file or directory")
}

//nolint:paralleltest
func TestConfig_App(t *testing.T) {
	// テスト用のアプリケーション設定
	appConfig := vo.AppConfig{
		Window: vo.WindowConfig{
			Width:  100,
			Height: 100,
		},
	}

	// テスト用のディレクトリと設定ファイルの作成
	err := os.Mkdir(ConfigDirName, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(ConfigDirName)

	// 書き込み：正常系
	config := &Config{}
	err = config.UpdateApp(appConfig)
	assert.NoError(t, err)

	// 取得：正常系
	loadedConfig, err := config.App()
	assert.NoError(t, err)
	assert.Equal(t, appConfig, loadedConfig)

	// 取得：異常系 存在しない場合
	err = os.Remove(filepath.Join(ConfigDirName, ConfigAppName))
	assert.NoError(t, err)
	_, err = config.App()
	assert.EqualError(t, err, "I400 Read open config/app.json: no such file or directory")
}
