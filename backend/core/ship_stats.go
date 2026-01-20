package core

type ShipStats struct {
	Battles         uint            `json:"battles"`
	Damage          RatingValue     `json:"damage"`
	MaxDamage       MaxDamage       `json:"maxDamage"`
	WinRate         RatingValue     `json:"winRate"`
	SurvivedRate    SurvivedRate    `json:"survivedRate"`
	KdRate          float64         `json:"kdRate"`
	Kill            float64         `json:"kill"`
	Exp             float64         `json:"exp"`
	PR              RatingValue     `json:"pr"`
	HitRate         HitRate         `json:"hitRate"`
	PlanesKilled    float64         `json:"planesKilled"`
	PlatoonRate     float64         `json:"platoonRate"`
	EfficiencyBadge EfficiencyBadge `json:"efficiencyBadge"`
}

func NewShipStats(
	statsPattern StatsPattern,
	playerShipStats PlayerShipStats,
	playerShipBadges PlayerShipBadges,
	shipID ShipID,
	warships Warships,
) ShipStats {
	shipStats, ok := playerShipStats[shipID]
	if !ok {
		return ShipStats{}
	}
	values := shipStats.shipStatsValues(statsPattern)

	battles := values.Battles
	if battles == 0 {
		return ShipStats{}
	}

	avgDamage := values.avgDamage()
	avgFrags := values.avgKills()
	winRate := values.winRate()
	winRateRating := NewRatingValue(
		winRate,
		NewRatingFromWinRate(winRate),
	)

	pr := NewRatingValueNone()
	damageRating := NewRatingValue(
		avgDamage,
		RatingNone,
	)

	if warship, ok := warships[shipID]; ok {
		if warship.ServerAverage != nil {
			sa := *warship.ServerAverage
			shipPR := calculatePR(
				avgDamage,
				avgFrags,
				winRate,
				sa.Damage,
				sa.Frags,
				sa.WinRate,
			)
			pr = NewRatingValue(
				shipPR,
				NewRatingFromPR(shipPR),
			)

			damageRating.Rating = NewRatingFromShipDamage(
				avgDamage,
				sa.Damage,
			)
		}
	}

	return ShipStats{
		Battles:         battles,
		Damage:          damageRating,
		MaxDamage:       MaxDamage{Value: values.MaxDamageDealt},
		WinRate:         winRateRating,
		SurvivedRate:    values.survivedRate(),
		KdRate:          values.kdRate(),
		Kill:            avgFrags,
		Exp:             values.exp(),
		PR:              pr,
		HitRate:         values.hitRate(),
		PlanesKilled:    values.avgPlanesKilled(),
		PlatoonRate:     shipStats.platoonRate(),
		EfficiencyBadge: playerShipBadges.efficiencyBadge(shipID),
	}
}
