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

func (p Players) WinRateOfShips() float64 {
    var sum float64 = 0
    var size int = 0
    for _, v := range p {
        if v.ShipStats.Battles != 0 {
            sum += v.ShipStats.WinRate
            size += 1
        }
    }
    if size == 0 {
        return 0
    }
    return sum / float64(size)
}

func (p Players) WinRateOfPlayers() float64 {
    var sum float64 = 0
    var size int = 0
    for _, v := range p {
        if v.PlayerStats.Battles != 0 {
            sum += v.PlayerStats.WinRate
            size += 1
        }
    }
    if size == 0 {
        return 0
    }
    return sum / float64(size)
}
