package infra

import (
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/vo"

	"github.com/stretchr/testify/assert"
)

const testInstallPath = "testdata"

//nolint:paralleltest
func TestTempArenaInfo_Get_正常系(t *testing.T) {
	tempArenaInfo := NewTempArenaInfo()
	expected := vo.TempArenaInfo{
		Vehicles: []vo.Vehicle{
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
		filepath.Join(testInstallPath, ReplayDir, TempArenaInfoName),
		filepath.Join(testInstallPath, ReplayDir, "12.4.0", TempArenaInfoName),
	}

	for _, path := range paths {
		func(path string) {
			defer os.RemoveAll(testInstallPath)

			err := writeJSON(path, expected)
			assert.NoError(t, err)

			actual, err := tempArenaInfo.Get(testInstallPath)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}(path)
	}
}

//nolint:paralleltest
func TestTempArenaInfo_Get_正常系_該当ファイルが複数存在する場合_最新を返す(t *testing.T) {
	tempArenaInfo := NewTempArenaInfo()

	older := vo.TempArenaInfo{
		Vehicles: []vo.Vehicle{
			{ShipID: 1, Relation: 0, ID: 100, Name: "player_1"},
			{ShipID: 2, Relation: 1, ID: 200, Name: "player_2"},
			{ShipID: 3, Relation: 2, ID: 300, Name: "player_3"},
		},
		DateTime:   "22.05.2022 12:34:56", // older than expected
		MapID:      10,
		MatchGroup: "pvp",
		PlayerName: "player_1",
	}

	expected := vo.TempArenaInfo{
		Vehicles: []vo.Vehicle{
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

	err := writeJSON(filepath.Join(installPath, ReplayDir, TempArenaInfoName), older)
	assert.NoError(t, err)
	err = writeJSON(filepath.Join(installPath, ReplayDir, "12.4.0", TempArenaInfoName), expected)
	assert.NoError(t, err)

	actual, err := tempArenaInfo.Get(installPath)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

//nolint:paralleltest
func TestTempArenaInfo_Get_異常系_該当ファイルなし(t *testing.T) {
	tempArenaInfo := NewTempArenaInfo()

	installPath := "testdata"
	paths := []string{
		filepath.Join(installPath, ReplayDir, "hoge.wowsreplay"),
		filepath.Join(installPath, ReplayDir, "12.4.0", "hoge.wowsreplay"),
	}

	for _, path := range paths {
		func(path string) {
			defer os.RemoveAll(testInstallPath)

			err := writeJSON(path, vo.TempArenaInfo{})
			assert.NoError(t, err)

			_, err = tempArenaInfo.Get(installPath)
			assert.Error(t, err)
		}(path)
	}
}
