package domain

import "changeme/backend/vo"


type Stats struct {
	Battles         uint
	SurvivedBattles uint
	DamageDealt     uint
	Frags           uint
	Xp              uint
	Wins            uint
}
type SummaryStats struct {
	Ship   Stats
	Player Stats
}

func (s *SummaryStats) ShipAvgDamage() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.DamageDealt) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *SummaryStats) ShipAvgExp() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.Xp) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *SummaryStats) ShipKdRate() float64 {
	if s.Ship.Battles-s.Ship.SurvivedBattles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles-s.Ship.SurvivedBattles)
	}
	return 0
}

func (s *SummaryStats) ShipAvgFrags() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *SummaryStats) ShipWinRate() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.Wins) / float64(s.Ship.Battles) * 100
	}
	return 0
}

func (s *SummaryStats) PlayerAvgDamage() float64 {
	if s.Player.Battles > 0 {
		return float64(s.Player.DamageDealt) / float64(s.Player.Battles)
	}
	return 0
}

func (s *SummaryStats) PlayerAvgExp() float64 {
	if s.Player.Battles > 0 {
		return float64(s.Player.Xp) / float64(s.Player.Battles)
	}
	return 0
}

func (s *SummaryStats) PlayerKdRate() float64 {
	if s.Player.Battles-s.Player.SurvivedBattles > 0 {
		return float64(s.Player.Frags) / float64(s.Player.Battles-s.Player.SurvivedBattles)
	}
	return 0
}

func (s *SummaryStats) PlayerWinRate() float64 {
	if s.Player.Battles != 0 {
		return float64(s.Player.Wins) / float64(s.Player.Battles) * 100
	}
	return 0
}

func (s *SummaryStats) PlayerAvgTier(accountID int, shipInfo map[int]vo.ShipInfo, shipStats map[int]vo.WGShipsStats) float64 {
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
