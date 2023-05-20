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

func (s *Stats) ShipPR() float64 {
	rating := Rating{}

	battles := s.ShipsStats.Pvp.Battles
	return rating.PersonalRating(
		RatingFactor{
			AvgDamage: avgDamage(s.ShipsStats.Pvp.DamageDealt, battles),
			AvgFrags:  avgFrags(s.ShipsStats.Pvp.Frags, battles),
			WinRate:   winRate(s.ShipsStats.Pvp.Wins, battles),
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
	battles := values.Battles
	survivedBattles := values.SurvivedBattles

	death := battles - survivedBattles
	if death < 1 {
		return float64(frags)
	}

	return float64(frags) / float64(death)
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

func (s *Stats) AvgTier(
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) float64 {
	var sum uint
	var battles uint
	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		shipID := ship.ShipID
		tier := shipInfo[shipID].Tier
		sum += ship.Pvp.Battles * tier
		battles += ship.Pvp.Battles
	}

	if battles == 0 {
		return 0
	}

	return float64(sum) / float64(battles)
}

func (s *Stats) UsingTierRate(
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) vo.TierGroup {
	lowKey := "low"
	middleKey := "middle"
	highKey := "high"
	tierMap := make(map[string]uint)

	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		battles := ship.Pvp.Battles
		if battles == 0 {
			continue
		}

		shipID := ship.ShipID
		tier := shipInfo[shipID].Tier

		switch {
		case tier >= 1 && tier <= 4:
			tierMap[lowKey] += battles
		case tier >= 5 && tier <= 7:
			tierMap[middleKey] += battles
		case tier >= 8:
			tierMap[highKey] += battles
		}
	}

	var allBattles uint
	for _, v := range tierMap {
		allBattles += v
	}
	if allBattles == 0 {
		return vo.TierGroup{}
	}

	return vo.TierGroup{
		Low:    float64(tierMap[lowKey]) / float64(allBattles) * 100,
		Middle: float64(tierMap[middleKey]) / float64(allBattles) * 100,
		High:   float64(tierMap[highKey]) / float64(allBattles) * 100,
	}
}

func (s *Stats) UsingShipTypeRate(
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) vo.ShipTypeGroup {
	battlesMap := make(map[vo.ShipType]uint)

	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		battles := ship.Pvp.Battles
		if battles == 0 {
			continue
		}

		shipID := ship.ShipID
		shipType := shipInfo[shipID].Type

		battlesMap[shipType] += battles
	}

	var allBattles uint
	for _, v := range battlesMap {
		allBattles += v
	}
	if allBattles == 0 {
		return vo.ShipTypeGroup{}
	}

	return vo.ShipTypeGroup{
		SS: float64(battlesMap[vo.SS]) / float64(allBattles) * 100,
		DD: float64(battlesMap[vo.DD]) / float64(allBattles) * 100,
		CL: float64(battlesMap[vo.CL]) / float64(allBattles) * 100,
		BB: float64(battlesMap[vo.BB]) / float64(allBattles) * 100,
		CV: float64(battlesMap[vo.CV]) / float64(allBattles) * 100,
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

func avgDamage(damageDealt uint, battles uint) float64 {
	if battles < 1 {
		return 0
	}

	return float64(damageDealt) / float64(battles)
}

func avgFrags(frags uint, battles uint) float64 {
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
