package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"wfs/backend/vo"
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

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(ConfigDirName, ConfigUserName))
	assert.NoError(t, err)

	actual, err = config.User()
	assert.NoError(t, err)
	assert.Equal(t, defaultUserConfig, actual)
}

//nolint:paralleltest
func TestConfig_App(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDirName)

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

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(ConfigDirName, ConfigAppName))
	assert.NoError(t, err)

	actual, err = config.App()
	assert.NoError(t, err)
	assert.Equal(t, vo.AppConfig{}, actual)
}

//nolint:paralleltest
func TestConfig_AlertPlayers(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDirName)

	expected1 := vo.AlertPlayer{
		AccountID: 100,
		Name:      "test",
		Pattern:   "bi-check-circle-fill",
		Message:   "hello",
	}
	expected2 := vo.AlertPlayer{
		AccountID: 200,
		Name:      "hoge",
		Pattern:   "bi-check-circle-fill",
		Message:   "memo",
	}

	config := NewConfig()

	// 書き込み：正常系
	err := config.UpdateAlertPlayer(expected1)
	assert.NoError(t, err)
	err = config.UpdateAlertPlayer(expected2)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := config.AlertPlayers()
	assert.NoError(t, err)
	assert.Equal(t, []vo.AlertPlayer{expected1, expected2}, actual)

	// 削除：正常系
	err = config.RemoveAlertPlayer(100)
	assert.NoError(t, err)

	// 取得
	actual, err = config.AlertPlayers()
	assert.NoError(t, err)
	assert.Equal(t, []vo.AlertPlayer{expected2}, actual)

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(ConfigDirName, ConfigAlertPlayerName))
	assert.NoError(t, err)

	actual, err = config.AlertPlayers()
	assert.NoError(t, err)
	assert.Equal(t, []vo.AlertPlayer{}, actual)
}
