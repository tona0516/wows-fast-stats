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

func (s *StatsCalculator) UsingTierRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.TierGroup[float64] {
	var result vo.TierGroup[float64]
	var allBattles uint = 0

	playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
		shipID := ship.ShipID
		tier := shipInfo[shipID].Tier
        battles := ship.Pvp.Battles
        if battles != 0 {
            switch {
            case tier >= 1 && tier <= 4:
                result.Low += float64(battles)
                allBattles += battles
            case tier >= 5 && tier <= 7:
                result.Middle += float64(battles)
                allBattles += battles
            case tier >= 8:
                result.High += float64(battles)
                allBattles += battles
            }
        }
	}

    if allBattles == 0 {
        return vo.TierGroup[float64]{}
    }

    result.Low = result.Low / float64(allBattles) * 100
    result.Middle = result.Middle / float64(allBattles) * 100
    result.High = result.High / float64(allBattles) * 100

    return result
}

func (s *StatsCalculator) UsingShipTypeRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.ShipTypeValue {
    var result vo.ShipTypeValue
    var allBattles uint

    playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
        shipID := ship.ShipID
        shipType := shipInfo[shipID].Type
        battles := ship.Pvp.Battles
        if battles != 0 {
            switch shipType {
            case "Submarine":
                result.Ss += float64(battles)
                allBattles += battles
            case "Destroyer":
                result.Dd += float64(battles)
                allBattles += battles
            case "Cruiser":
                result.Cl += float64(battles)
                allBattles += battles
            case "Battleship":
                result.Bb += float64(battles)
                allBattles += battles
            case "AirCarrier":
                result.Cv += float64(battles)
                allBattles += battles
            }
        }
	}

    if allBattles == 0 {
        return vo.ShipTypeValue{}
    }

    result.Ss = result.Ss / float64(allBattles) * 100
    result.Dd = result.Dd / float64(allBattles) * 100
    result.Cl = result.Cl / float64(allBattles) * 100
    result.Bb = result.Bb / float64(allBattles) * 100
    result.Cv = result.Cv / float64(allBattles) * 100

    return result
}
