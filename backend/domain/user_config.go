package domain

type UserConfig struct {
	Version      int          `json:"version"`
	InstallPath  string       `json:"install_path"`
	ZoomRate     int          `json:"zoom_rate"`
	StatsExtra   string       `json:"stats_extra"`
	IsSendReport bool         `json:"is_send_report"`
	Column       ColumnConfig `json:"column"`
}

type ColumnConfig struct {
	Player PlayerColumnConfig `json:"player"`
	Ship   ShipColumnConfig   `json:"ship"`
	Stats  StatsColumnConfig  `json:"stats"`
}

type PlayerColumnConfig struct {
	EnableNationFlag bool   `json:"enable_nation_flag"`
	ColorPattern     string `json:"color_pattern"`
}

type ShipColumnConfig struct {
	EnableNationFlag bool `json:"enable_nation_flag"`
	IsColored        bool `json:"is_colored"`
}

type StatsColumnConfig struct {
	Battles           DetailStatsColumnConfig `json:"battles"`
	Damage            DetailStatsColumnConfig `json:"damage"`
	MaxDamage         DetailStatsColumnConfig `json:"max_damage"`
	WinRate           DetailStatsColumnConfig `json:"win_rate"`
	SurvivedRate      DetailStatsColumnConfig `json:"survived_rate"`
	KdRate            DetailStatsColumnConfig `json:"kd_rate"`
	Kill              DetailStatsColumnConfig `json:"kill"`
	Exp               DetailStatsColumnConfig `json:"exp"`
	PR                DetailStatsColumnConfig `json:"pr"`
	HitRate           DetailStatsColumnConfig `json:"hit_rate"`
	PlanesKilled      DetailStatsColumnConfig `json:"planes_killed"`
	PlatoonRate       DetailStatsColumnConfig `json:"platoon_rate"`
	EfficiencyBadge   DetailStatsColumnConfig `json:"efficiency_badge"`
	ThreatLevel       DetailStatsColumnConfig `json:"threat_level"`
	AvgTier           DetailStatsColumnConfig `json:"avg_tier"`
	UsingShipTypeRate DetailStatsColumnConfig `json:"using_ship_type_rate"`
	UsingTierRate     DetailStatsColumnConfig `json:"using_tier_rate"`
}

type DetailStatsColumnConfig struct {
	IsShowShip    bool `json:"is_show_ship"`
	IsShowOverall bool `json:"is_show_overall"`
	Digit         int  `json:"digit"`
}
