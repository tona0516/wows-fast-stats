//nolint:paralleltest
package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigV0_User(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		defer os.RemoveAll(ConfigDir)

		expected := data.UserConfig{
			FontSize: "large",
			Displays: data.Displays{
				Ship:    data.Ship{PR: true},
				Overall: data.Overall{PR: true},
			},
		}

		instance := NewConfigV0()
		err := writeJSON(instance.userConfigPath, expected)
		require.NoError(t, err)

		// 取得
		actual, err := instance.User()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)

		// 削除
		err = os.Remove(instance.userConfigPath)
		require.NoError(t, err)

		// 取得 存在しない場合 デフォルト値を返却する
		actual, err = instance.User()
		require.NoError(t, err)
		assert.Equal(t, data.DefaultUserConfig(), actual)
	})

	t.Run("正常系_ファイルに新規パラメータが存在しない", func(t *testing.T) {
		// テストで生成したディレクトリを削除
		defer os.RemoveAll(ConfigDir)

		// 必須項目のみconfig.jsonに書き込む
		installPath := "dir/"
		appid := "abc"
		saved := fmt.Sprintf(`{"install_path": "%s","appid": "%s"}`, installPath, appid)
		err := os.Mkdir(ConfigDir, os.ModePerm)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(ConfigDir, UserConfigFile), []byte(saved), 0o600)
		require.NoError(t, err)

		instance := NewConfigV0()
		actual, err := instance.User()
		require.NoError(t, err)

		// 存在するパラメータはその値、存在しないパラメータはデフォルト値が格納されていること
		expected := data.DefaultUserConfig()
		expected.InstallPath = installPath
		expected.Appid = appid
		require.Equal(t, expected, actual)
	})
}

func TestConfigV0_AlertPlayers(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDir)

	expected := []data.AlertPlayer{
		{
			AccountID: 100,
			Name:      "test",
			Pattern:   "bi-check-circle-fill",
			Message:   "hello",
		},
		{
			AccountID: 200,
			Name:      "hoge",
			Pattern:   "bi-check-circle-fill",
			Message:   "memo",
		},
	}

	instance := NewConfigV0()
	err := writeJSON(instance.alertPlayerPath, expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err := instance.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(ConfigDir, AlertPlayerFile))
	require.NoError(t, err)

	actual, err = instance.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, []data.AlertPlayer{}, actual)
}
