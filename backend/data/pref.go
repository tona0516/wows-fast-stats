package data

type Pref struct {
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

func DefaultPref() Pref {
	return Pref{
		Version:      1,
		InstallPath:  "",
		ZoomRate:     100,
		StatsExtra:   "pvp_all",
		IsSendReport: true,
		Column: ColumnConfig{
			Player: PlayerColumnConfig{
				EnableNationFlag: false,
				ColorPattern:     "none",
			},
			Ship: ShipColumnConfig{
				EnableNationFlag: true,
				IsColored:        false,
			},
			Stats: StatsColumnConfig{
				Battles: DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				Damage: DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				MaxDamage: DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				WinRate: DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         1,
				},
				SurvivedRate: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				KdRate: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         2,
				},
				Kill: DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         2,
				},
				Exp: DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				PR: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				HitRate: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				PlanesKilled: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				PlatoonRate: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				EfficiencyBadge: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				ThreatLevel: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				AvgTier: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingShipTypeRate: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingTierRate: DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
			},
		},
	}
}
