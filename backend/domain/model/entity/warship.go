package entity

import "wfs/backend/domain/model/vo"

type Warship struct {
	id        vo.ShipID
	name      string
	tier      vo.ShipTier
	types     vo.ShipType
	nation    vo.Nation
	isPremium bool
	avgDamage vo.UFloat
}

func NewWarship[T vo.ShipID](
	id vo.ShipID,
	name string,
	tier vo.ShipTier,
	types vo.ShipType,
	nation vo.Nation,
	isPremium bool,
	avgDamage vo.UFloat,
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

func (w *Warship) less(another *Warship) bool {
	if w.types != another.types {
		return w.types.Priority() < another.types.Priority()
	}

	if w.tier != another.tier {
		return w.tier.Value() > another.tier.Value()
	}

	if w.nation != another.nation {
		return w.nation.Priority() < another.nation.Priority()
	}

	return w.name < another.name
}

// Entity interface

func (w *Warship) ID() vo.ShipID {
	return w.id
}

func (w *Warship) Equals(another *Warship) bool {
	return w.id == another.id
}
