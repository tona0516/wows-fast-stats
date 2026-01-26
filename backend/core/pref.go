package core

type Pref struct {
	Version        int        `json:"version"`
	GameClientPath string     `json:"gameClientPath"`
	ZoomRate       int        `json:"zoomRate"`
	StatsExtra     string     `json:"statsExtra"`
	IsSendReport   bool       `json:"isSendReport"`
	Column         ColumnPref `json:"column"`
}

type ColumnPref struct {
	Player PlayerColumnPref `json:"player"`
	Ship   ShipColumnPref   `json:"ship"`
	Stats  StatsColumnPref  `json:"stats"`
}

type PlayerColumnPref struct {
	EnableNationFlag bool   `json:"enableNationFlag"`
	ColorPattern     string `json:"colorPattern"`
}

type ShipColumnPref struct {
	EnableNationFlag bool `json:"enableNationFlag"`
	IsColored        bool `json:"isColored"`
}

type StatsColumnPref struct {
	Battles           DetailStatsColumnPref `json:"battles"`
	Damage            DetailStatsColumnPref `json:"damage"`
	MaxDamage         DetailStatsColumnPref `json:"maxDamage"`
	WinRate           DetailStatsColumnPref `json:"winRate"`
	SurvivedRate      DetailStatsColumnPref `json:"survivedRate"`
	KdRate            DetailStatsColumnPref `json:"kdRate"`
	Kill              DetailStatsColumnPref `json:"kill"`
	Exp               DetailStatsColumnPref `json:"exp"`
	PR                DetailStatsColumnPref `json:"pr"`
	HitRate           DetailStatsColumnPref `json:"hitRate"`
	PlanesKilled      DetailStatsColumnPref `json:"planesKilled"`
	PlatoonRate       DetailStatsColumnPref `json:"platoonRate"`
	EfficiencyBadge   DetailStatsColumnPref `json:"efficiencyBadge"`
	ThreatLevel       DetailStatsColumnPref `json:"threatLevel"`
	AvgTier           DetailStatsColumnPref `json:"avgTier"`
	UsingShipTypeRate DetailStatsColumnPref `json:"usingShipTypeRate"`
	UsingTierRate     DetailStatsColumnPref `json:"usingTierRate"`
}

type DetailStatsColumnPref struct {
	IsShowShip    bool `json:"isShowShip"`
	IsShowOverall bool `json:"isShowOverall"`
	Digit         int  `json:"digit"`
}

func DefaultPref() Pref {
	return Pref{
		Version:        1,
		GameClientPath: "",
		ZoomRate:       100,
		StatsExtra:     StatsPatternPvPAll,
		IsSendReport:   true,
		Column: ColumnPref{
			Player: PlayerColumnPref{
				EnableNationFlag: false,
				ColorPattern:     "none",
			},
			Ship: ShipColumnPref{
				EnableNationFlag: true,
				IsColored:        false,
			},
			Stats: StatsColumnPref{
				Battles: DetailStatsColumnPref{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				Damage: DetailStatsColumnPref{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				MaxDamage: DetailStatsColumnPref{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				WinRate: DetailStatsColumnPref{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         1,
				},
				SurvivedRate: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				KdRate: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         2,
				},
				Kill: DetailStatsColumnPref{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         2,
				},
				Exp: DetailStatsColumnPref{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				PR: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				HitRate: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				PlanesKilled: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				PlatoonRate: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				EfficiencyBadge: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				ThreatLevel: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				AvgTier: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingShipTypeRate: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingTierRate: DetailStatsColumnPref{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
			},
		},
	}
}
