package vo

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

	shipTypes := shipTypes()

	if one.Type != second.Type {
		if _, ok := shipTypes[one.Type]; !ok {
			return true
		}

		if _, ok := shipTypes[second.Type]; !ok {
			return false
		}

		return shipTypes[one.Type] < shipTypes[second.Type]
	}

	if one.Tier != second.Tier {
		return one.Tier > second.Tier
	}

	if one.Nation != second.Nation {
		return one.Nation < second.Nation
	}

	return one.Name < second.Name
}

func shipTypes() map[string]int {
	return map[string]int{
		"AirCarrier": 0,
		"Battleship": 1,
		"Cruiser":    2,
		"Destroyer":  3,
		"Submarine":  4,
		"Auxiliary":  5,
	}
}
