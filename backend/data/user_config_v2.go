package data

func DefaultUserConfigV2() UserConfigV2 {
	return UserConfigV2{
		Version:         2,
		FontSize:        "medium",
		SendReport:      true,
		NotifyUpdatable: true,
		StatsPattern:    StatsPatternPvPAll,
		Display: UCDisplay{
			Ship: UCDisplayShip{
				PR:      true,
				Damage:  true,
				WinRate: true,
				Battles: true,
			},
			Overall: UCDisplayOverall{
				Damage:  true,
				WinRate: true,
				Battles: true,
			},
		},
		Color: UCColor{
			PlayerName: PlayerNameColorShip,
		},
		Digit: UCDigit{
			PR:                0,
			Damage:            0,
			WinRate:           1,
			KdRate:            2,
			Kill:              2,
			PlanesKilled:      1,
			Exp:               0,
			Battles:           0,
			SurvivedRate:      1,
			HitRate:           1,
			AvgTier:           2,
			UsingShipTypeRate: 1,
			UsingTierRate:     1,
			PlatoonRate:       1,
			ThreatLevel:       0,
		},
	}
}

type UserConfigV2 struct {
	Version uint `json:"version"`
	// required
	InstallPath string `json:"install_path"`
	// display
	FontSize         string    `json:"font_size"`
	Display          UCDisplay `json:"display"`
	Color            UCColor   `json:"color"`
	Digit            UCDigit   `json:"digit"`
	ShowLanguageFrag bool      `json:"show_language_frag"`
	// other
	SaveScreenshot    bool         `json:"save_screenshot"`
	SaveTempArenaInfo bool         `json:"save_temp_arena_info"`
	SendReport        bool         `json:"send_report"`
	NotifyUpdatable   bool         `json:"notify_updatable"`
	StatsPattern      StatsPattern `json:"stats_pattern"`
}

func FromUserConfigV1(v1 UserConfig) UserConfigV2 {
	return UserConfigV2{
		Version:     2,
		InstallPath: v1.InstallPath,
		FontSize:    v1.FontSize,
		Display: UCDisplay{
			Ship: UCDisplayShip(v1.Displays.Ship),
			Overall: UCDisplayOverall{
				ThreatLevel:       false,
				PR:                v1.Displays.Overall.PR,
				Damage:            v1.Displays.Overall.Damage,
				MaxDamage:         v1.Displays.Overall.MaxDamage,
				WinRate:           v1.Displays.Overall.WinRate,
				KdRate:            v1.Displays.Overall.KdRate,
				Kill:              v1.Displays.Overall.Kill,
				Exp:               v1.Displays.Overall.Exp,
				Battles:           v1.Displays.Overall.Battles,
				SurvivedRate:      v1.Displays.Overall.SurvivedRate,
				AvgTier:           v1.Displays.Overall.AvgTier,
				UsingShipTypeRate: v1.Displays.Overall.UsingShipTypeRate,
				UsingTierRate:     v1.Displays.Overall.UsingTierRate,
				PlatoonRate:       v1.Displays.Overall.PlatoonRate,
			},
		},
		Color: UCColor{
			PlayerName: v1.CustomColor.PlayerName,
		},
		Digit: UCDigit{
			PR:                v1.CustomDigit.PR,
			Damage:            v1.CustomDigit.Damage,
			MaxDamage:         v1.CustomDigit.MaxDamage,
			WinRate:           v1.CustomDigit.WinRate,
			KdRate:            v1.CustomDigit.KdRate,
			Kill:              v1.CustomDigit.Kill,
			PlanesKilled:      v1.CustomDigit.PlanesKilled,
			Exp:               v1.CustomDigit.Exp,
			Battles:           v1.CustomDigit.Battles,
			SurvivedRate:      v1.CustomDigit.SurvivedRate,
			HitRate:           v1.CustomDigit.HitRate,
			AvgTier:           v1.CustomDigit.AvgTier,
			UsingShipTypeRate: v1.CustomDigit.UsingShipTypeRate,
			UsingTierRate:     v1.CustomDigit.UsingTierRate,
			PlatoonRate:       v1.CustomDigit.PlatoonRate,
			ThreatLevel:       0,
		},
		SaveScreenshot:    v1.SaveScreenshot,
		SaveTempArenaInfo: v1.SaveTempArenaInfo,
		SendReport:        v1.SendReport,
		NotifyUpdatable:   v1.NotifyUpdatable,
		StatsPattern:      v1.StatsPattern,
	}
}

type UCDisplay struct {
	Ship    UCDisplayShip    `json:"ship"`
	Overall UCDisplayOverall `json:"overall"`
}

type UCDisplayShip struct {
	PR           bool `json:"pr"`
	Damage       bool `json:"damage"`
	MaxDamage    bool `json:"max_damage"`
	WinRate      bool `json:"win_rate"`
	KdRate       bool `json:"kd_rate"`
	Kill         bool `json:"kill"`
	PlanesKilled bool `json:"planes_killed"`
	Exp          bool `json:"exp"`
	Battles      bool `json:"battles"`
	SurvivedRate bool `json:"survived_rate"`
	HitRate      bool `json:"hit_rate"`
	PlatoonRate  bool `json:"platoon_rate"`
}

type UCDisplayOverall struct {
	PR                bool `json:"pr"`
	Damage            bool `json:"damage"`
	MaxDamage         bool `json:"max_damage"`
	WinRate           bool `json:"win_rate"`
	KdRate            bool `json:"kd_rate"`
	Kill              bool `json:"kill"`
	Exp               bool `json:"exp"`
	Battles           bool `json:"battles"`
	SurvivedRate      bool `json:"survived_rate"`
	AvgTier           bool `json:"avg_tier"`
	UsingShipTypeRate bool `json:"using_ship_type_rate"`
	UsingTierRate     bool `json:"using_tier_rate"`
	PlatoonRate       bool `json:"platoon_rate"`
	ThreatLevel       bool `json:"threat_level"`
}

type UCColor struct {
	PlayerName PlayerNameColor `json:"player_name"`
}

type UCDigit struct {
	PR                uint `json:"pr"`
	Damage            uint `json:"damage"`
	MaxDamage         uint `json:"max_damage"`
	WinRate           uint `json:"win_rate"`
	KdRate            uint `json:"kd_rate"`
	Kill              uint `json:"kill"`
	PlanesKilled      uint `json:"planes_killed"`
	Exp               uint `json:"exp"`
	Battles           uint `json:"battles"`
	SurvivedRate      uint `json:"survived_rate"`
	HitRate           uint `json:"hit_rate"`
	AvgTier           uint `json:"avg_tier"`
	UsingShipTypeRate uint `json:"using_ship_type_rate"`
	UsingTierRate     uint `json:"using_tier_rate"`
	PlatoonRate       uint `json:"platoon_rate"`
	ThreatLevel       uint `json:"threat_level"`
}
