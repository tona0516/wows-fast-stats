package infra

import (
	"io/fs"
	"path/filepath"
	"testing"
	"wfs/backend/core"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplayReader_TempArenaInfo(t *testing.T) {
	const (
		replaysDir        = "replays"
		tempArenaInfoFile = "tempArenaInfo.json"
	)

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := core.TempArenaInfo{
			Vehicles: []core.Vehicle{
				{ShipID: 1, Relation: 0, Name: "player_1"},
				{ShipID: 2, Relation: 1, Name: "player_2"},
				{ShipID: 3, Relation: 2, Name: "player_3"},
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

			injector := do.New()
			instance, err := NewReplayReader(injector)
			require.NoError(t, err)
			actual, err := instance.TempArenaInfo(testInstallDir)

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("正常系_該当ファイルが複数存在する場合_最新を返す", func(t *testing.T) {
		t.Parallel()

		older := core.TempArenaInfo{
			Vehicles: []core.Vehicle{
				{ShipID: 1, Relation: 0, Name: "player_1"},
				{ShipID: 2, Relation: 1, Name: "player_2"},
				{ShipID: 3, Relation: 2, Name: "player_3"},
			},
			DateTime:   "22.05.2022 12:34:56", // older than expected
			MapID:      10,
			MatchGroup: "pvp",
			PlayerName: "player_1",
		}

		expected := core.TempArenaInfo{
			Vehicles: []core.Vehicle{
				{ShipID: 1, Relation: 0, Name: "player_1"},
				{ShipID: 2, Relation: 1, Name: "player_2"},
				{ShipID: 3, Relation: 2, Name: "player_3"},
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

		injector := do.New()
		instance, err := NewReplayReader(injector)
		require.NoError(t, err)
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
			err = writeJSON(path, core.TempArenaInfo{})
			require.NoError(t, err)

			injector := do.New()
			instance, err := NewReplayReader(injector)
			require.NoError(t, err)
			_, err = instance.TempArenaInfo(testInstallDir)

			assert.Error(t, err)
		}
	})
	t.Run("異常系_replayフォルダなし", func(t *testing.T) {
		testInstallDir := t.TempDir()

		injector := do.New()
		instance, err := NewReplayReader(injector)
		require.NoError(t, err)
		_, err = instance.TempArenaInfo(testInstallDir)

		assert.Error(t, err, fs.ErrNotExist)
	})
}
