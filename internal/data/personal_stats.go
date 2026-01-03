package data

import (
	"math"
)

type PersonalStats struct {
	useShipID        int
	accountInfo      WGAccountInfoData
	useShipStats     WGShipsStatsData
	allShipsStats    []WGShipsStatsData
	shipsBadges      []WGShipsBadgesData
	allExpectedStats ExpectedStats
	warships         Warships
	tempArenaInfo    TempArenaInfo
}

func NewPersonalStats(
	useShipID int,
	accountInfo WGAccountInfoData,
	allShipsStats []WGShipsStatsData,
	shipsBadges []WGShipsBadgesData,
	expectedStats ExpectedStats,
	warships Warships,
	tempArenaInfo TempArenaInfo,
) *PersonalStats {
	var useShipStats WGShipsStatsData
	for _, v := range allShipsStats {
		if v.ShipID == useShipID {
			useShipStats = v
			break
		}
	}

	return &PersonalStats{
		useShipID:        useShipID,
		accountInfo:      accountInfo,
		useShipStats:     useShipStats,
		allShipsStats:    allShipsStats,
		shipsBadges:      shipsBadges,
		allExpectedStats: expectedStats,
		warships:         warships,
		tempArenaInfo:    tempArenaInfo,
	}
}

func (s *PersonalStats) PR(category StatsCategory, pattern StatsPattern) RatingValue {
	switch category {
	case StatsCategoryShip:
		values, _ := s.statsValues(pattern)
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

	return NewRatingValueNone()
}

func (s *PersonalStats) Battles(category StatsCategory, pattern StatsPattern) uint {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return ship.Battles
	case StatsCategoryOverall:
		return player.Battles
	}

	return 0
}

func (s *PersonalStats) AvgDamage(category StatsCategory, pattern StatsPattern) RatingValue {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		value := avgDamage(ship.DamageDealt, ship.Battles)
		rating := NewRatingFromShipDamage(value, s.allExpectedStats[s.useShipID].AverageDamageDealt)
		return NewRatingValue(value, rating)
	case StatsCategoryOverall:
		value := avgDamage(player.DamageDealt, player.Battles)
		return NewRatingValue(value, RatingNone)
	}

	return NewRatingValueNone()
}

func (s *PersonalStats) MaxDamage(category StatsCategory, pattern StatsPattern) MaxDamage {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return MaxDamage{
			Value: ship.MaxDamageDealt,
		}
	case StatsCategoryOverall:
		shipID := player.MaxDamageDealtShipID
		warship := s.warships[shipID]
		return MaxDamage{
			ShipID:   shipID,
			ShipName: warship.Name,
			ShipTier: warship.Tier,
			Value:    player.MaxDamageDealt,
		}
	}

	return MaxDamage{}
}

func (s *PersonalStats) KdRate(category StatsCategory, pattern StatsPattern) float64 {
	var (
		survivedBattles uint
		frags           uint
		battles         uint
	)

	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		survivedBattles = ship.SurvivedBattles
		frags = ship.Frags
		battles = ship.Battles
	case StatsCategoryOverall:
		survivedBattles = player.SurvivedBattles
		frags = player.Frags
		battles = player.Battles
	}

	death := max(battles-survivedBattles, 1)

	return float64(frags) / float64(death)
}

func (s *PersonalStats) AvgKill(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return avgKill(ship.Frags, ship.Battles)
	case StatsCategoryOverall:
		return avgKill(player.Frags, player.Battles)
	}

	return 0
}

func (s *PersonalStats) AvgExp(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return div(ship.Xp, ship.Battles)
	case StatsCategoryOverall:
		return div(player.Xp, player.Battles)
	}

	return 0
}

func (s *PersonalStats) WinRate(category StatsCategory, pattern StatsPattern) RatingValue {
	ship, player := s.statsValues(pattern)

	var value float64
	switch category {
	case StatsCategoryShip:
		value = winRate(ship.Wins, ship.Battles)
	case StatsCategoryOverall:
		value = winRate(player.Wins, player.Battles)
	}

	rating := NewRatingFromWinRate(value)
	return NewRatingValue(value, rating)
}

func (s *PersonalStats) SurvivedRate(category StatsCategory, pattern StatsPattern) SurvivedRate {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return SurvivedRate{
			All:  percentage(ship.SurvivedBattles, ship.Battles),
			Win:  percentage(ship.SurvivedWins, ship.Wins),
			Lose: percentage(ship.SurvivedBattles-ship.SurvivedWins, ship.Battles-ship.Wins),
		}
	case StatsCategoryOverall:
		return SurvivedRate{
			All:  percentage(player.SurvivedBattles, player.Battles),
			Win:  percentage(player.SurvivedWins, player.Wins),
			Lose: percentage(player.SurvivedBattles-player.SurvivedWins, player.Battles-player.Wins),
		}
	}

	return SurvivedRate{}
}

