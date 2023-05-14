package vo

type ShipInfo struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Nation    string   `json:"nation"`
	Tier      uint     `json:"tier"`
	Type      ShipType `json:"type"`
	AvgDamage float64  `json:"avg_damage"`
}

type ShipStats struct {
	Battles            uint    `json:"battles"`
	Damage             float64 `json:"damage"`
	WinRate            float64 `json:"win_rate"`
	WinSurvivedRate    float64 `json:"win_survived_rate"`
	LoseSurvivedRate   float64 `json:"lose_survived_rate"`
	KdRate             float64 `json:"kd_rate"`
	Exp                float64 `json:"exp"`
	MainBatteryHitRate float64 `json:"main_battery_hit_rate"`
	TorpedoesHitRate   float64 `json:"torpedoes_hit_rate"`
	PR                 float64 `json:"pr"`
}

type PlayerInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Clan     Clan   `json:"clan"`
	IsHidden bool   `json:"is_hidden"`
}

type OverallStats struct {
	Battles           uint          `json:"battles"`
	Damage            float64       `json:"damage"`
	WinRate           float64       `json:"win_rate"`
	WinSurvivedRate   float64       `json:"win_survived_rate"`
	LoseSurvivedRate  float64       `json:"lose_survived_rate"`
	KdRate            float64       `json:"kd_rate"`
	Exp               float64       `json:"exp"`
	AvgTier           float64       `json:"avg_tier"`
	UsingShipTypeRate ShipTypeGroup `json:"using_ship_type_rate"`
	UsingTierRate     TierGroup     `json:"using_tier_rate"`
}

type Player struct {
	ShipInfo     ShipInfo     `json:"ship_info"`
	ShipStats    ShipStats    `json:"ship_stats"`
	PlayerInfo   PlayerInfo   `json:"player_info"`
	OverallStats OverallStats `json:"overall_stats"`
}
