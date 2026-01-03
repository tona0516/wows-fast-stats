package data

type WGPlayerStatsValues struct {
	Wins                 uint `json:"wins"`
	Battles              uint `json:"battles"`
	DamageDealt          uint `json:"damage_dealt"`
	MaxDamageDealt       uint `json:"max_damage_dealt"`
	MaxDamageDealtShipID int  `json:"max_damage_dealt_ship_id"`
	Frags                uint `json:"frags"`
	SurvivedWins         uint `json:"survived_wins"`
	SurvivedBattles      uint `json:"survived_battles"`
	Xp                   uint `json:"xp"`
}
