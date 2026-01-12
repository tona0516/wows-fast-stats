package data

type ThreatLevelInput struct {
	AccountID        int
	Vehicles         []Vehicle
	Warships         Warships
	ShipID           int
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
	accountID int,
	vehicles []Vehicle,
	warships Warships,
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
) ThreatLevelInput {
	return ThreatLevelInput{
		AccountID:        accountID,
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
