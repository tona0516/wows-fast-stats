package core

type Pref struct {
	Version      int          `json:"version"`
	InstallPath  string       `json:"installPath"`
	ZoomRate     int          `json:"zoomRate"`
	StatsExtra   string       `json:"statsExtra"`
	IsSendReport bool         `json:"isSendReport"`
	Column       ColumnConfig `json:"column"`
}

type ColumnConfig struct {
	Player PlayerColumnConfig `json:"player"`
	Ship   ShipColumnConfig   `json:"ship"`
	Stats  StatsColumnConfig  `json:"stats"`
}

type PlayerColumnConfig struct {
	EnableNationFlag bool   `json:"enableNationFlag"`
	ColorPattern     string `json:"colorPattern"`
}

type ShipColumnConfig struct {
	EnableNationFlag bool `json:"enableNationFlag"`
	IsColored        bool `json:"isColored"`
}

type StatsColumnConfig struct {
	Battles           DetailStatsColumnConfig `json:"battles"`
	Damage            DetailStatsColumnConfig `json:"damage"`
	MaxDamage         DetailStatsColumnConfig `json:"maxDamage"`
	WinRate           DetailStatsColumnConfig `json:"winRate"`
	SurvivedRate      DetailStatsColumnConfig `json:"survivedRate"`
	KdRate            DetailStatsColumnConfig `json:"kdRate"`
	Kill              DetailStatsColumnConfig `json:"kill"`
	Exp               DetailStatsColumnConfig `json:"exp"`
	PR                DetailStatsColumnConfig `json:"pr"`
	HitRate           DetailStatsColumnConfig `json:"hitRate"`
	PlanesKilled      DetailStatsColumnConfig `json:"planesKilled"`
	PlatoonRate       DetailStatsColumnConfig `json:"platoonRate"`
	EfficiencyBadge   DetailStatsColumnConfig `json:"efficiencyBadge"`
	ThreatLevel       DetailStatsColumnConfig `json:"threatLevel"`
	AvgTier           DetailStatsColumnConfig `json:"avgTier"`
	UsingShipTypeRate DetailStatsColumnConfig `json:"usingShipTypeRate"`
	UsingTierRate     DetailStatsColumnConfig `json:"usingTierRate"`
}

type DetailStatsColumnConfig struct {
	IsShowShip    bool `json:"isShowShip"`
	IsShowOverall bool `json:"isShowOverall"`
	Digit         int  `json:"digit"`
}

func DefaultPref() Pref {
	return Pref{
		Version:      1,
		InstallPath:  "",
		ZoomRate:     100,
		StatsExtra:   StatsPatternPvPAll,
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
