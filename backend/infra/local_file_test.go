//nolint:paralleltest
package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testInstallPath = "testdata"

func TestLocalFile_User(t *testing.T) {
	defer os.RemoveAll(ConfigDir)

	expected := domain.UserConfig{
		FontSize: "large",
		Displays: domain.Displays{
			Ship: domain.Ship{
				PR: true,
			},
			Overall: domain.Overall{
				PR: false,
			},
		},
	}

	localFile := NewLocalFile()
	err := writeJSON(localFile.userConfigPath, expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err := localFile.User()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(localFile.userConfigPath)
	require.NoError(t, err)

	actual, err = localFile.User()
	require.NoError(t, err)
	assert.Equal(t, domain.DefaultUserConfig, actual)
}

func TestLocalFile_User_異常系_ファイルに新規パラメータが存在しない(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDir)

	// 必須項目のみconfig.jsonに書き込む
	installPath := "dir/"
	appid := "abc"
	saved := fmt.Sprintf(`{"install_path": "%s","appid": "%s"}`, installPath, appid)
	err := os.Mkdir(ConfigDir, os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(ConfigDir, UserConfigFile), []byte(saved), os.ModePerm)
	require.NoError(t, err)

	localFile := NewLocalFile()
	actual, err := localFile.User()
	require.NoError(t, err)

	// 存在するパラメータはその値、存在しないパラメータはデフォルト値が格納されていること
	expected := domain.DefaultUserConfig
	expected.InstallPath = installPath
	expected.Appid = appid
	require.Equal(t, expected, actual)
}

func TestLocalFile_AlertPlayers(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(ConfigDir)

	expected := []domain.AlertPlayer{
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

	localFile := NewLocalFile()
	err := writeJSON(localFile.alertPlayerPath, expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err := localFile.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(ConfigDir, AlertPlayerFile))
	require.NoError(t, err)

	actual, err = localFile.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, []domain.AlertPlayer{}, actual)
}

func TestLocalFile_SaveScreenshot_正常系(t *testing.T) {
	t.Parallel()

	// テストデータの作成
	rawData := "Hello, world!"
	base64Data := "SGVsbG8sIHdvcmxkIQ=="
	path := "screenshot_test/screenshot.png"

	// テストで生成したディレクトリを削除
	defer os.RemoveAll(filepath.Dir(path))

	// テスト
	localFile := LocalFile{}
	err := localFile.SaveScreenshot(path, base64Data)

	// アサーション
	require.NoError(t, err)
	assert.FileExists(t, path)
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, []byte(rawData))
}

func TestLocalFile_GetTempArenaInfo_正常系(t *testing.T) {
	localFile := NewLocalFile()
	expected := domain.TempArenaInfo{
		Vehicles: []domain.Vehicle{
			{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
			{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
			{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
		},
		DateTime:   "22.05.2023 12:34:56",
		MapID:      10,
		MatchGroup: "pvp",
		PlayerName: "player_1",
	}

	paths := []string{
		filepath.Join(testInstallPath, replaysDir, tempArenaInfoFile),
		filepath.Join(testInstallPath, replaysDir, "12.4.0", tempArenaInfoFile),
	}

	for _, path := range paths {
		func(path string) {
			defer os.RemoveAll(testInstallPath)

			err := writeJSON(path, expected)
			require.NoError(t, err)

			actual, err := localFile.TempArenaInfo(testInstallPath)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		}(path)
	}
}

func TestLocalFile_GetTempArenaInfo_正常系_該当ファイルが複数存在する場合_最新を返す(t *testing.T) {
	localFile := NewLocalFile()

	older := domain.TempArenaInfo{
		Vehicles: []domain.Vehicle{
			{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
			{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
			{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
		},
		DateTime:   "22.05.2022 12:34:56", // older than expected
		MapID:      10,
		MatchGroup: "pvp",
		PlayerName: "player_1",
	}

	expected := domain.TempArenaInfo{
		Vehicles: []domain.Vehicle{
			{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
			{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
			{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
		},
		DateTime:   "22.05.2023 12:34:56",
		MapID:      10,
		MatchGroup: "pvp",
		PlayerName: "player_1",
	}

	installPath := "testdata"
	defer os.RemoveAll(installPath)

	err := writeJSON(filepath.Join(installPath, replaysDir, tempArenaInfoFile), older)
	require.NoError(t, err)
	err = writeJSON(filepath.Join(installPath, replaysDir, "12.4.0", tempArenaInfoFile), expected)
	require.NoError(t, err)

	actual, err := localFile.TempArenaInfo(installPath)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestLocalFile_GetTempArenaInfo_異常系_該当ファイルなし(t *testing.T) {
	localFile := NewLocalFile()

	installPath := "testdata"
	paths := []string{
		filepath.Join(installPath, replaysDir, "hoge.wowsreplay"),
		filepath.Join(installPath, replaysDir, "12.4.0", "hoge.wowsreplay"),
	}

	for _, path := range paths {
		func(path string) {
			defer os.RemoveAll(testInstallPath)

			err := writeJSON(path, domain.TempArenaInfo{})
			require.NoError(t, err)

			_, err = localFile.TempArenaInfo(installPath)
			require.Error(t, err)
		}(path)
	}
}
