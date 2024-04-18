package data

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
		return one.Type.Priority() < second.Type.Priority()
	}

	if one.Tier != second.Tier {
		return one.Tier > second.Tier
	}

	if one.Nation != second.Nation {
		return one.Nation.Priority() < second.Nation.Priority()
	}

	return one.Name < second.Name
}
