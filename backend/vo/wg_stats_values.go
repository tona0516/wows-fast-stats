package vo

type WGStatsValues struct {
	Wins            uint `json:"wins"`
	Battles         uint `json:"battles"`
	DamageDealt     uint `json:"damage_dealt"`
	Frags           uint `json:"frags"`
	SurvivedWins    uint `json:"survived_wins"`
	SurvivedBattles uint `json:"survived_battles"`
	Xp              uint `json:"xp"`
	MainBattery     struct {
		Hits  uint `json:"hits"`
		Shots uint `json:"shots"`
	} `json:"main_battery"`
	Torpedoes struct {
		Hits  uint `json:"hits"`
		Shots uint `json:"shots"`
	} `json:"torpedoes"`
}
