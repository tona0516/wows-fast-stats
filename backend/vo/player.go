package vo

type ShipInfo struct {
	Name     string `json:"name"`
	Nation   string `json:"nation"`
	Tier     uint `json:"tier"`
	Type     string `json:"type"`
	StatsURL string `json:"stats_url"`
}

type ShipStats struct {
	Battles        uint `json:"battles"`
	Damage      float64 `json:"damage"`
	WinRate        float64 `json:"win_rate"`
    WinSurvivedRate float64 `json:"win_survived_rate"`
    LoseSurvivedRate float64 `json:"lose_survived_rate"`
	KdRate         float64 `json:"kd_rate"`
    Exp float64 `json:"exp"`
	PR float64 `json:"pr"`
}

type PlayerInfo struct {
    ID       int `json:"id"`
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	IsHidden bool `json:"is_hidden"`
	StatsURL string `json:"stats_url"`
}

type PlayerStats struct {
	Battles   uint `json:"battles"`
	Damage float64 `json:"damage"`
	WinRate   float64 `json:"win_rate"`
    WinSurvivedRate float64 `json:"win_survived_rate"`
    LoseSurvivedRate float64 `json:"lose_survived_rate"`
	KdRate    float64 `json:"kd_rate"`
    Exp float64 `json:"exp"`
	AvgTier   float64 `json:"avg_tier"`
    UsingShipTypeRate ShipTypeGroup[float64] `json:"using_ship_type_rate"`
    UsingTierRate TierGroup[float64] `json:"using_tier_rate"`
}

type Player struct {
	ShipInfo    ShipInfo `json:"ship_info"`
	ShipStats   ShipStats `json:"ship_stats"`
	PlayerInfo  PlayerInfo `json:"player_info"`
	PlayerStats PlayerStats `json:"player_stats"`
}
