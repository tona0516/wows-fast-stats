package data

import (
	"wfs/backend/util"
)

// BuildPlayerStats は単一プレイヤーの統計を構築する。
func BuildPlayerStats(
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

// CalculateTeamAverageStats はチームの平均統計を計算する。
func CalculateTeamAverageStats(
	players Players,
	statsPattern StatsPattern,
) TeamAverageStats {
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
		ShipPR:         util.SafeDivide(shipPRSum, shipStatsCount),
		ShipDamage:     util.SafeDivide(shipDamageSum, shipStatsCount),
		ShipWinRate:    util.SafeDivide(shipWinRateSum, shipStatsCount),
		ShipBattles:    uint(util.SafeDivide(float64(shipBattlesSum), shipStatsCount)),
		OverallPR:      util.SafeDivide(overallPRSum, overallStatsCount),
		OverallDamage:  util.SafeDivide(overallDamageSum, overallStatsCount),
		OverallWinRate: util.SafeDivide(overallWinRateSum, overallStatsCount),
		OverallBattles: uint(util.SafeDivide(float64(overallBattlesSum), overallStatsCount)),
	}
}
