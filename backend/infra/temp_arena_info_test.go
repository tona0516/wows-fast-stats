package infra_test

import (
	"changeme/backend/infra"
	"changeme/backend/vo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempArenaInfo_Get_正常系(t *testing.T) {
	t.Parallel()

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

	actual, err := tempArenaInfo.Get("testdata")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTempArenaInfo_Get_異常系_不正なパス(t *testing.T) {
	t.Parallel()

	tempArenaInfo := infra.NewTempArenaInfo()
	installPath := "./invalid"
	actual, err := tempArenaInfo.Get(installPath)

	assert.Error(t, err)
	assert.Empty(t, actual)
}
