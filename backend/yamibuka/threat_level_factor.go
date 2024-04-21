package yamibuka

import (
	"wfs/backend/data"
)

type ThreatLevelFactor struct {
	accountID        int
	tempArenaInfo    data.TempArenaInfo
	warships         data.Warships
	shipID           int
	shipBattles      uint
	shipDamage       float64
	shipWinRate      float64
	shipSurvivedRate float64
	shipPlanesKilled float64
	overallBattles   uint
	overallDamage    float64
	overallWinRate   float64
	overallKill      float64
	overallKdRate    float64
}

func NewThreatLevelFactor(
	accountID int,
	tempArenaInfo data.TempArenaInfo,
	warships data.Warships,
	shipID int,
	shipBattles uint,
	shipDamage float64,
	shipWinRate float64,
	shipSurvivedRate float64,
	shipPlanesKilled float64,
	overallBattles uint,
	overallDamage float64,
	overallWinRate float64,
	overallKill float64,
	overallKdRate float64,
) ThreatLevelFactor {
	return ThreatLevelFactor{
		accountID:        accountID,
		tempArenaInfo:    tempArenaInfo,
		warships:         warships,
		shipID:           shipID,
		shipBattles:      shipBattles,
		shipDamage:       shipDamage,
		shipWinRate:      shipWinRate,
		shipSurvivedRate: shipSurvivedRate,
		shipPlanesKilled: shipPlanesKilled,
		overallBattles:   overallBattles,
		overallDamage:    overallDamage,
		overallWinRate:   overallWinRate,
		overallKill:      overallKill,
		overallKdRate:    overallKdRate,
	}
}
