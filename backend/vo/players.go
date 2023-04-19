package vo

var SHIP_TYPES = map[string]int{
    "AirCarrier": 0,
    "Battleship": 1,
    "Cruiser": 2,
    "Destroyer": 3,
    "Submarine": 4,
    "Auxiliary": 5,
};

type Players []Player

func (p Players) Len() int {
    return len(p)
}

func (p Players) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func (p Players) Less(i, j int) bool {
    one := p[i].ShipInfo
    second := p[j].ShipInfo
    if one.Type != second.Type {
        if _, ok := SHIP_TYPES[one.Type]; !ok {
            return true
        }
        if _, ok := SHIP_TYPES[second.Type]; !ok {
            return false
        }
        return SHIP_TYPES[one.Type] < SHIP_TYPES[second.Type]
    }
    if one.Tier != second.Tier {
        return one.Tier > second.Tier
    }
    if one.Nation != second.Nation {
        return one.Nation < second.Nation
    }
    return one.Name < second.Name
}

func (p Players) TeamAverage() TeamAverage {
    var result TeamAverage
    var nShip uint
    var nPlayer uint
    for _, v := range p {
        if v.ShipStats.Battles != 0 {
            result.PersonalRating += v.ShipStats.PersonalRating
            result.DamageByShip += v.ShipStats.AvgDamage
            result.WinRateByShip += v.ShipStats.WinRate
            result.KdRateByShip += v.ShipStats.KdRate
            nShip += 1
        }
        if v.PlayerStats.Battles != 0 {
            result.DamageByPlayer += v.PlayerStats.AvgDamage
            result.WinRateByPlayer += v.PlayerStats.WinRate
            result.KdRateByPlayer += v.PlayerStats.KdRate
            nPlayer += 1
        }
    }

    if nShip != 0 {
        result.PersonalRating = result.PersonalRating / float64(nShip)
        result.DamageByShip = result.DamageByShip / float64(nShip)
        result.WinRateByShip = result.WinRateByShip / float64(nShip)
        result.KdRateByShip = result.KdRateByShip / float64(nShip)
    }

    if nPlayer != 0 {
        result.DamageByPlayer = result.DamageByPlayer / float64(nPlayer)
        result.WinRateByPlayer = result.WinRateByPlayer / float64(nPlayer)
        result.KdRateByPlayer = result.KdRateByPlayer / float64(nPlayer)
    }

    return result
}
