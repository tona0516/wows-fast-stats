package domain

import (
	"changeme/backend/vo"
)


type Stats struct {
	Battles         uint
	SurvivedBattles uint
	DamageDealt     uint
	Frags           uint
	Wins            uint
    SurvivedWins    uint
    Xp              uint
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

func (s *StatsCalculator) ShipAvgExp() float64 {
    if s.Ship.Battles > 0 {
        return float64(s.Ship.Xp) / float64(s.Ship.Battles)
    }
    return 0
}

func (s *StatsCalculator) ShipWinRate() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.Wins) / float64(s.Ship.Battles) * 100
	}
	return 0
}

func (s *StatsCalculator) ShipWinSurvivedRate() float64 {
    if s.Ship.Wins > 0 {
        return float64(s.Ship.SurvivedWins) / float64(s.Ship.Wins) * 100
    }
    return 0
}

func (s *StatsCalculator) ShipLoseSurvivedRate() float64 {
    loses := s.Ship.Battles - s.Ship.Wins
    if loses > 0 {
        return float64(s.Ship.SurvivedBattles - s.Ship.SurvivedWins) / float64(loses) * 100
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

func (s *StatsCalculator) PlayerAvgExp() float64 {
    if s.Player.Battles > 0 {
        return float64(s.Player.Xp) / float64(s.Player.Battles)
    }
    return 0
}

func (s *StatsCalculator) PlayerWinRate() float64 {
	if s.Player.Battles != 0 {
		return float64(s.Player.Wins) / float64(s.Player.Battles) * 100
	}
	return 0
}

func (s *StatsCalculator) PlayerWinSurvivedRate() float64 {
    if s.Player.Wins > 0 {
        return float64(s.Player.SurvivedWins) / float64(s.Player.Wins) * 100
    }
    return 0
}

func (s *StatsCalculator) PlayerLoseSurvivedRate() float64 {
    loses := s.Player.Battles - s.Player.Wins
    if loses > 0 {
        return float64(s.Player.SurvivedBattles - s.Player.SurvivedWins) / float64(loses) * 100
    }
    return 0
}

func (s *StatsCalculator) PlayerAvgTier(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) float64 {
	var sum uint = 0
	var battles uint = 0
	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		shipID := ship.ShipID
		tier := shipInfo[shipID].Tier
		sum += ship.Pvp.Battles * tier
		battles += ship.Pvp.Battles
	}

	if battles == 0 {
		return 0
	} else {
		return float64(sum) / float64(battles)
	}
}

func (s *StatsCalculator) UsingShipTypeRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.ShipTypeValue {
    var result vo.ShipTypeValue
    shipTypeMap := make(map[string]float64, 0)
    var allBattles uint

    playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
        if ship.Pvp.Battles != 0 {
            shipID := ship.ShipID
            shipType := shipInfo[shipID].Type
            shipTypeMap[shipType] += 1
            allBattles += 1
        }
	}

    if allBattles == 0 {
        return result
    }

    for k, v := range shipTypeMap {
        shipTypeMap[k] = v / float64(allBattles) * 100
    }

    result.Ss = shipTypeMap["Submarine"]
    result.Dd = shipTypeMap["Destroyer"]
    result.Cl = shipTypeMap["Cruiser"]
    result.Bb = shipTypeMap["Battleship"]
    result.Cv = shipTypeMap["AirCarrier"]

    return result
}
