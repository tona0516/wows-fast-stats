package domain

type ShipTier struct {
	ValueObject[uint]
}

func NewShipTier(value uint) ShipTier {
	return ShipTier{ValueObject[uint]{value}}
}

func (t ShipTier) TierGroup() TierGroup {
	v := t.value
	switch {
	case v >= 1 && v <= 4:
		return ShipTierLow
	case v >= 5 && v <= 7:
		return ShipTierMiddle
	case v >= 8:
		return ShipTierHigh
	default:
		return ShipTierNone
	}
}
