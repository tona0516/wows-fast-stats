package domain

import (
	"changeme/backend/vo"
)

type StatsMode string

const (
	ModeShip        StatsMode = "ship"
	ModeShipSolo    StatsMode = "ship_solo"
	ModeOverall     StatsMode = "overall"
	ModeOverallSolo StatsMode = "overall_solo"
)

type Stats struct {
	AccountInfo vo.WGAccountInfoData
	ShipsStats  vo.WGShipsStatsData
	Expected    vo.NSExpectedStatsData
}

func (s *Stats) SetShipStats(shipStats vo.WGShipsStatsData) {
	s.ShipsStats = shipStats
}

func (s *Stats) ShipPR(mode StatsMode) float64 {
	rating := Rating{}
	values := statsValuesFrom(mode, s.ShipsStats)
	battles := values.Battles

	if battles < 1 {
		return -1
	}

	return rating.PersonalRating(
		RatingFactor{
			AvgDamage: avgDamage(values.DamageDealt, battles),
			AvgFrags:  avgKill(values.Frags, battles),
			WinRate:   winRate(values.Wins, battles),
		},
		RatingFactor{
			AvgDamage: s.Expected.AverageDamageDealt,
			AvgFrags:  s.Expected.AverageFrags,
			WinRate:   s.Expected.WinRate,
		},
	)
}

func (s *Stats) Battles(mode StatsMode) uint {
	values := s.statsValues(mode)
	return values.Battles
}

func (s *Stats) AvgDamage(mode StatsMode) float64 {
	values := s.statsValues(mode)
	damageDealt := values.DamageDealt
	battles := values.Battles
	return avgDamage(damageDealt, battles)
}

func (s *Stats) KdRate(mode StatsMode) float64 {
	values := s.statsValues(mode)
	frags := values.Frags

	death := values.Battles - values.SurvivedBattles
	if death < 1 {
		return float64(frags)
	}

	return float64(frags) / float64(death)
}

func (s *Stats) AvgKill(mode StatsMode) float64 {
	values := s.statsValues(mode)
	return avgKill(values.Frags, values.Battles)
}

func (s *Stats) AvgExp(mode StatsMode) float64 {
	values := s.statsValues(mode)
	xp := values.Xp
	battles := values.Battles

	if battles < 1 {
		return 0
	}

	return float64(xp) / float64(battles)
}

func (s *Stats) WinRate(mode StatsMode) float64 {
	values := s.statsValues(mode)
	wins := values.Wins
	battles := values.Battles

	return winRate(wins, battles)
}

func (s *Stats) WinSurvivedRate(mode StatsMode) float64 {
	values := s.statsValues(mode)
	wins := values.Wins
	survivedWins := values.SurvivedWins

	if wins < 1 {
		return 0
	}

	return float64(survivedWins) / float64(wins) * 100
}

func (s *Stats) LoseSurvivedRate(mode StatsMode) float64 {
	values := s.statsValues(mode)
	wins := values.Wins
	survivedWins := values.SurvivedWins
	battles := values.Battles
	survivedBattles := values.SurvivedBattles

	loses := battles - wins
	if loses < 1 {
		return 0
	}

	return float64(survivedBattles-survivedWins) / float64(loses) * 100
}

func (s *Stats) MainBatteryHitRate(mode StatsMode) float64 {
	values := s.statsValues(mode)
	hits := values.MainBattery.Hits
	shots := values.MainBattery.Shots

	if shots < 1 {
		return 0
	}

	return float64(hits) / float64(shots) * 100
}

func (s *Stats) TorpedoesHitRate(mode StatsMode) float64 {
	values := s.statsValues(mode)
	hits := values.Torpedoes.Hits
	shots := values.Torpedoes.Shots

	if shots < 1 {
		return 0
	}

	return float64(hits) / float64(shots) * 100
}

