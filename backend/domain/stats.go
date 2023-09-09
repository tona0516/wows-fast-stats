package domain

import (
	"math"
	"wfs/backend/apperr"
	"wfs/backend/logger"

	"github.com/morikuni/failure"
)

type Stats struct {
	useShipID        int
	accountInfo      WGAccountInfoData
	useShipStats     WGShipsStatsData
	allShipsStats    []WGShipsStatsData
	allExpectedStats AllExpectedStats
	warships         Warships
}

func NewStats(
	useShipID int,
	accountInfo WGAccountInfoData,
	allShipsStats []WGShipsStatsData,
	expectedStats AllExpectedStats,
	warships Warships,
) *Stats {
	var useShipStats WGShipsStatsData
	for _, v := range allShipsStats {
		if v.ShipID == useShipID {
			useShipStats = v
			break
		}
	}

	return &Stats{
		useShipID:        useShipID,
		accountInfo:      accountInfo,
		useShipStats:     useShipStats,
		allShipsStats:    allShipsStats,
		allExpectedStats: expectedStats,
		warships:         warships,
	}
}

func (s *Stats) PR(category StatsCategory, pattern StatsPattern) float64 {
	switch category {
	case StatsCategoryShip:
		values := s.statsValues(StatsCategoryShip, pattern)
		battles := values.Battles

		return s.pr(
			PRFactor{
				damage: avgDamage(values.DamageDealt, battles),
				frags:  avgKill(values.Frags, battles),
				wins:   winRate(values.Wins, battles),
			},
			PRFactor{
				damage: s.allExpectedStats[s.useShipID].AverageDamageDealt,
				frags:  s.allExpectedStats[s.useShipID].AverageFrags,
				wins:   s.allExpectedStats[s.useShipID].WinRate,
			},
			battles,
		)

	case StatsCategoryOverall:
		var (
			actual     PRFactor
			expected   PRFactor
			allBattles uint
		)

		for _, ship := range s.allShipsStats {
			values := s.statsValuesForm(ship, pattern)
			battles := values.Battles

			es, ok := s.allExpectedStats[ship.ShipID]
			if !ok {
				continue
			}

			actual.damage += float64(values.DamageDealt)
			actual.frags += float64(values.Frags)
			actual.wins += float64(values.Wins)

			expected.damage += es.AverageDamageDealt * float64(battles)
			expected.frags += es.AverageFrags * float64(battles)
			expected.wins += es.WinRate / 100 * float64(battles)

			allBattles += battles
		}

		return s.pr(actual, expected, allBattles)
	}

	logger.Error(failure.New(apperr.UnexpectedError))
	return -1
}

func (s *Stats) Battles(category StatsCategory, pattern StatsPattern) uint {
	return s.statsValues(category, pattern).Battles
}

func (s *Stats) AvgDamage(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return avgDamage(values.DamageDealt, values.Battles)
}

func (s *Stats) MaxDamage(category StatsCategory, pattern StatsPattern) uint {
	return s.statsValues(category, pattern).MaxDamageDealt
}

func (s *Stats) KdRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	death := values.Battles - values.SurvivedBattles
	if death < 1 {
		death = 1
	}

	return float64(values.Frags) / float64(death)
}

func (s *Stats) AvgKill(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return avgKill(values.Frags, values.Battles)
}

func (s *Stats) AvgExp(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return div(values.Xp, values.Battles)
}

func (s *Stats) WinRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return winRate(values.Wins, values.Battles)
}

func (s *Stats) WinSurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return percentage(values.SurvivedWins, values.Wins)
}

func (s *Stats) LoseSurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	loses := values.Battles - values.Wins
	return percentage(values.SurvivedBattles-values.SurvivedWins, loses)
}

func (s *Stats) MainBatteryHitRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return percentage(values.MainBattery.Hits, values.MainBattery.Shots)
}

func (s *Stats) TorpedoesHitRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return percentage(values.Torpedoes.Hits, values.Torpedoes.Shots)
}

func (s *Stats) PlanesKilled(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return div(values.PlanesKilled, values.Battles)
}

