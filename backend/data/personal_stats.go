package data

import (
	"math"
)

type PersonalStats struct {
	useShipID        ShipID
	accountInfo      WGAccountInfoData
	playerShipStats  PlayerShipStats
	playerShipBadges PlayerShipBadges
	warships         Warships
}

func NewPersonalStats(
	useShipID ShipID,
	accountInfo WGAccountInfoData,
	playerShipStats PlayerShipStats,
	playerShipBadges PlayerShipBadges,
	warships Warships,
) *PersonalStats {
	return &PersonalStats{
		useShipID:        useShipID,
		accountInfo:      accountInfo,
		playerShipStats:  playerShipStats,
		playerShipBadges: playerShipBadges,
		warships:         warships,
	}
}

func (s *PersonalStats) PR(category StatsCategory, pattern StatsPattern) RatingValue {
	switch category {
	case StatsCategoryShip:
		values, _ := s.statsValues(pattern)
		battles := values.Battles
		if battles == 0 {
			return NewRatingValueNone()
		}

		warship, ok := s.warships[s.useShipID]
		if !ok {
			return NewRatingValueNone()
		}

		if warship.ServerAverage == nil {
			return NewRatingValueNone()
		}
		serverAverage := *warship.ServerAverage

		return s.pr(
			PRFactor{
				damage: safeDivide(values.DamageDealt, battles),
				frags:  safeDivide(values.Frags, battles),
				wins:   safeDivide(values.Wins, battles) * 100,
			},
			PRFactor{
				damage: serverAverage.Damage,
				frags:  serverAverage.Frags,
				wins:   serverAverage.WinRate,
			},
		)

	case StatsCategoryOverall:
		var (
			actual     PRFactor
			expected   PRFactor
			allBattles uint
		)

		for _, ship := range s.playerShipStats {
			values := s.statsValuesForm(ship, pattern)
			battles := values.Battles

			if battles == 0 {
				continue
			}

			warship, ok := s.warships[ship.ShipID]
			if !ok {
				continue
			}

			if warship.ServerAverage == nil {
				continue
			}
			serverAverage := *warship.ServerAverage

			actual.damage += float64(values.DamageDealt)
			actual.frags += float64(values.Frags)
			actual.wins += float64(values.Wins)

			expected.damage += serverAverage.Damage * float64(battles)
			expected.frags += serverAverage.Frags * float64(battles)
			expected.wins += serverAverage.WinRate / 100 * float64(battles)

			allBattles += battles
		}

		if allBattles == 0 {
			return NewRatingValueNone()
		}

		return s.pr(actual, expected)
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
		value := safeDivide(ship.DamageDealt, ship.Battles)
		rating := NewRatingFromShipDamage(value, s.warships[s.useShipID].ServerAverage.Damage)
		return NewRatingValue(value, rating)
	case StatsCategoryOverall:
		value := safeDivide(player.DamageDealt, player.Battles)
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
		return safeDivide(ship.Frags, ship.Battles)
	case StatsCategoryOverall:
		return safeDivide(player.Frags, player.Battles)
	}

	return 0
}

func (s *PersonalStats) AvgExp(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return safeDivide(ship.Xp, ship.Battles)
	case StatsCategoryOverall:
		return safeDivide(player.Xp, player.Battles)
	}

	return 0
}

func (s *PersonalStats) WinRate(category StatsCategory, pattern StatsPattern) RatingValue {
	ship, player := s.statsValues(pattern)

	var value float64
	switch category {
	case StatsCategoryShip:
		value = safeDivide(ship.Wins, ship.Battles) * 100
	case StatsCategoryOverall:
		value = safeDivide(player.Wins, player.Battles) * 100
	}

	rating := NewRatingFromWinRate(value)
	return NewRatingValue(value, rating)
}

func (s *PersonalStats) SurvivedRate(category StatsCategory, pattern StatsPattern) SurvivedRate {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return SurvivedRate{
			All:  safeDivide(ship.SurvivedBattles, ship.Battles) * 100,
			Win:  safeDivide(ship.SurvivedWins, ship.Wins) * 100,
			Lose: safeDivide(ship.SurvivedBattles-ship.SurvivedWins, ship.Battles-ship.Wins) * 100,
		}
	case StatsCategoryOverall:
		return SurvivedRate{
			All:  safeDivide(player.SurvivedBattles, player.Battles) * 100,
			Win:  safeDivide(player.SurvivedWins, player.Wins) * 100,
			Lose: safeDivide(player.SurvivedBattles-player.SurvivedWins, player.Battles-player.Wins) * 100,
		}
	}

	return SurvivedRate{}
}