func (s *Stats) PlanesKilled(mode StatsMode) float64 {
	values := s.statsValues(mode)
	killed := values.PlanesKilled
	battles := values.Battles

	if battles < 1 {
		return 0
	}

	return float64(killed) / float64(battles)
}

func (s *Stats) AvgTier(
	mode StatsMode,
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) float64 {
	var sum uint
	var battles uint
	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		values := statsValuesFrom(mode, ship)
		shipID := ship.ShipID
		tier := shipInfo[shipID].Tier
		sum += values.Battles * tier
		battles += values.Battles
	}

	if battles < 1 {
		return 0
	}

	return float64(sum) / float64(battles)
}

func (s *Stats) UsingTierRate(
	mode StatsMode,
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) vo.TierGroup {
	lowKey := "low"
	middleKey := "middle"
	highKey := "high"
	tierMap := make(map[string]float64)

	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		values := statsValuesFrom(mode, ship)
		battles := values.Battles
		if battles == 0 {
			continue
		}

		shipID := ship.ShipID
		tier := shipInfo[shipID].Tier

		var key string
		switch {
		case tier >= 1 && tier <= 4:
			key = lowKey
		case tier >= 5 && tier <= 7:
			key = middleKey
		case tier >= 8:
			key = highKey
		}
		tierMap[key] += float64(battles)
	}

	var allBattles float64
	for _, v := range tierMap {
		allBattles += v
	}
	if allBattles < 1 {
		return vo.TierGroup{}
	}

	return vo.TierGroup{
		Low:    tierMap[lowKey] / allBattles * 100,
		Middle: tierMap[middleKey] / allBattles * 100,
		High:   tierMap[highKey] / allBattles * 100,
	}
}

func (s *Stats) UsingShipTypeRate(
	mode StatsMode,
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) vo.ShipTypeGroup {
	battlesMap := make(map[vo.ShipType]float64)

	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		values := statsValuesFrom(mode, ship)
		battles := values.Battles
		if battles == 0 {
			continue
		}

		shipID := ship.ShipID
		shipType := shipInfo[shipID].Type

		battlesMap[shipType] += float64(battles)
	}

	var allBattles float64
	for _, v := range battlesMap {
		allBattles += v
	}
	if allBattles < 1 {
		return vo.ShipTypeGroup{}
	}

	return vo.ShipTypeGroup{
		SS: battlesMap[vo.SS] / allBattles * 100,
		DD: battlesMap[vo.DD] / allBattles * 100,
		CL: battlesMap[vo.CL] / allBattles * 100,
		BB: battlesMap[vo.BB] / allBattles * 100,
		CV: battlesMap[vo.CV] / allBattles * 100,
	}
}

func (s *Stats) statsValues(mode StatsMode) vo.WGStatsValues {
	switch mode {
	case ModeShip:
		return s.ShipsStats.Pvp
	case ModeShipSolo:
		return s.ShipsStats.PvpSolo
	case ModeOverall:
		return s.AccountInfo.Statistics.Pvp
	case ModeOverallSolo:
		return s.AccountInfo.Statistics.PvpSolo
	}

	return vo.WGStatsValues{}
}

func statsValuesFrom(mode StatsMode, s vo.WGShipsStatsData) vo.WGStatsValues {
	switch mode {
	case ModeShip:
		return s.Pvp
	case ModeShipSolo:
		return s.PvpSolo
	case ModeOverall:
		return s.Pvp
	case ModeOverallSolo:
		return s.PvpSolo
	}

	return vo.WGStatsValues{}
}

func avgDamage(damageDealt uint, battles uint) float64 {
	if battles < 1 {
		return 0
	}

	return float64(damageDealt) / float64(battles)
}

func avgKill(frags uint, battles uint) float64 {
	if battles < 1 {
		return 0
	}

	return float64(frags) / float64(battles)
}

func winRate(wins uint, battles uint) float64 {
	if battles < 1 {
		return 0
	}

	return float64(wins) / float64(battles) * 100
}
