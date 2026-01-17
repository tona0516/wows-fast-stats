package core

type ThreatLevelInput struct {
	Vehicles         []Vehicle
	Warships         Warships
	ShipID           ShipID
	ShipBattles      uint
	ShipDamage       float64
	ShipWinRate      float64
	ShipSurvivedRate float64
	ShipPlanesKilled float64
	OverallBattles   uint
	OverallDamage    float64
	OverallWinRate   float64
	OverallKill      float64
	OverallKdRate    float64
}

func NewThreatLevelInput(
	vehicles []Vehicle,
	warships Warships,
	shipID ShipID,
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
) ThreatLevelInput {
	return ThreatLevelInput{
		Vehicles:         vehicles,
		Warships:         warships,
		ShipID:           shipID,
		ShipBattles:      shipBattles,
		ShipDamage:       shipDamage,
		ShipWinRate:      shipWinRate,
		ShipSurvivedRate: shipSurvivedRate,
		ShipPlanesKilled: shipPlanesKilled,
		OverallBattles:   overallBattles,
		OverallDamage:    overallDamage,
		OverallWinRate:   overallWinRate,
		OverallKill:      overallKill,
		OverallKdRate:    overallKdRate,
	}
}
