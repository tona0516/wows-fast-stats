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
	AvgDamage      float64 `json:"avg_damage"`
	WinRate        float64 `json:"win_rate"`
    WinSurvivedRate float64 `json:"win_survived_rate"`
    LoseSurvivedRate float64 `json:"lose_survived_rate"`
	KdRate         float64 `json:"kd_rate"`
    Exp float64 `json:"exp"`
	PersonalRating float64 `json:"personal_rating"`
}

type PlayerPlayerInfo struct {
    ID       int `json:"id"`
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	IsHidden bool `json:"is_hidden"`
	StatsURL string `json:"stats_url"`
}

type PlayerPlayerStats struct {
	Battles   uint `json:"battles"`
	AvgDamage float64 `json:"avg_damage"`
	WinRate   float64 `json:"win_rate"`
    WinSurvivedRate float64 `json:"win_survived_rate"`
    LoseSurvivedRate float64 `json:"lose_survived_rate"`
	KdRate    float64 `json:"kd_rate"`
    Exp float64 `json:"exp"`
	AvgTier   float64 `json:"avg_tier"`
    UsingShipTypeRate ShipTypeValue `json:"using_ship_type_rate"`
}

type Player struct {
	ShipInfo    PlayerShipInfo `json:"player_ship_info"`
	ShipStats   PlayerShipStats `json:"player_ship_stats"`
	PlayerInfo  PlayerPlayerInfo `json:"player_player_info"`
	PlayerStats PlayerPlayerStats `json:"player_player_stats"`
}
