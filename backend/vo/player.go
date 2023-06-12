package vo

type PlayerInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Clan     Clan   `json:"clan"`
	IsHidden bool   `json:"is_hidden"`
}

type ShipInfo struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Nation    string   `json:"nation"`
	Tier      uint     `json:"tier"`
	Type      ShipType `json:"type"`
	AvgDamage float64  `json:"avg_damage"`
}

type PlayerStats struct {
	ShipStats    ShipStats    `json:"ship"`
	OverallStats OverallStats `json:"overall"`
}

type ShipStats struct {
	Battles            uint    `json:"battles"`
	Damage             float64 `json:"damage"`
	WinRate            float64 `json:"win_rate"`
	WinSurvivedRate    float64 `json:"win_survived_rate"`
	LoseSurvivedRate   float64 `json:"lose_survived_rate"`
	KdRate             float64 `json:"kd_rate"`
	Kill               float64 `json:"kill"`
	Death              float64 `json:"death"`
	Exp                float64 `json:"exp"`
	MainBatteryHitRate float64 `json:"main_battery_hit_rate"`
	TorpedoesHitRate   float64 `json:"torpedoes_hit_rate"`
	PR                 float64 `json:"pr"`
}

type OverallStats struct {
	Battles           uint          `json:"battles"`
	Damage            float64       `json:"damage"`
	WinRate           float64       `json:"win_rate"`
	WinSurvivedRate   float64       `json:"win_survived_rate"`
	LoseSurvivedRate  float64       `json:"lose_survived_rate"`
	KdRate            float64       `json:"kd_rate"`
	Kill              float64       `json:"kill"`
	Death             float64       `json:"death"`
	Exp               float64       `json:"exp"`
	AvgTier           float64       `json:"avg_tier"`
	UsingShipTypeRate ShipTypeGroup `json:"using_ship_type_rate"`
	UsingTierRate     TierGroup     `json:"using_tier_rate"`
}

type Player struct {
	PlayerInfo PlayerInfo  `json:"player_info"`
	ShipInfo   ShipInfo    `json:"ship_info"`
	PvPSolo    PlayerStats `json:"pvp_solo"`
	PvPAll     PlayerStats `json:"pvp_all"`
}
