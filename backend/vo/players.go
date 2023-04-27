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

func (p Players) TeamAverage() Player {
    var result Player
    var nShip uint
    var nPlayer uint
    for _, v := range p {
        if v.ShipStats.Battles != 0 {
            result.ShipStats.PR += v.ShipStats.PR
            result.ShipStats.Damage += v.ShipStats.Damage
            result.ShipStats.WinRate += v.ShipStats.WinRate
            result.ShipStats.KdRate += v.ShipStats.KdRate
            nShip += 1
        }
        if v.PlayerStats.Battles != 0 {
            result.PlayerStats.Damage += v.PlayerStats.Damage
            result.PlayerStats.WinRate += v.PlayerStats.WinRate
            result.PlayerStats.KdRate += v.PlayerStats.KdRate
            nPlayer += 1
        }
    }

    if nShip != 0 {
        result.ShipStats.PR = result.ShipStats.PR / float64(nShip)
        result.ShipStats.Damage = result.ShipStats.Damage / float64(nShip)
        result.ShipStats.WinRate = result.ShipStats.WinRate / float64(nShip)
        result.ShipStats.KdRate = result.ShipStats.KdRate / float64(nShip)
    }

    if nPlayer != 0 {
        result.PlayerStats.Damage = result.PlayerStats.Damage / float64(nPlayer)
        result.PlayerStats.WinRate = result.PlayerStats.WinRate / float64(nPlayer)
        result.PlayerStats.KdRate = result.PlayerStats.KdRate / float64(nPlayer)
    }

    return result
}
