package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempArenaInfo_AccountNames(t *testing.T) {
	t.Parallel()
	vehicles := []Vehicle{
		{ShipID: 1, Relation: 0, ID: 1001, Name: "Player1"},
		{ShipID: 2, Relation: 0, ID: 1002, Name: ":Bot1:"},
		{ShipID: 3, Relation: 1, ID: 1003, Name: "Player2"},
		{ShipID: 4, Relation: 0, ID: 1004, Name: "IDS_OP1"},
	}

	info := &TempArenaInfo{
		Vehicles: vehicles,
	}

	expectedNames := []string{"Player1", "Player2"}
	actualNames := info.AccountNames()

	assert.ElementsMatch(t, expectedNames, actualNames)
}

func TestTempArenaInfo_Unixtime(t *testing.T) {
	t.Parallel()
	info := &TempArenaInfo{
		DateTime: "22.05.2023 12:34:56",
	}

	var expected int64 = 1684726496 // 2023-05-22 12:34:56
	actual := info.Unixtime()

	assert.Equal(t, expected, actual)
}

func TestTempArenaInfo_BattleArena(t *testing.T) {
	t.Parallel()
	w := WGBattleArenas{
		1: {Name: "Arena1"},
		2: {Name: "Arena2"},
	}

	info := &TempArenaInfo{
		MapID: 2,
	}

	expectedBattleArena := "Arena2"
	actualBattleArena := info.BattleArena(w)

	assert.Equal(t, expectedBattleArena, actualBattleArena)
}

func TestTempArenaInfo_BattleType(t *testing.T) {
	t.Parallel()
	w := WGBattleTypes{
		"RANDOM": {Name: "Random Battle"},
		"COOP":   {Name: "Co-op Battle"},
	}

	info := &TempArenaInfo{
		MatchGroup: "random",
	}

	expectedBattleType := "RandomBattle"
	actualBattleType := info.BattleType(w)

	assert.Equal(t, expectedBattleType, actualBattleType)
}
