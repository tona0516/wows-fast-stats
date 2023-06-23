package vo

type UserConfig struct {
	InstallPath       string       `json:"install_path"`
	Appid             string       `json:"appid"`
	FontSize          string       `json:"font_size"`
	SaveScreenshot    bool         `json:"save_screenshot"`
	SaveTempArenaInfo bool         `json:"save_temp_arena_info"`
	SendReport        bool         `json:"send_report"`
	StatsPattern      StatsPattern `json:"stats_pattern"`
	Displays          Displays     `json:"displays"`
	CustomColor       CustomColor  `json:"custom_color"`
	CustomDigit       CustomDigit  `json:"custom_digit"`
}

type Displays struct {
	Basic   Basic   `json:"basic"`
	Ship    Ship    `json:"ship"`
	Overall Overall `json:"overall"`
}

type Basic struct {
	IsInAvg    bool `json:"is_in_avg"`
	PlayerName bool `json:"player_name"`
	ShipInfo   bool `json:"ship_info"`
}

type Ship struct {
	PR           bool `json:"pr"`
	Damage       bool `json:"damage"`
	WinRate      bool `json:"win_rate"`
	KdRate       bool `json:"kd_rate"`
	Kill         bool `json:"kill"`
	PlanesKilled bool `json:"planes_killed"`
	Exp          bool `json:"exp"`
	Battles      bool `json:"battles"`
	SurvivedRate bool `json:"survived_rate"`
	HitRate      bool `json:"hit_rate"`
}

type Overall struct {
	Damage            bool `json:"damage"`
	WinRate           bool `json:"win_rate"`
	KdRate            bool `json:"kd_rate"`
	Kill              bool `json:"kill"`
	Exp               bool `json:"exp"`
	Battles           bool `json:"battles"`
	SurvivedRate      bool `json:"survived_rate"`
	AvgTier           bool `json:"avg_tier"`
	UsingShipTypeRate bool `json:"using_ship_type_rate"`
	UsingTierRate     bool `json:"using_tier_rate"`
}

type CustomColor struct {
	Skill    SkillColor    `json:"skill"`
	Tier     TierColor     `json:"tier"`
	ShipType ShipTypeColor `json:"ship_type"`
}

type TierColor struct {
	Own   TierColorCode `json:"own"`
	Other TierColorCode `json:"other"`
}

type SkillColor struct {
	Text       SkillColorCode `json:"text"`
	Background SkillColorCode `json:"background"`
}

type ShipTypeColor struct {
	Own   ShipTypeColorCode `json:"own"`
	Other ShipTypeColorCode `json:"other"`
}

type SkillColorCode struct {
	Bad         string `json:"bad"`
	BelowAvg    string `json:"below_avg"`
	Avg         string `json:"avg"`
	Good        string `json:"good"`
	VeryGood    string `json:"very_good"`
	Great       string `json:"great"`
	Unicum      string `json:"unicum"`
	SuperUnicum string `json:"super_unicum"`
}

type TierColorCode struct {
	Low    string `json:"low"`    // tier 1~4
	Middle string `json:"middle"` // tier 5~7
	High   string `json:"high"`   // tier 8~★
}

type ShipTypeColorCode struct {
	SS string `json:"ss"`
	DD string `json:"dd"`
	CL string `json:"cl"`
	BB string `json:"bb"`
	CV string `json:"cv"`
}

type CustomDigit struct {
	PR                uint `json:"pr"`
	Damage            uint `json:"damage"`
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
}
