package core

type TeamAverageStats struct {
	ShipPR         float64 `json:"ship_pr"`
	ShipDamage     float64 `json:"ship_damage"`
	ShipWinRate    float64 `json:"ship_win_rate"`
	ShipBattles    uint    `json:"ship_battles"`
	OverallPR      float64 `json:"overall_pr"`
	OverallDamage  float64 `json:"overall_damage"`
	OverallWinRate float64 `json:"overall_win_rate"`
	OverallBattles uint    `json:"overall_battles"`
}

func NewTeamAverageStats(players Players, statsPattern StatsPattern) TeamAverageStats {
	var shipPRSum, shipDamageSum, shipWinRateSum float64
	var shipBattlesSum, shipStatsCount uint

	var overallPRSum, overallDamageSum, overallWinRateSum float64
	var overallBattlesSum, overallStatsCount uint

	for _, player := range players {
		var shipStats ShipStats
		var overallStats OverallStats
		switch statsPattern {
		case StatsPatternPvPSolo:
			shipStats = player.PvPSolo.ShipStats
			overallStats = player.PvPSolo.OverallStats
		case StatsPatternPvPAll:
			shipStats = player.PvPAll.ShipStats
			overallStats = player.PvPAll.OverallStats
		case StatsPatternRankSolo:
			shipStats = player.RankSolo.ShipStats
			overallStats = player.RankSolo.OverallStats
		}

		if shipStats.Battles > 0 {
			shipPRSum += shipStats.PR.Value
			shipDamageSum += shipStats.Damage.Value
			shipWinRateSum += shipStats.WinRate.Value
			shipBattlesSum += shipStats.Battles

			shipStatsCount++
		}

		if overallStats.Battles > 0 {
			overallPRSum += overallStats.PR.Value
			overallDamageSum += overallStats.Damage.Value
			overallWinRateSum += overallStats.WinRate.Value
			overallBattlesSum += overallStats.Battles

			overallStatsCount++
		}
	}

	return TeamAverageStats{
		ShipPR:         safeDivide(shipPRSum, shipStatsCount),
		ShipDamage:     safeDivide(shipDamageSum, shipStatsCount),
		ShipWinRate:    safeDivide(shipWinRateSum, shipStatsCount),
		ShipBattles:    uint(safeDivide(shipBattlesSum, shipStatsCount)),
		OverallPR:      safeDivide(overallPRSum, overallStatsCount),
		OverallDamage:  safeDivide(overallDamageSum, overallStatsCount),
		OverallWinRate: safeDivide(overallWinRateSum, overallStatsCount),
		OverallBattles: uint(safeDivide(overallBattlesSum, overallStatsCount)),
	}
}
