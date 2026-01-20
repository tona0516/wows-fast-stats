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
