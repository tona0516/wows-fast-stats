package core

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

func NewOverallStats(
	statsPattern StatsPattern,
	accountInfoData WGAccountInfoData,
	playerShipStats PlayerShipStats,
	playerShipBadges PlayerShipBadges,
	shipID ShipID,
	accountID AccountID,
	vehicles []Vehicle,
	warships Warships,
) OverallStats {
	values := accountInfoData.playerStatsValue(statsPattern)
	if values.Battles == 0 {
		return OverallStats{}
	}

	maxDamage := MaxDamage{
		ShipID: values.MaxDamageDealtShipID,
		Value:  values.MaxDamageDealt,
	}
	if warship, ok := warships[values.MaxDamageDealtShipID]; ok {
		maxDamage.ShipName = warship.Name
		maxDamage.ShipTier = warship.Tier
	}

	winRateRating := NewRatingValue(
		values.winRate(),
		NewRatingFromWinRate(values.winRate()),
	)

	overallPR := calculateOverallPR(
		statsPattern,
		playerShipStats,
		warships,
	)
	pr := NewRatingValue(
		overallPR,
		NewRatingFromPR(overallPR),
	)

	var threatLevel ThreatLevel
	shipStats, ok := playerShipStats[shipID]
	if ok {
		shipValues := shipStats.shipStatsValues(statsPattern)
		tlc := NewThreatLevelCalculator()

		tti := ThreatLevelInput{
			Vehicles:         vehicles,
			Warships:         warships,
			ShipID:           shipID,
			ShipBattles:      shipValues.Battles,
			ShipDamage:       shipValues.avgDamage(),
			ShipWinRate:      shipValues.winRate(),
			ShipSurvivedRate: shipValues.survivedRate().All,
			ShipPlanesKilled: shipValues.avgPlanesKilled(),
			OverallBattles:   values.Battles,
			OverallDamage:    values.avgDamage(),
			OverallWinRate:   values.winRate(),
			OverallKill:      values.avgKills(),
			OverallKdRate:    values.kdRate(),
		}
		threatLevel = tlc.Calculate(tti)
	}

	return OverallStats{
		Battles: values.Battles,
		Damage: NewRatingValue(
			values.avgDamage(),
			RatingNone,
		),
		MaxDamage:         maxDamage,
		WinRate:           winRateRating,
		SurvivedRate:      values.survivedRate(),
		KdRate:            values.kdRate(),
		Kill:              values.avgKills(),
		Exp:               values.exp(),
		PR:                pr,
		ThreatLevel:       threatLevel,
		AvgTier:           playerShipStats.avgTier(statsPattern, warships),
		UsingShipTypeRate: playerShipStats.usingShipTypeRate(statsPattern, warships),
		UsingTierRate:     playerShipStats.usingTierRate(statsPattern, warships),
		PlatoonRate:       accountInfoData.platoonRate(),
		EfficiencyBadge:   playerShipBadges.efficiencyBadges(),
	}
}
