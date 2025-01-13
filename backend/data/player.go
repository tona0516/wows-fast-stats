package data

import "wfs/backend/domain/model"

type PlayerInfo struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Clan     model.Clan `json:"clan"`
	IsHidden bool       `json:"is_hidden"`
}

type PlayerStats struct {
	ShipStats    ShipStats    `json:"ship"`
	OverallStats OverallStats `json:"overall"`
}

type ShipStats struct {
	Battles            uint      `json:"battles"`
	Damage             float64   `json:"damage"`
	MaxDamage          MaxDamage `json:"max_damage"`
	WinRate            float64   `json:"win_rate"`
	WinSurvivedRate    float64   `json:"win_survived_rate"`
	LoseSurvivedRate   float64   `json:"lose_survived_rate"`
	KdRate             float64   `json:"kd_rate"`
	Kill               float64   `json:"kill"`
	Exp                float64   `json:"exp"`
	PR                 float64   `json:"pr"`
	MainBatteryHitRate float64   `json:"main_battery_hit_rate"`
	TorpedoesHitRate   float64   `json:"torpedoes_hit_rate"`
	PlanesKilled       float64   `json:"planes_killed"`
	PlatoonRate        float64   `json:"platoon_rate"`
}

type OverallStats struct {
	Battles           uint              `json:"battles"`
	Damage            float64           `json:"damage"`
	MaxDamage         MaxDamage         `json:"max_damage"`
	WinRate           float64           `json:"win_rate"`
	WinSurvivedRate   float64           `json:"win_survived_rate"`
	LoseSurvivedRate  float64           `json:"lose_survived_rate"`
	KdRate            float64           `json:"kd_rate"`
	Kill              float64           `json:"kill"`
	Exp               float64           `json:"exp"`
	PR                float64           `json:"pr"`
	ThreatLevel       model.ThreatLevel `json:"threat_level"`
	AvgTier           float64           `json:"avg_tier"`
	UsingShipTypeRate ShipTypeGroup     `json:"using_ship_type_rate"`
	UsingTierRate     TierGroup         `json:"using_tier_rate"`
	PlatoonRate       float64           `json:"platoon_rate"`
}

type Player struct {
	PlayerInfo PlayerInfo    `json:"player_info"`
	Warship    model.Warship `json:"warship"`
	PvPSolo    PlayerStats   `json:"pvp_solo"`
	PvPAll     PlayerStats   `json:"pvp_all"`
	RankSolo   PlayerStats   `json:"rank_solo"`
}
