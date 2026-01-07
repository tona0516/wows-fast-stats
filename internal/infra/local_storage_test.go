package infra

import (
	"io/fs"
	"path/filepath"
	"testing"
	"wfs/internal/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage_TempArenaInfo(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := data.TempArenaInfo{
			Vehicles: []data.Vehicle{
				{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
				{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
				{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
			},
			DateTime:   "22.05.2023 12:34:56",
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		testInstallDir := t.TempDir()
		paths := []string{
			filepath.Join(testInstallDir, replaysDir, tempArenaInfoFile),
			filepath.Join(testInstallDir, replaysDir, "12.4.0", tempArenaInfoFile),
		}
		for _, path := range paths {
			err := writeJSON(path, expected)
			require.NoError(t, err)

			instance := NewLocalStorage("", "")
			actual, err := instance.TempArenaInfo(testInstallDir)

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("正常系_該当ファイルが複数存在する場合_最新を返す", func(t *testing.T) {
		t.Parallel()

		older := data.TempArenaInfo{
			Vehicles: []data.Vehicle{
				{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
				{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
				{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
			},
			DateTime:   "22.05.2022 12:34:56", // older than expected
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		expected := data.TempArenaInfo{
			Vehicles: []data.Vehicle{
				{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
				{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
				{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
			},
			DateTime:   "22.05.2023 12:34:56",
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		testInstallDir := t.TempDir()
		var err error
		err = writeJSON(filepath.Join(testInstallDir, replaysDir, tempArenaInfoFile), older)
		require.NoError(t, err)
		err = writeJSON(filepath.Join(testInstallDir, replaysDir, "12.4.0", tempArenaInfoFile), expected)
		require.NoError(t, err)

		instance := NewLocalStorage("", "")
		actual, err := instance.TempArenaInfo(testInstallDir)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("異常系_該当ファイルなし", func(t *testing.T) {
		t.Parallel()

		testInstallDir := t.TempDir()
		paths := []string{
			filepath.Join(testInstallDir, replaysDir, "hoge.wowsreplay"),
			filepath.Join(testInstallDir, replaysDir, "12.4.0", "hoge.wowsreplay"),
		}

		for _, path := range paths {
			var err error
			err = writeJSON(path, data.TempArenaInfo{})
			require.NoError(t, err)

			instance := NewLocalStorage("", "")
			_, err = instance.TempArenaInfo(testInstallDir)

			assert.Error(t, err)
		}
	})
	t.Run("異常系_replayフォルダなし", func(t *testing.T) {
		testInstallDir := t.TempDir()

		instance := NewLocalStorage("", "")
		_, err := instance.TempArenaInfo(testInstallDir)

		assert.Error(t, err, fs.ErrNotExist)
	})
}
