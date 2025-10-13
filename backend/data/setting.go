package data

type RequiredSetting struct {
	Version     int    `json:"version"`
	InstallPath string `json:"install_path"`
}

type OptionalSetting struct {
	Version      int    `json:"version"`
	ZoomRate     int    `json:"zoom_rate"`
	StatsExtra   string `json:"stats_extra"`
	IsSendReport bool   `json:"is_send_report"`
}

type BasicColumnSetting struct {
	Version int                 `json:"version"`
	Player  PlayerColumnSetting `json:"player"`
	Ship    ShipColumnSetting   `json:"ship"`
}

type PlayerColumnSetting struct {
	EnableNationFlag bool   `json:"enable_nation_flag"`
	ColorPattern     string `json:"color_pattern"`
}

type ShipColumnSetting struct {
	EnableNationFlag bool `json:"enable_nation_flag"`
	IsColored        bool `json:"is_colored"`
}

type StatsColumnSettings struct {
	Version           int                `json:"version"`
	Battles           StatsColumnSetting `json:"battles"`
	Damage            StatsColumnSetting `json:"damage"`
	MaxDamage         StatsColumnSetting `json:"max_damage"`
	WinRate           StatsColumnSetting `json:"win_rate"`
	SurvivedRate      StatsColumnSetting `json:"survived_rate"`
	KdRate            StatsColumnSetting `json:"kd_rate"`
	Kill              StatsColumnSetting `json:"kill"`
	Exp               StatsColumnSetting `json:"exp"`
	PR                StatsColumnSetting `json:"pr"`
	HitRate           StatsColumnSetting `json:"hit_rate"`
	PlanesKilled      StatsColumnSetting `json:"planes_killed"`
	PlatoonRate       StatsColumnSetting `json:"platoon_rate"`
	EfficiencyBadge   StatsColumnSetting `json:"efficiency_badge"`
	ThreatLevel       StatsColumnSetting `json:"threat_level"`
	AvgTier           StatsColumnSetting `json:"avg_tier"`
	UsingShipTypeRate StatsColumnSetting `json:"using_ship_type_rate"`
	UsingTierRate     StatsColumnSetting `json:"using_tier_rate"`
}

type StatsColumnSetting struct {
	IsShowShip    bool `json:"is_show_ship"`
	IsShowOverall bool `json:"is_show_overall"`
	Digit         int  `json:"digit"`
}
