package domain

import (
	"changeme/backend/vo"
)

type StatsFactor struct {
	Battles         uint
	SurvivedBattles uint
	DamageDealt     uint
	Frags           uint
	Wins            uint
	SurvivedWins    uint
	Xp              uint
	MainBattery     struct {
		Hits  uint
		Shots uint
	}
	Torpedoes struct {
		Hits  uint
		Shots uint
	}
}
type Stats struct {
	Ship    StatsFactor
	Overall StatsFactor
}

func (s *Stats) SetShipStats(ship StatsFactor) {
	s.Ship = ship
}

func (s *Stats) ShipAvgDamage() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.DamageDealt) / float64(s.Ship.Battles)
	}

	return 0
}

func (s *Stats) ShipKdRate() float64 {
	if s.Ship.Battles-s.Ship.SurvivedBattles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles-s.Ship.SurvivedBattles)
	}

	return 0
}

func (s *Stats) ShipAvgFrags() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles)
	}

	return 0
}

func (s *Stats) ShipAvgExp() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.Xp) / float64(s.Ship.Battles)
	}

	return 0
}

func (s *Stats) ShipWinRate() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.Wins) / float64(s.Ship.Battles) * 100
	}

	return 0
}

func (s *Stats) ShipWinSurvivedRate() float64 {
	if s.Ship.Wins > 0 {
		return float64(s.Ship.SurvivedWins) / float64(s.Ship.Wins) * 100
	}

	return 0
}

func (s *Stats) ShipLoseSurvivedRate() float64 {
	loses := s.Ship.Battles - s.Ship.Wins
	if loses > 0 {
		return float64(s.Ship.SurvivedBattles-s.Ship.SurvivedWins) / float64(loses) * 100
	}

	return 0
}

func (s *Stats) ShipMainBatteryHitRate() float64 {
	if s.Ship.MainBattery.Shots > 0 {
		return float64(s.Ship.MainBattery.Hits) / float64(s.Ship.MainBattery.Shots) * 100
	}

	return 0
}

func (s *Stats) ShipTorpedoesHitRate() float64 {
	if s.Ship.Torpedoes.Shots > 0 {
		return float64(s.Ship.Torpedoes.Hits) / float64(s.Ship.Torpedoes.Shots) * 100
	}

	return 0
}

func (s *Stats) OverallAvgDamage() float64 {
	if s.Overall.Battles > 0 {
		return float64(s.Overall.DamageDealt) / float64(s.Overall.Battles)
	}

	return 0
}

func (s *Stats) OverallKdRate() float64 {
	if s.Overall.Battles-s.Overall.SurvivedBattles > 0 {
		return float64(s.Overall.Frags) / float64(s.Overall.Battles-s.Overall.SurvivedBattles)
	}

	return 0
}

func (s *Stats) OverallAvgExp() float64 {
	if s.Overall.Battles > 0 {
		return float64(s.Overall.Xp) / float64(s.Overall.Battles)
	}

	return 0
}

func (s *Stats) OverallWinRate() float64 {
	if s.Overall.Battles > 0 {
		return float64(s.Overall.Wins) / float64(s.Overall.Battles) * 100
	}

	return 0
}

func (s *Stats) OverallWinSurvivedRate() float64 {
	if s.Overall.Wins > 0 {
		return float64(s.Overall.SurvivedWins) / float64(s.Overall.Wins) * 100
	}

	return 0
}

func (s *Stats) OverallLoseSurvivedRate() float64 {
	loses := s.Overall.Battles - s.Overall.Wins
	if loses > 0 {
		return float64(s.Overall.SurvivedBattles-s.Overall.SurvivedWins) / float64(loses) * 100
	}

	return 0
}

func (s *Stats) OverallAvgTier(
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

func (s *Stats) OverallUsingTierRate(
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) vo.TierGroup {
	var (
		low        uint
		middle     uint
		high       uint
		allBattles uint
	)

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
			low += battles
			allBattles += battles
		case tier >= 5 && tier <= 7:
			middle += battles
			allBattles += battles
		case tier >= 8:
			high += battles
			allBattles += battles
		}
	}

	if allBattles == 0 {
		return vo.TierGroup{}
	}

	return vo.TierGroup{
		Low:    float64(low) / float64(allBattles) * 100,
		Middle: float64(middle) / float64(allBattles) * 100,
		High:   float64(high) / float64(allBattles) * 100,
	}
}

//nolint:cyclop
func (s *Stats) OverallUsingShipTypeRate(
	accountID int,
	shipInfo map[int]vo.Warship,
	shipStats map[int]vo.WGShipsStats,
) vo.ShipTypeGroup {
	var (
		ss         uint
		dd         uint
		cl         uint
		bb         uint
		cv         uint
		allBattles uint
	)

	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		battles := ship.Pvp.Battles
		if battles == 0 {
			continue
		}

		shipID := ship.ShipID
		shipType := shipInfo[shipID].Type

		switch shipType {
		case vo.SS:
			ss += battles
			allBattles += battles
		case vo.DD:
			dd += battles
			allBattles += battles
		case vo.CL:
			cl += battles
			allBattles += battles
		case vo.BB:
			bb += battles
			allBattles += battles
		case vo.CV:
			cv += battles
			allBattles += battles
		case vo.AUX:
			continue
		case vo.NONE:
			continue
		}
	}

	if allBattles == 0 {
		return vo.ShipTypeGroup{}
	}

	return vo.ShipTypeGroup{
		SS: float64(ss) / float64(allBattles) * 100,
		DD: float64(dd) / float64(allBattles) * 100,
		CL: float64(cl) / float64(allBattles) * 100,
		BB: float64(bb) / float64(allBattles) * 100,
		CV: float64(cv) / float64(allBattles) * 100,
	}
}
