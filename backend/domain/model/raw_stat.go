package model

type RawStats map[int]RawStat

type RawStat struct {
	Ship    map[int]ShipStat
	Overall OverallStat
}

type ShipStat struct {
	Pvp     ShipStatsValues
	PvpSolo ShipStatsValues
	PvpDiv2 struct {
		Battles uint
	}
	PvpDiv3 struct {
		Battles uint
	}
	RankSolo ShipStatsValues
}

type OverallStat struct {
	IsHidden bool
	Pvp      OverallStatsValues
	PvpSolo  OverallStatsValues
	PvpDiv2  OverallStatsValues
	PvpDiv3  OverallStatsValues
	RankSolo OverallStatsValues
}

type OverallStatsValues struct {
	Wins                 uint
	Battles              uint
	DamageDealt          uint
	MaxDamageDealt       uint
	MaxDamageDealtShipID int
	Frags                uint
	SurvivedWins         uint
	SurvivedBattles      uint
	Xp                   uint
}

type ShipStatsValues struct {
	Wins            uint
	Battles         uint
	DamageDealt     uint
	MaxDamageDealt  uint
	Frags           uint
	SurvivedWins    uint
	SurvivedBattles uint
	Xp              uint
	MainBattery     struct {
		Hits  uint
		Shots uint
	}
	Torpedoes struct {
		Hits  uint
		Shots uint
	}
	PlanesKilled uint
}