func (s *PersonalStats) HitRate(pattern StatsPattern) HitRate {
	ship, _ := s.statsValues(pattern)
	return HitRate{
		MainBattery: percentage(ship.MainBattery.Hits, ship.MainBattery.Shots),
		Torpedoes:   percentage(ship.Torpedoes.Hits, ship.Torpedoes.Shots),
	}
}

func (s *PersonalStats) PlanesKilled(pattern StatsPattern) float64 {
	ship, _ := s.statsValues(pattern)
	return div(ship.PlanesKilled, ship.Battles)
}

func (s *PersonalStats) AvgTier(
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

func (s *PersonalStats) UsingTierRate(
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

func (s *PersonalStats) UsingShipTypeRate(
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
		SS: percentage(shipTypeMap[ShipTypeSS], allBattles),
		DD: percentage(shipTypeMap[ShipTypeDD], allBattles),
		CL: percentage(shipTypeMap[ShipTypeCL], allBattles),
		BB: percentage(shipTypeMap[ShipTypeBB], allBattles),
		CV: percentage(shipTypeMap[ShipTypeCV], allBattles),
	}
}

func (s *PersonalStats) PlatoonRate(
	category StatsCategory,
) float64 {
	switch category {
	case StatsCategoryShip:
		stats := s.useShipStats
		return platoonRate(
			stats.Pvp.Battles,
			stats.PvpSolo.Battles,
			stats.PvpDiv2.Battles,
			stats.PvpDiv3.Battles,
		)
	case StatsCategoryOverall:
		stats := s.accountInfo.Statistics
		return platoonRate(
			stats.Pvp.Battles,
			stats.PvpSolo.Battles,
			stats.PvpDiv2.Battles,
			stats.PvpDiv3.Battles,
		)
	}

	return 0
}

func (s *PersonalStats) EfficiencyBadge() EfficiencyBadge {
	for _, b := range s.shipsBadges {
		if b.ShipID == s.useShipID {
			return s.toEfficiencyBadge(b.TopGradeClass)
		}
	}

	return EfficiencyBadgeNone
}

func (s *PersonalStats) EfficiencyBadges() EfficiencyBadgeGroup {
	var badges EfficiencyBadgeGroup

	for _, b := range s.shipsBadges {
		switch s.toEfficiencyBadge(b.TopGradeClass) {
		case EfficiencyBadgeExpert:
			badges.Expert++
		case EfficiencyBadgeFirst:
			badges.First++
		case EfficiencyBadgeSecond:
			badges.Second++
		case EfficiencyBadgeThird:
			badges.Third++
		}
	}

	return badges
}

func (s *PersonalStats) toEfficiencyBadge(value int) EfficiencyBadge {
	switch value {
	case 1:
		return EfficiencyBadgeExpert
	case 2:
		return EfficiencyBadgeFirst
	case 3:
		return EfficiencyBadgeSecond
	case 4:
		return EfficiencyBadgeThird
	default:
		return EfficiencyBadgeNone
	}
}

func (s *PersonalStats) statsValues(pattern StatsPattern) (WGShipStatsValues, WGPlayerStatsValues) {
	switch pattern {
	case StatsPatternPvPAll:
		return s.useShipStats.Pvp, s.accountInfo.Statistics.Pvp
	case StatsPatternPvPSolo:
		return s.useShipStats.PvpSolo, s.accountInfo.Statistics.PvpSolo
	case StatsPatternRankSolo:
		return s.useShipStats.RankSolo, s.accountInfo.Statistics.RankSolo
	}

	return WGShipStatsValues{}, WGPlayerStatsValues{}
}

func (s *PersonalStats) statsValuesForm(statsData WGShipsStatsData, pattern StatsPattern) WGShipStatsValues {
	switch pattern {
	case StatsPatternPvPAll:
		return statsData.Pvp
	case StatsPatternPvPSolo:
		return statsData.PvpSolo
	case StatsPatternRankSolo:
		return statsData.RankSolo
	}

	return WGShipStatsValues{}
}

func (s *PersonalStats) pr(
	actual PRFactor,
	expected PRFactor,
	battles uint,
) RatingValue {
	if battles < 1 {
		return NewRatingValueNone()
	}

	ratio := PRFactor{
		damage: actual.damage / expected.damage,
		frags:  actual.frags / expected.frags,
		wins:   actual.wins / expected.wins,
	}

	if !ratio.Valid() {
		return NewRatingValueNone()
	}

	norm := PRFactor{
		damage: math.Max(0, (ratio.damage-0.4)/(1-0.4)),
		frags:  math.Max(0, (ratio.frags-0.1)/(1-0.1)),
		wins:   math.Max(0, (ratio.wins-0.7)/(1-0.7)),
	}

	value := 700*norm.damage + 300*norm.frags + 150*norm.wins
	rating := NewRatingFromPR(value)

	return NewRatingValue(value, rating)
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

func platoonRate(allBattles uint, soloBattles uint, div2Battles uint, div3Battles uint) float64 {
	soloRate := div(soloBattles, allBattles) * 1
	div2Rate := div(div2Battles, allBattles) * 2
	div3Rate := div(div3Battles, allBattles) * 3
	return soloRate + div2Rate + div3Rate
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
