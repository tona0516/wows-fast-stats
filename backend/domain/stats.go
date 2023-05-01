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
	Overall StatsFactor
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
	if s.Overall.Battles != 0 {
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
        return float64(s.Overall.SurvivedBattles - s.Overall.SurvivedWins) / float64(loses) * 100
    }
    return 0
}

func (s *Stats) OverallAvgTier(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) float64 {
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

func (s *Stats) OverallUsingTierRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.TierGroup {
	var result vo.TierGroup
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
        return vo.TierGroup{}
    }

    result.Low = result.Low / float64(allBattles) * 100
    result.Middle = result.Middle / float64(allBattles) * 100
    result.High = result.High / float64(allBattles) * 100

    return result
}

func (s *Stats) OverallUsingShipTypeRate(accountID int, shipInfo map[int]vo.Warship, shipStats map[int]vo.WGShipsStats) vo.ShipTypeGroup {
    var result vo.ShipTypeGroup
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
        return vo.ShipTypeGroup{}
    }

    result.SS = result.SS / float64(allBattles) * 100
    result.DD = result.DD / float64(allBattles) * 100
    result.CL = result.CL / float64(allBattles) * 100
    result.BB = result.BB / float64(allBattles) * 100
    result.CV = result.CV / float64(allBattles) * 100

    return result
}
