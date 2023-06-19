package infra_test

import (
	"changeme/backend/infra"
	"changeme/backend/vo"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testInstallPath = "testdata"

//nolint:paralleltest
func TestTempArenaInfo_Get_正常系(t *testing.T) {
	tempArenaInfo := infra.NewTempArenaInfo()
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
		filepath.Join(testInstallPath, infra.ReplayDir, infra.TempArenaInfoName),
		filepath.Join(testInstallPath, infra.ReplayDir, "12.4.0", infra.TempArenaInfoName),
	}

	for _, path := range paths {
		func(path string) {
			defer os.RemoveAll(testInstallPath)

			err := create(expected, path)
			assert.NoError(t, err)

			actual, err := tempArenaInfo.Get(testInstallPath)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}(path)
	}
}

//nolint:paralleltest
func TestTempArenaInfo_Get_正常系_該当ファイルが複数存在する場合_最新を返す(t *testing.T) {
	tempArenaInfo := infra.NewTempArenaInfo()

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

	err := create(older, filepath.Join(installPath, infra.ReplayDir, infra.TempArenaInfoName))
	assert.NoError(t, err)
	err = create(expected, filepath.Join(installPath, infra.ReplayDir, "12.4.0", infra.TempArenaInfoName))
	assert.NoError(t, err)

	actual, err := tempArenaInfo.Get(installPath)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

//nolint:paralleltest
func TestTempArenaInfo_Get_異常系_該当ファイルなし(t *testing.T) {
	tempArenaInfo := infra.NewTempArenaInfo()

	installPath := "testdata"
	paths := []string{
		filepath.Join(installPath, infra.ReplayDir, "hoge.wowsreplay"),
		filepath.Join(installPath, infra.ReplayDir, "12.4.0", "hoge.wowsreplay"),
	}

	for _, path := range paths {
		func(path string) {
			defer os.RemoveAll(testInstallPath)

			err := create(vo.TempArenaInfo{}, path)
			assert.NoError(t, err)

			_, err = tempArenaInfo.Get(installPath)
			assert.Error(t, err)
		}(path)
	}
}

func create(tempArenaInfo vo.TempArenaInfo, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), fs.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(tempArenaInfo)
}