func (s *Stats) AvgTier(
	pattern StatsPattern,
) float64 {
	var (
		sum        uint
		allBattles uint
	)

	for _, stats := range s.allShipsStats {
		warship, ok := s.warships[stats.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(stats, pattern)
		sum += values.Battles * warship.Tier
		allBattles += values.Battles
	}

	return div(sum, allBattles)
}

func (s *Stats) UsingTierRate(
	pattern StatsPattern,
) TierGroup {
	tierGroupMap := make(map[string]uint)

	for _, ship := range s.allShipsStats {
		warship, ok := s.warships[ship.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(ship, pattern)
		tier := warship.Tier
		battles := values.Battles
		switch {
		case tier >= 1 && tier <= 4:
			tierGroupMap["low"] += battles
		case tier >= 5 && tier <= 7:
			tierGroupMap["middle"] += battles
		case tier >= 8:
			tierGroupMap["high"] += battles
		}
	}

	var allBattles uint
	for _, v := range tierGroupMap {
		allBattles += v
	}

	return TierGroup{
		Low:    percentage(tierGroupMap["low"], allBattles),
		Middle: percentage(tierGroupMap["middle"], allBattles),
		High:   percentage(tierGroupMap["high"], allBattles),
	}
}

func (s *Stats) UsingShipTypeRate(
	pattern StatsPattern,
) ShipTypeGroup {
	shipTypeMap := make(map[ShipType]uint)

	for _, ship := range s.allShipsStats {
		warship, ok := s.warships[ship.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(ship, pattern)
		shipTypeMap[warship.Type] += values.Battles
	}

	var allBattles uint
	for _, v := range shipTypeMap {
		allBattles += v
	}

	return ShipTypeGroup{
		SS: percentage(shipTypeMap[SS], allBattles),
		DD: percentage(shipTypeMap[DD], allBattles),
		CL: percentage(shipTypeMap[CL], allBattles),
		BB: percentage(shipTypeMap[BB], allBattles),
		CV: percentage(shipTypeMap[CV], allBattles),
	}
}

func (s *Stats) statsValues(category StatsCategory, pattern StatsPattern) WGStatsValues {
	switch category {
	case StatsCategoryShip:
		switch pattern {
		case StatsPatternPvPAll:
			return s.useShipStats.Pvp
		case StatsPatternPvPSolo:
			return s.useShipStats.PvpSolo
		}
	case StatsCategoryOverall:
		switch pattern {
		case StatsPatternPvPAll:
			return s.accountInfo.Statistics.Pvp
		case StatsPatternPvPSolo:
			return s.accountInfo.Statistics.PvpSolo
		}
	}

	logger.Error(failure.New(apperr.UnexpectedError))
	return WGStatsValues{}
}

func (s *Stats) statsValuesForm(statsData WGShipsStatsData, pattern StatsPattern) WGStatsValues {
	switch pattern {
	case StatsPatternPvPAll:
		return statsData.Pvp
	case StatsPatternPvPSolo:
		return statsData.PvpSolo
	}

	logger.Error(failure.New(apperr.UnexpectedError))
	return WGStatsValues{}
}

func (s *Stats) pr(
	actual PRFactor,
	expected PRFactor,
	battles uint,
) float64 {
	if battles < 1 {
		return -1
	}

	ratio := PRFactor{
		damage: actual.damage / expected.damage,
		frags:  actual.frags / expected.frags,
		wins:   actual.wins / expected.wins,
	}

	if !ratio.Valid() {
		return -1
	}

	norm := PRFactor{
		damage: math.Max(0, (ratio.damage-0.4)/(1-0.4)),
		frags:  math.Max(0, (ratio.frags-0.1)/(1-0.1)),
		wins:   math.Max(0, (ratio.wins-0.7)/(1-0.7)),
	}

	return 700*norm.damage + 300*norm.frags + 150*norm.wins
}

func avgDamage(damageDealt uint, battles uint) float64 {
	return div(damageDealt, battles)
}

func avgKill(frags uint, battles uint) float64 {
	return div(frags, battles)
}

func winRate(wins uint, battles uint) float64 {
	return percentage(wins, battles)
}

func div(a uint, b uint) float64 {
	if b <= 0 {
		return 0
	}

	return float64(a) / float64(b)
}

func percentage(a uint, b uint) float64 {
	return div(a, b) * 100
}
