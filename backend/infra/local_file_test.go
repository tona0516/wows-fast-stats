//nolint:paralleltest
package infra

import (
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testInstallPath = "testdata"

func TestLocalFile_SaveScreenshot(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// テストデータの作成
		rawData := "Hello, world!"
		base64Data := "SGVsbG8sIHdvcmxkIQ=="
		path := "screenshot_test/screenshot.png"

		// テストで生成したディレクトリを削除
		defer os.RemoveAll(filepath.Dir(path))

		// テスト
		instance := LocalFile{}
		err := instance.SaveScreenshot(path, base64Data)

		// アサーション
		assert.NoError(t, err)
		assert.FileExists(t, path)
		content, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, content, []byte(rawData))
	})
}

func TestLocalFile_ReadTempArenaInfo(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		expected := model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
				{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
				{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
			},
			DateTime:   "22.05.2023 12:34:56",
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		paths := [][]string{
			{testInstallPath, replaysDir, tempArenaInfoFile},
			{testInstallPath, replaysDir, "12.4.0", tempArenaInfoFile},
			{"ほげほげ", replaysDir, tempArenaInfoFile},
		}

		for _, path := range paths {
			func(path []string) {
				defer os.RemoveAll(path[0])

				filePath := filepath.Join(path...)
				err := writeJSON(filePath, expected)
				assert.NoError(t, err)

				fetcher := NewLocalFile()
				actual, err := fetcher.ReadTempArenaInfo(path[0])
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			}(path)
		}
	})

	t.Run("正常系_該当ファイルが複数存在する場合_最新を返す", func(t *testing.T) {
		older := model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
				{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
				{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
			},
			DateTime:   "22.05.2022 12:34:56", // older than expected
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		expected := model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
				{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
				{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
			},
			DateTime:   "22.05.2023 12:34:56",
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		defer os.RemoveAll(testInstallPath)

		err := writeJSON(filepath.Join(testInstallPath, replaysDir, tempArenaInfoFile), older)
		require.NoError(t, err)
		err = writeJSON(filepath.Join(testInstallPath, replaysDir, "12.4.0", tempArenaInfoFile), expected)
		require.NoError(t, err)

		fetcher := NewLocalFile()
		actual, err := fetcher.ReadTempArenaInfo(testInstallPath)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("異常系_該当ファイルなし", func(t *testing.T) {
		paths := []string{
			filepath.Join(testInstallPath, replaysDir, "hoge.wowsreplay"),
			filepath.Join(testInstallPath, replaysDir, "12.4.0", "hoge.wowsreplay"),
		}

		for _, path := range paths {
			func(path string) {
				defer os.RemoveAll(testInstallPath)

				err := writeJSON(path, model.TempArenaInfo{})
				require.NoError(t, err)

				fetcher := NewLocalFile()
				_, err = fetcher.ReadTempArenaInfo(testInstallPath)

				assert.True(t, failure.Is(err, apperr.FileNotExist))
			}(path)
		}
	})
	t.Run("異常系_replayフォルダなし", func(t *testing.T) {
		err := os.Mkdir(testInstallPath, os.ModePerm)
		require.NoError(t, err)
		defer os.RemoveAll(testInstallPath)

		fetcher := NewLocalFile()
		_, err = fetcher.ReadTempArenaInfo(testInstallPath)

		assert.True(t, failure.Is(err, apperr.ReplayDirNotFoundError))
	})
}
