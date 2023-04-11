package vo

type PlayerShipInfo struct {
	Name     string `json:"name"`
	Nation   string `json:"nation"`
	Tier     uint `json:"tier"`
	Type     string `json:"type"`
	StatsURL string `json:"stats_url"`
}

type PlayerShipStats struct {
	Battles        uint `json:"battles"`
	AvgDamage      uint `json:"avg_damage"`
	AvgExp         uint `json:"avg_exp"`
	WinRate        float64 `json:"win_rate"`
	KdRate         float64 `json:"kd_rate"`
	CombatPower    uint `json:"combat_power"`
	PersonalRating uint `json:"personal_rating"`
}

type PlayerPlayerInfo struct {
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	IsHidden bool `json:"is_hidden"`
	StatsURL string `json:"stats_url"`
}

type PlayerPlayerStats struct {
	Battles   uint `json:"battles"`
	AvgDamage uint `json:"avg_damage"`
	AvgExp    uint `json:"avg_exp"`
	WinRate   float64 `json:"win_rate"`
	KdRate    float64 `json:"kd_rate"`
	AvgTier   float64 `json:"avg_tier"`
}

type Player struct {
	ShipInfo    PlayerShipInfo `json:"player_ship_info"`
	ShipStats   PlayerShipStats `json:"player_ship_stats"`
	PlayerInfo  PlayerPlayerInfo `json:"player_player_info"`
	PlayerStats PlayerPlayerStats `json:"player_player_stats"`
}

type Team struct {
    Friends []Player `json:"friends"`
    Enemies []Player `json:"enemies"`
}
