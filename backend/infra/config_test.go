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
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDirName)

	expected := vo.UserConfig{
		FontSize: "large",
		Displays: vo.Displays{
			Basic: vo.Basic{
				IsInAvg:    false,
				PlayerName: true,
				ShipInfo:   false,
			},
		},
	}

	config := NewConfig()

	// 書き込み：正常系
	err := config.UpdateUser(expected)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := config.User()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：異常系 存在しない場合
	err = os.Remove(filepath.Join(ConfigDirName, ConfigUserName))
	assert.NoError(t, err)
	_, err = config.User()
	assert.EqualError(t, err, "I400 Read open config/user.json: no such file or directory")
}

//nolint:paralleltest
func TestConfig_App(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDirName)

	// テスト用のアプリケーション設定
	expected := vo.AppConfig{
		Window: vo.WindowConfig{
			Width:  100,
			Height: 100,
		},
	}

	config := NewConfig()

	// 書き込み：正常系
	err := config.UpdateApp(expected)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := config.App()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：異常系 存在しない場合
	err = os.Remove(filepath.Join(ConfigDirName, ConfigAppName))
	assert.NoError(t, err)
	_, err = config.App()
	assert.EqualError(t, err, "I400 Read open config/app.json: no such file or directory")
}
