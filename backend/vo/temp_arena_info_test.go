package vo

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

func TestTempArenaInfo_FormattedDateTime(t *testing.T) {
	t.Parallel()
	info := &TempArenaInfo{
		DateTime: "22.05.2023 12:34:56",
	}

	expectedFormattedDateTime := "2023-05-22 12:34:56"
	actualFormattedDateTime := info.FormattedDateTime()

	assert.Equal(t, expectedFormattedDateTime, actualFormattedDateTime)
}

func TestTempArenaInfo_BattleArena(t *testing.T) {
	t.Parallel()
	battleArenas := WGBattleArenas{
		Data: map[int]WGBattleArenasData{
			1: {Name: "Arena1"},
			2: {Name: "Arena2"},
		},
	}

	info := &TempArenaInfo{
		MapID: 2,
	}

	expectedBattleArena := "Arena2"
	actualBattleArena := info.BattleArena(battleArenas)

	assert.Equal(t, expectedBattleArena, actualBattleArena)
}

func TestTempArenaInfo_BattleType(t *testing.T) {
	t.Parallel()
	battleTypes := WGBattleTypes{
		Data: map[string]WGBattleTypesData{
			"RANDOM": {Name: "Random Battle"},
			"COOP":   {Name: "Co-op Battle"},
		},
	}

	info := &TempArenaInfo{
		MatchGroup: "random",
	}

	expectedBattleType := "RandomBattle"
	actualBattleType := info.BattleType(battleTypes)

	assert.Equal(t, expectedBattleType, actualBattleType)
}
