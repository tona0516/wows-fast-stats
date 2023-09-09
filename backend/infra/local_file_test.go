package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"wfs/backend/application/vo"
	"wfs/backend/domain"
)

const testInstallPath = "testdata"

//nolint:paralleltest
func TestLocalFile_User(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(configDir)

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

	// 書き込み：正常系
	err := localFile.UpdateUser(expected)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := localFile.User()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(configDir, userConfigFile))
	assert.NoError(t, err)

	actual, err = localFile.User()
	assert.NoError(t, err)
	assert.Equal(t, DefaultUserConfig, actual)
}

//nolint:paralleltest
func TestLocalFile_User_異常系_ファイルに新規パラメータが存在しない(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(configDir)

	// 必須項目のみconfig.jsonに書き込む
	installPath := "dir/"
	appid := "abc"
	saved := fmt.Sprintf(`{"install_path": "%s","appid": "%s"}`, installPath, appid)
	err := os.Mkdir(configDir, os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(configDir, userConfigFile), []byte(saved), os.ModePerm)
	assert.NoError(t, err)

	localFile := NewLocalFile()
	actual, err := localFile.User()
	assert.NoError(t, err)

	// 存在するパラメータはその値、存在しないパラメータはデフォルト値が格納されていること
	expected := DefaultUserConfig
	expected.InstallPath = installPath
	expected.Appid = appid
	assert.Equal(t, expected, actual)
}

//nolint:paralleltest
func TestLocalFile_App(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(configDir)

	expected := vo.AppConfig{
		Window: vo.WindowConfig{
			Width:  100,
			Height: 100,
		},
	}

	localFile := NewLocalFile()

	// 書き込み：正常系
	err := localFile.UpdateApp(expected)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := localFile.App()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(configDir, appConfigFile))
	assert.NoError(t, err)

	actual, err = localFile.App()
	assert.NoError(t, err)
	assert.Equal(t, vo.AppConfig{}, actual)
}

//nolint:paralleltest
func TestLocalFile_AlertPlayers(t *testing.T) {
	// テストで生成したディレクトリを削除
	defer os.RemoveAll(configDir)

	expected1 := domain.AlertPlayer{
		AccountID: 100,
		Name:      "test",
		Pattern:   "bi-check-circle-fill",
		Message:   "hello",
	}
	expected2 := domain.AlertPlayer{
		AccountID: 200,
		Name:      "hoge",
		Pattern:   "bi-check-circle-fill",
		Message:   "memo",
	}

	localFile := NewLocalFile()

	// 書き込み：正常系
	err := localFile.UpdateAlertPlayer(expected1)
	assert.NoError(t, err)
	err = localFile.UpdateAlertPlayer(expected2)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := localFile.AlertPlayers()
	assert.NoError(t, err)
	assert.Equal(t, []domain.AlertPlayer{expected1, expected2}, actual)

	// 削除：正常系
	err = localFile.RemoveAlertPlayer(100)
	assert.NoError(t, err)

	// 取得
	actual, err = localFile.AlertPlayers()
	assert.NoError(t, err)
	assert.Equal(t, []domain.AlertPlayer{expected2}, actual)

	// 取得：異常系 存在しない場合 デフォルト値を返却する
	err = os.Remove(filepath.Join(configDir, alertPlayerFile))
	assert.NoError(t, err)

	actual, err = localFile.AlertPlayers()
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.FileExists(t, path)
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, content, []byte(rawData))
}

//nolint:paralleltest
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
			assert.NoError(t, err)

			actual, err := localFile.TempArenaInfo(testInstallPath)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}(path)
	}
}

//nolint:paralleltest
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
	assert.NoError(t, err)
	err = writeJSON(filepath.Join(installPath, replaysDir, "12.4.0", tempArenaInfoFile), expected)
	assert.NoError(t, err)

	actual, err := localFile.TempArenaInfo(installPath)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

//nolint:paralleltest
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
			assert.NoError(t, err)

			_, err = localFile.TempArenaInfo(installPath)
			assert.Error(t, err)
		}(path)
	}
}
