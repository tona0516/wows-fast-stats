package domain

import "changeme/backend/vo"


type Stats struct {
	Battles         uint
	SurvivedBattles uint
	DamageDealt     uint
	Frags           uint
	Wins            uint
}
type StatsCalculator struct {
	Ship   Stats
	Player Stats
}

func (s *StatsCalculator) SetShipStats(ship Stats) {
    s.Ship = ship
}

func (s *StatsCalculator) ShipAvgDamage() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.DamageDealt) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *StatsCalculator) ShipKdRate() float64 {
	if s.Ship.Battles-s.Ship.SurvivedBattles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles-s.Ship.SurvivedBattles)
	}
	return 0
}

func (s *StatsCalculator) ShipAvgFrags() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *StatsCalculator) ShipWinRate() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.Wins) / float64(s.Ship.Battles) * 100
	}
	return 0
}

func (s *StatsCalculator) PlayerAvgDamage() float64 {
	if s.Player.Battles > 0 {
		return float64(s.Player.DamageDealt) / float64(s.Player.Battles)
	}
	return 0
}

func (s *StatsCalculator) PlayerKdRate() float64 {
	if s.Player.Battles-s.Player.SurvivedBattles > 0 {
		return float64(s.Player.Frags) / float64(s.Player.Battles-s.Player.SurvivedBattles)
	}
	return 0
}

func (s *StatsCalculator) PlayerWinRate() float64 {
	if s.Player.Battles != 0 {
		return float64(s.Player.Wins) / float64(s.Player.Battles) * 100
	}
	return 0
}

func (s *StatsCalculator) PlayerAvgTier(accountID int, shipInfo map[int]vo.ShipInfo, shipStats map[int]vo.WGShipsStats) float64 {
	var sum uint = 0
	var battles uint = 0
	playerShipStats := shipStats[accountID].Data[accountID]
	for i := range playerShipStats {
		shipID := playerShipStats[i].ShipID
		tier := shipInfo[shipID].Tier
		sum += playerShipStats[i].Pvp.Battles * tier
		battles += playerShipStats[i].Pvp.Battles
	}

	if battles == 0 {
		return 0
	} else {
		return float64(sum) / float64(battles)
	}
}
