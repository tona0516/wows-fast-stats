package core

import "math"

// ref: https://asia.wows-numbers.com/personal/rating

func calculatePR(
	avgDamage float64,
	avgFrags float64,
	winRate float64,
	expectedDamage float64,
	expectedFrags float64,
	expectedWinRate float64,
) float64 {
	rDamage := safeDivide(avgDamage, expectedDamage)
	rFrags := safeDivide(avgFrags, expectedFrags)
	rWins := safeDivide(winRate, expectedWinRate)

	nDamage := math.Max(0, (rDamage-0.4)/(1-0.4))
	nFrags := math.Max(0, (rFrags-0.1)/(1-0.1))
	nWins := math.Max(0, (rWins-0.7)/(1-0.7))

	return 700*nDamage + 300*nFrags + 150*nWins
}

func calculateOverallPR(
	pattern StatsPattern,
	playerShipStats PlayerShipStats,
	warships Warships,
) float64 {
	var (
		battles         uint
		avgDamage       uint
		avgFrags        uint
		winRate         uint
		expectedDamage  float64
		expectedFrags   float64
		expectedWinRate float64
	)

	for shipID, ship := range playerShipStats {
		values := ship.shipStatsValues(pattern)

		if values.Battles == 0 {
			continue
		}

		warship, ok := warships[shipID]
		if !ok {
			continue
		}

		if warship.ServerAverage == nil {
			continue
		}
		serverAvg := warship.ServerAverage

		battles += values.Battles
		avgDamage += values.DamageDealt
		avgFrags += values.Frags
		winRate += values.Wins

		expectedDamage += serverAvg.Damage * float64(values.Battles)
		expectedFrags += serverAvg.Frags * float64(values.Battles)
		expectedWinRate += serverAvg.WinRate / 100 * float64(values.Battles)
	}

	return calculatePR(
		float64(avgDamage),
		float64(avgFrags),
		float64(winRate),
		expectedDamage,
		expectedFrags,
		expectedWinRate,
	)
}
