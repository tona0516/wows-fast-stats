package core

type OverallStats struct {
	Battles           uint                 `json:"battles"`
	Damage            RatingValue          `json:"damage"`
	MaxDamage         MaxDamage            `json:"maxDamage"`
	WinRate           RatingValue          `json:"winRate"`
	SurvivedRate      SurvivedRate         `json:"survivedRate"`
	KdRate            float64              `json:"kdRate"`
	Kill              float64              `json:"kill"`
	Exp               float64              `json:"exp"`
	PR                RatingValue          `json:"pr"`
	AvgTier           float64              `json:"avgTier"`
	UsingShipTypeRate ShipTypeGroup        `json:"usingShipTypeRate"`
	UsingTierRate     TierGroup            `json:"usingTierRate"`
	PlatoonRate       float64              `json:"platoonRate"`
	EfficiencyBadge   EfficiencyBadgeGroup `json:"efficiencyBadge"`
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
		AvgTier:           playerShipStats.avgTier(statsPattern, warships),
		UsingShipTypeRate: playerShipStats.usingShipTypeRate(statsPattern, warships),
		UsingTierRate:     playerShipStats.usingTierRate(statsPattern, warships),
		PlatoonRate:       accountInfoData.platoonRate(),
		EfficiencyBadge:   playerShipBadges.efficiencyBadges(),
	}
}
