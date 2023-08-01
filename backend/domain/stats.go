package domain

import (
	"math"
)

type Stats struct {
	useShipID     int
	accountInfo   WGAccountInfoData
	shipsStats    []WGShipsStatsData
	expectedStats map[int]NSExpectedStatsData
	useShipStats  WGShipsStatsData
	nickname      string
}

func NewStats(
	useShipID int,
	accountInfo WGAccountInfoData,
	shipsStats []WGShipsStatsData,
	expectedStats map[int]NSExpectedStatsData,
	nickname string,
) *Stats {
	var useShipStats WGShipsStatsData
	for _, v := range shipsStats {
		if v.ShipID == useShipID {
			useShipStats = v
			break
		}
	}

	return &Stats{
		useShipID:     useShipID,
		accountInfo:   accountInfo,
		shipsStats:    shipsStats,
		expectedStats: expectedStats,
		useShipStats:  useShipStats,
		nickname:      nickname,
	}
}

func (s *Stats) PR(category StatsCategory, pattern StatsPattern) float64 {
	switch category {
	case StatsCategoryShip:
		values := s.statsValues(StatsCategoryShip, pattern)
		battles := values.Battles

		actual := PRFactor{
			damage: avgDamage(values.DamageDealt, battles),
			frags:  avgKill(values.Frags, battles),
			wins:   winRate(values.Wins, battles),
		}
		expected := PRFactor{
			damage: s.expectedStats[s.useShipID].AverageDamageDealt,
			frags:  s.expectedStats[s.useShipID].AverageFrags,
			wins:   s.expectedStats[s.useShipID].WinRate,
		}

		return pr(actual, expected, battles)

	case StatsCategoryOverall:
		actual := PRFactor{}
		expected := PRFactor{}
		var allBattles uint

		for _, ship := range s.shipsStats {
			values := s.statsValuesForm(ship, pattern)
			battles := values.Battles

			shipID := ship.ShipID
			es, ok := s.expectedStats[shipID]
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

		return pr(actual, expected, allBattles)
	}

	return -1
}

func (s *Stats) Battles(category StatsCategory, pattern StatsPattern) uint {
	values := s.statsValues(category, pattern)
	return values.Battles
}

func (s *Stats) AvgDamage(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	damageDealt := values.DamageDealt
	battles := values.Battles
	return avgDamage(damageDealt, battles)
}

func (s *Stats) KdRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	frags := values.Frags

	death := values.Battles - values.SurvivedBattles
	if death < 1 {
		return float64(frags)
	}

	return float64(frags) / float64(death)
}

func (s *Stats) AvgKill(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return avgKill(values.Frags, values.Battles)
}

func (s *Stats) AvgExp(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	xp := values.Xp
	battles := values.Battles

	if battles < 1 {
		return 0
	}

	return float64(xp) / float64(battles)
}

func (s *Stats) WinRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	return winRate(values.Wins, values.Battles)
}

func (s *Stats) WinSurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	wins := values.Wins
	if wins < 1 {
		return 0
	}

	return float64(values.SurvivedWins) / float64(wins) * 100
}

func (s *Stats) LoseSurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	loses := values.Battles - values.Wins
	if loses < 1 {
		return 0
	}

	return float64(values.SurvivedBattles-values.SurvivedWins) / float64(loses) * 100
}

func (s *Stats) MainBatteryHitRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	shots := values.MainBattery.Shots
	if shots < 1 {
		return 0
	}

	return float64(values.MainBattery.Hits) / float64(shots) * 100
}

func (s *Stats) TorpedoesHitRate(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	shots := values.Torpedoes.Shots
	if shots < 1 {
		return 0
	}

	return float64(values.Torpedoes.Hits) / float64(shots) * 100
}

func (s *Stats) PlanesKilled(category StatsCategory, pattern StatsPattern) float64 {
	values := s.statsValues(category, pattern)
	battles := values.Battles
	if battles < 1 {
		return 0
	}

	return float64(values.PlanesKilled) / float64(battles)
}

func (s *Stats) AvgTier(
	pattern StatsPattern,
	warships map[int]Warship,
) float64 {
	var sum uint
	var battles uint

	for _, stats := range s.shipsStats {
		warship, ok := warships[stats.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(stats, pattern)
		sum += values.Battles * warship.Tier
		battles += values.Battles
	}

	if battles < 1 {
		return 0
	}

	return float64(sum) / float64(battles)
}

func (s *Stats) UsingTierRate(
	pattern StatsPattern,
	warships map[int]Warship,
) TierGroup {
	var tierGroup TierGroup

	for _, ship := range s.shipsStats {
		values := s.statsValuesForm(ship, pattern)
		battles := values.Battles
		if battles < 1 {
			continue
		}

		warship, ok := warships[ship.ShipID]
		if !ok {
			continue
		}
		tier := warship.Tier

		switch {
		case tier >= 1 && tier <= 4:
			tierGroup.Low += float64(battles)
		case tier >= 5 && tier <= 7:
			tierGroup.Middle += float64(battles)
		case tier >= 8:
			tierGroup.High += float64(battles)
		}
	}

	allBattles := tierGroup.Low + tierGroup.Middle + tierGroup.High
	if allBattles < 1 {
		return TierGroup{}
	}

	return TierGroup{
		Low:    tierGroup.Low / allBattles * 100,
		Middle: tierGroup.Middle / allBattles * 100,
		High:   tierGroup.High / allBattles * 100,
	}
}

func (s *Stats) UsingShipTypeRate(
	pattern StatsPattern,
	warships map[int]Warship,
) ShipTypeGroup {
	shipTypeMap := make(map[ShipType]float64)

	for _, ship := range s.shipsStats {
		values := s.statsValuesForm(ship, pattern)
		battles := values.Battles
		if battles < 1 {
			continue
		}

		warship, ok := warships[ship.ShipID]
		if !ok {
			continue
		}

		shipTypeMap[warship.Type] += float64(battles)
	}

	var allBattles float64
	for _, v := range shipTypeMap {
		allBattles += v
	}
	if allBattles < 1 {
		return ShipTypeGroup{}
	}

	return ShipTypeGroup{
		SS: shipTypeMap[SS] / allBattles * 100,
		DD: shipTypeMap[DD] / allBattles * 100,
		CL: shipTypeMap[CL] / allBattles * 100,
		BB: shipTypeMap[BB] / allBattles * 100,
		CV: shipTypeMap[CV] / allBattles * 100,
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

	return WGStatsValues{}
}

func (s *Stats) statsValuesForm(statsData WGShipsStatsData, pattern StatsPattern) WGStatsValues {
	switch pattern {
	case StatsPatternPvPAll:
		return statsData.Pvp
	case StatsPatternPvPSolo:
		return statsData.PvpSolo
	}

	return WGStatsValues{}
}

func avgDamage(damageDealt uint, battles uint) float64 {
	if battles < 1 {
		return -1
	}

	return float64(damageDealt) / float64(battles)
}

func avgKill(frags uint, battles uint) float64 {
	if battles < 1 {
		return -1
	}

	return float64(frags) / float64(battles)
}

func winRate(wins uint, battles uint) float64 {
	if battles < 1 {
		return -1
	}

	return float64(wins) / float64(battles) * 100
}

func pr(
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
