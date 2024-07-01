package domain

type Warship struct {
	id        ShipID
	name      string
	tier      ShipTier
	types     ShipType
	nation    Nation
	isPremium bool
	avgDamage ufloat
}

func NewWarship[T ShipID](
	id ShipID,
	name string,
	tier ShipTier,
	types ShipType,
	nation Nation,
	isPremium bool,
	avgDamage ufloat,
) *Warship {
	return &Warship{
		id:        id,
		name:      name,
		tier:      tier,
		types:     types,
		nation:    nation,
		isPremium: isPremium,
		avgDamage: avgDamage,
	}
}

func (w *Warship) Less(another *Warship) bool {
	if w.types != another.types {
		return w.types.Priority() < another.types.Priority()
	}

	if w.tier != another.tier {
		return w.tier.value > another.tier.value
	}

	if w.nation != another.nation {
		return w.nation.Priority() < another.nation.Priority()
	}

	return w.name < another.name
}

// Entity interface

func (w *Warship) ID() ShipID {
	return w.id
}

func (w *Warship) Equals(another *Warship) bool {
	return w.id == another.id
}
