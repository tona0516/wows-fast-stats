package data

type PlayerInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Clan     Clan   `json:"clan"`
	IsHidden bool   `json:"is_hidden"`
}

type PlayerStats struct {
	ShipStats    ShipStats    `json:"ship"`
	OverallStats OverallStats `json:"overall"`
}

func NewPlayerStats(
	statsPattern StatsPattern,
	stats *PersonalStats,
	accountID int,
	shipID int,
	tempArenaInfo TempArenaInfo,
	warships Warships,
) PlayerStats {
	threatLevel := CalculateThreatLevel(NewThreatLevelCalculatorFactor(
		accountID,
		tempArenaInfo,
		warships,
		shipID,
		stats.Battles(StatsCategoryShip, statsPattern),
		stats.AvgDamage(StatsCategoryShip, statsPattern).Value,
		stats.WinRate(StatsCategoryShip, statsPattern).Value,
		stats.SurvivedRate(StatsCategoryShip, statsPattern).All,
		stats.PlanesKilled(StatsCategoryShip),
		stats.Battles(StatsCategoryOverall, statsPattern),
		stats.AvgDamage(StatsCategoryOverall, statsPattern).Value,
		stats.WinRate(StatsCategoryOverall, statsPattern).Value,
		stats.AvgKill(StatsCategoryOverall, statsPattern),
		stats.KdRate(StatsCategoryOverall, statsPattern),
	))

	return PlayerStats{
		ShipStats: ShipStats{
			Battles:         stats.Battles(StatsCategoryShip, statsPattern),
			Damage:          stats.AvgDamage(StatsCategoryShip, statsPattern),
			MaxDamage:       stats.MaxDamage(StatsCategoryShip, statsPattern),
			WinRate:         stats.WinRate(StatsCategoryShip, statsPattern),
			SurvivedRate:    stats.SurvivedRate(StatsCategoryShip, statsPattern),
			KdRate:          stats.KdRate(StatsCategoryShip, statsPattern),
			Kill:            stats.AvgKill(StatsCategoryShip, statsPattern),
			Exp:             stats.AvgExp(StatsCategoryShip, statsPattern),
			PR:              stats.PR(StatsCategoryShip, statsPattern),
			HitRate:         stats.HitRate(statsPattern),
			PlanesKilled:    stats.PlanesKilled(statsPattern),
			PlatoonRate:     stats.PlatoonRate(StatsCategoryShip),
			EfficiencyBadge: stats.EfficiencyBadge(),
		},
		OverallStats: OverallStats{
			Battles:           stats.Battles(StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(StatsCategoryOverall, statsPattern),
			SurvivedRate:      stats.SurvivedRate(StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(StatsCategoryOverall, statsPattern),
			PR:                stats.PR(StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(StatsCategoryOverall),
			EfficiencyBadge:   stats.EfficiencyBadges(),
			ThreatLevel:       threatLevel,
		},
	}
}

type ShipStats struct {
	Battles         uint            `json:"battles"`
	Damage          RatingValue     `json:"damage"`
	MaxDamage       MaxDamage       `json:"max_damage"`
	WinRate         RatingValue     `json:"win_rate"`
	SurvivedRate    SurvivedRate    `json:"survived_rate"`
	KdRate          float64         `json:"kd_rate"`
	Kill            float64         `json:"kill"`
	Exp             float64         `json:"exp"`
	PR              RatingValue     `json:"pr"`
	HitRate         HitRate         `json:"hit_rate"`
	PlanesKilled    float64         `json:"planes_killed"`
	PlatoonRate     float64         `json:"platoon_rate"`
	EfficiencyBadge EfficiencyBadge `json:"efficiency_badge"`
}

type OverallStats struct {
	Battles           uint                 `json:"battles"`
	Damage            RatingValue          `json:"damage"`
	MaxDamage         MaxDamage            `json:"max_damage"`
	WinRate           RatingValue          `json:"win_rate"`
	SurvivedRate      SurvivedRate         `json:"survived_rate"`
	KdRate            float64              `json:"kd_rate"`
	Kill              float64              `json:"kill"`
	Exp               float64              `json:"exp"`
	PR                RatingValue          `json:"pr"`
	ThreatLevel       ThreatLevel          `json:"threat_level"`
	AvgTier           float64              `json:"avg_tier"`
	UsingShipTypeRate ShipTypeGroup        `json:"using_ship_type_rate"`
	UsingTierRate     TierGroup            `json:"using_tier_rate"`
	PlatoonRate       float64              `json:"platoon_rate"`
	EfficiencyBadge   EfficiencyBadgeGroup `json:"efficiency_badge"`
}

type Player struct {
	PlayerInfo PlayerInfo  `json:"player_info"`
	Warship    Warship     `json:"warship"`
	PvPSolo    PlayerStats `json:"pvp_solo"`
	PvPAll     PlayerStats `json:"pvp_all"`
	RankSolo   PlayerStats `json:"rank_solo"`
}