func (s *PersonalStats) HitRate(pattern StatsPattern) HitRate {
	ship, _ := s.statsValues(pattern)
	return HitRate{
		MainBattery: safeDivide(ship.MainBattery.Hits, ship.MainBattery.Shots) * 100,
		Torpedoes:   safeDivide(ship.Torpedoes.Hits, ship.Torpedoes.Shots) * 100,
	}
}

func (s *PersonalStats) PlanesKilled(pattern StatsPattern) float64 {
	ship, _ := s.statsValues(pattern)
	return safeDivide(ship.PlanesKilled, ship.Battles)
}

func (s *PersonalStats) AvgTier(
	pattern StatsPattern,
) float64 {
	var (
		sum        uint
		allBattles uint
	)

	for _, stats := range s.playerShipStats {
		warship, ok := s.warships[stats.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(stats, pattern)
		sum += values.Battles * warship.Tier
		allBattles += values.Battles
	}

	return safeDivide(sum, allBattles)
}

func (s *PersonalStats) UsingTierRate(
	pattern StatsPattern,
) TierGroup {
	tierGroupMap := make(map[string]uint)

	for _, ship := range s.playerShipStats {
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
		Low:    safeDivide(tierGroupMap["low"], allBattles) * 100,
		Middle: safeDivide(tierGroupMap["middle"], allBattles) * 100,
		High:   safeDivide(tierGroupMap["high"], allBattles) * 100,
	}
}

func (s *PersonalStats) UsingShipTypeRate(
	pattern StatsPattern,
) ShipTypeGroup {
	shipTypeMap := make(map[ShipType]uint)

	for _, ship := range s.playerShipStats {
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
		SS: safeDivide(shipTypeMap[ShipTypeSS], allBattles) * 100,
		DD: safeDivide(shipTypeMap[ShipTypeDD], allBattles) * 100,
		CL: safeDivide(shipTypeMap[ShipTypeCL], allBattles) * 100,
		BB: safeDivide(shipTypeMap[ShipTypeBB], allBattles) * 100,
		CV: safeDivide(shipTypeMap[ShipTypeCV], allBattles) * 100,
	}
}

func (s *PersonalStats) PlatoonRate(
	category StatsCategory,
) float64 {
	var (
		allBattles  uint
		soloBattles uint
		div2Battles uint
		div3Battles uint
	)

	switch category {
	case StatsCategoryShip:
		stats := s.playerShipStats[s.useShipID]
		allBattles = stats.Pvp.Battles
		soloBattles = stats.PvpSolo.Battles
		div2Battles = stats.PvpDiv2.Battles
		div3Battles = stats.PvpDiv3.Battles
	case StatsCategoryOverall:
		stats := s.accountInfo.Statistics
		allBattles = stats.Pvp.Battles
		soloBattles = stats.PvpSolo.Battles
		div2Battles = stats.PvpDiv2.Battles
		div3Battles = stats.PvpDiv3.Battles
	}

	soloRate := safeDivide(soloBattles, allBattles) * 1
	div2Rate := safeDivide(div2Battles, allBattles) * 2
	div3Rate := safeDivide(div3Battles, allBattles) * 3
	return soloRate + div2Rate + div3Rate
}

func (s *PersonalStats) EfficiencyBadge() EfficiencyBadge {
	for _, b := range s.playerShipBadges {
		if b.ShipID == s.useShipID {
			return s.toEfficiencyBadge(b.TopGradeClass)
		}
	}

	return EfficiencyBadgeNone
}

func (s *PersonalStats) EfficiencyBadges() EfficiencyBadgeGroup {
	var badges EfficiencyBadgeGroup

	for _, b := range s.playerShipBadges {
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
		return s.playerShipStats[s.useShipID].Pvp, s.accountInfo.Statistics.Pvp
	case StatsPatternPvPSolo:
		return s.playerShipStats[s.useShipID].PvpSolo, s.accountInfo.Statistics.PvpSolo
	case StatsPatternRankSolo:
		return s.playerShipStats[s.useShipID].RankSolo, s.accountInfo.Statistics.RankSolo
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
) RatingValue {
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
