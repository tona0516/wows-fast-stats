package persistence

type UserSetting struct {
	Version  uint `json:"version"`
	Required struct {
		InstallPath string `json:"install_path"`
	} `json:"required"`
	Optional struct {
		Stats struct {
			PR struct {
				Ship             bool   `json:"ship"`
				Overall          bool   `json:"overall"`
				Digit            uint   `json:"digit"`
				OwnNameColorType string `json:"own_name_color_type"`
			} `json:"pr"`
			Damage struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"damage"`
			MaxDamage struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"max_damage"`
			WinRate struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"win_rate"`
			KdRate struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"kd_rate"`
			Kill struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"kill"`
			PlanesKilled struct {
				Ship  bool `json:"ship"`
				Digit uint `json:"digit"`
			} `json:"planes_killed"`
			Exp struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"exp"`
			Battles struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"battles"`
			SurvivedRate struct {
				Ship    bool `json:"ship"`
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"survived_rate"`
			HitRate struct {
				Ship  bool `json:"ship"`
				Digit uint `json:"digit"`
			} `json:"hit_rate"`
			AvgTier struct {
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"avg_tier"`
			UsingShipTypeRate struct {
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"using_ship_type_rate"`
			UsingTierRate struct {
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"using_tier_rate"`
			PlatoonRate struct {
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"platoon_rate"`
			ThreatLevel struct {
				Overall bool `json:"overall"`
				Digit   uint `json:"digit"`
			} `json:"threat_level"`
		} `json:"stats"`
		Skill struct {
			Bad struct {
				TextColor string `json:"text_color"`
			} `json:"bad"`
			BelowAvg struct {
				TextColor string `json:"text_color"`
			} `json:"below_avg"`
			Avg struct {
				TextColor string `json:"text_color"`
			} `json:"avg"`
			Good struct {
				TextColor string `json:"text_color"`
			} `json:"good"`
			VeryGood struct {
				TextColor string `json:"text_color"`
			} `json:"very_good"`
			Great struct {
				TextColor string `json:"text_color"`
			} `json:"great"`
			Unicum struct {
				TextColor string `json:"text_color"`
			} `json:"unicum"`
			SuperUnicum struct {
				TextColor string `json:"text_color"`
			} `json:"super_unicum"`
		} `json:"skill"`
		TierGroup struct {
			Low struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"low"` // tier 1~4
			Middle struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"middle"` // tier 5~7
			High struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"high"` // tier 8~★
		}
		ShipType struct {
			SS struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"ss"`
			DD struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"dd"`
			CL struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"cl"`
			BB struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"bb"`
			CV struct {
				OwnColor   string `json:"own_color"`
				OtherColor string `json:"other_color"`
			} `json:"cv"`
		} `json:"optional"`
		TeamSummary struct {
			MinShipBattles    uint `json:"min_ship_battles"`
			MinOverallBattles uint `json:"min_overall_battles"`
		} `json:"team_summary"`
		FontSize          string `json:"font_size"`
		StatsPattern      string `json:"stats_pattern"`
		ShowLanguageFrag  bool   `json:"show_language_frag"`
		SaveScreenshot    bool   `json:"save_screenshot"`
		SaveTempArenaInfo bool   `json:"save_temp_arena_info"`
		SendReport        bool   `json:"send_report"`
		NotifyUpdatable   bool   `json:"notify_updatable"`
	}
}
