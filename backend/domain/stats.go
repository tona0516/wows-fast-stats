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
}
type Stats struct {
	Ship   StatsFactor
	Player StatsFactor
}

func (s *Stats) SetShipStats(ship StatsFactor) {
    s.Ship = ship
}

func (s *Stats) ShipAvgDamage() float64 {
	if s.Ship.Battles != 0 {
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
	if s.Ship.Battles != 0 {
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
        return float64(s.Ship.SurvivedBattles - s.Ship.SurvivedWins) / float64(loses) * 100
    }
    return 0
}

func (s *Stats) PlayerAvgDamage() float64 {
	if s.Player.Battles > 0 {
		return float64(s.Player.DamageDealt) / float64(s.Player.Battles)
	}
	return 0
}

func (s *Stats) PlayerKdRate() float64 {
	if s.Player.Battles-s.Player.SurvivedBattles > 0 {
		return float64(s.Player.Frags) / float64(s.Player.Battles-s.Player.SurvivedBattles)
	}
	return 0
}

func (s *Stats) PlayerAvgExp() float64 {
    if s.Player.Battles > 0 {
        return float64(s.Player.Xp) / float64(s.Player.Battles)
    }
    return 0
}

func (s *Stats) PlayerWinRate() float64 {
	if s.Player.Battles != 0 {
		return float64(s.Player.Wins) / float64(s.Player.Battles) * 100
	}
	return 0
}

func (s *Stats) PlayerWinSurvivedRate() float64 {
    if s.Player.Wins > 0 {
        return float64(s.Player.SurvivedWins) / float64(s.Player.Wins) * 100
    }
    return 0
}

func (s *Stats) PlayerLoseSurvivedRate() float64 {
    loses := s.Player.Battles - s.Player.Wins
    if loses > 0 {
        return float64(s.Player.SurvivedBattles - s.Player.SurvivedWins) / float64(loses) * 100
    }
    return 0
}

func (s *Stats) PlayerAvgTier(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) float64 {
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

func (s *Stats) UsingTierRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.TierGroup[float64] {
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

func (s *Stats) UsingShipTypeRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.ShipTypeGroup[float64] {
    var result vo.ShipTypeGroup[float64]
    var allBattles uint

    playerShipStats := shipStats[accountID].Data[accountID]
	for _, ship := range playerShipStats {
        shipID := ship.ShipID
        shipType := shipInfo[shipID].Type
        battles := ship.Pvp.Battles
        if battles != 0 {
            switch shipType {
            case "Submarine":
                result.SS += float64(battles)
                allBattles += battles
            case "Destroyer":
                result.DD += float64(battles)
                allBattles += battles
            case "Cruiser":
                result.CL += float64(battles)
                allBattles += battles
            case "Battleship":
                result.BB += float64(battles)
                allBattles += battles
            case "AirCarrier":
                result.CV += float64(battles)
                allBattles += battles
            }
        }
	}

    if allBattles == 0 {
        return vo.ShipTypeGroup[float64]{}
    }

    result.SS = result.SS / float64(allBattles) * 100
    result.DD = result.DD / float64(allBattles) * 100
    result.CL = result.CL / float64(allBattles) * 100
    result.BB = result.BB / float64(allBattles) * 100
    result.CV = result.CV / float64(allBattles) * 100

    return result
}
