package model

//nolint:gochecknoglobals
var DefaultUserConfig = UserConfig{
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
		Skill: UCSkillColor{
			Text: UCSkillColorCode{
				Bad:         "#ff382d",
				BelowAvg:    "#fd9234",
				Avg:         "#ffd351",
				Good:        "#57e500",
				VeryGood:    "#44b200",
				Great:       "#02f7da",
				Unicum:      "#da6ff5",
				SuperUnicum: "#bf15ee",
			},
		},
		Tier: UCTierColor{
			Own: UCTierColorCode{
				Low:    "#8CA113",
				Middle: "#205B85",
				High:   "#990F4F",
			},
			Other: UCTierColorCode{
				Low:    "#E6F5B0",
				Middle: "#B3D7DD",
				High:   "#E3ADD5",
			},
		},
		ShipType: UCShipTypeColor{
			Own: UCShipTypeColorCode{
				CV: "#5E2883",
				BB: "#CA1028",
				CL: "#27853F",
				DD: "#D9760F",
				SS: "#233B8B",
			},
			Other: UCShipTypeColorCode{
				CV: "#CAB2D6",
				BB: "#FBB4C4",
				CL: "#CCEBC5",
				DD: "#FEE6AA",
				SS: "#B3CDE3",
			},
		},
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
	},
	TeamSummary: UCTeamSummary{
		MinShipBattles:    1,
		MinOverallBattles: 10,
	},
}

type UserConfig struct {
	// required
	InstallPath string `json:"install_path"`
	Appid       string `json:"appid"`
	// display
	FontSize string    `json:"font_size"`
	Display  UCDisplay `json:"display"`
	Color    UCColor   `json:"color"`
	Digit    UCDigit   `json:"digit"`
	// team summary
	TeamSummary UCTeamSummary `json:"team_summary"`
	// other
	SaveScreenshot    bool         `json:"save_screenshot"`
	SaveTempArenaInfo bool         `json:"save_temp_arena_info"`
	SendReport        bool         `json:"send_report"`
	NotifyUpdatable   bool         `json:"notify_updatable"`
	StatsPattern      StatsPattern `json:"stats_pattern"`
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
}

type UCColor struct {
	Skill      UCSkillColor    `json:"skill"`
	Tier       UCTierColor     `json:"tier"`
	ShipType   UCShipTypeColor `json:"ship_type"`
	PlayerName PlayerNameColor `json:"player_name"`
}

type UCTierColor struct {
	Own   UCTierColorCode `json:"own"`
	Other UCTierColorCode `json:"other"`
}

type UCSkillColor struct {
	Text UCSkillColorCode `json:"text"`
}

type UCShipTypeColor struct {
	Own   UCShipTypeColorCode `json:"own"`
	Other UCShipTypeColorCode `json:"other"`
}

type UCSkillColorCode struct {
	Bad         string `json:"bad"`
	BelowAvg    string `json:"below_avg"`
	Avg         string `json:"avg"`
	Good        string `json:"good"`
	VeryGood    string `json:"very_good"`
	Great       string `json:"great"`
	Unicum      string `json:"unicum"`
	SuperUnicum string `json:"super_unicum"`
}

type UCTierColorCode struct {
	Low    string `json:"low"`    // tier 1~4
	Middle string `json:"middle"` // tier 5~7
	High   string `json:"high"`   // tier 8~★
}

type UCShipTypeColorCode struct {
	SS string `json:"ss"`
	DD string `json:"dd"`
	CL string `json:"cl"`
	BB string `json:"bb"`
	CV string `json:"cv"`
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
}

type UCTeamSummary struct {
	MinShipBattles    uint `json:"min_ship_battles"`
	MinOverallBattles uint `json:"min_overall_battles"`
}
