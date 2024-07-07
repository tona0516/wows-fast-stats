package vo

type TierGroup int

const (
	ShipTierLow TierGroup = iota + 1
	ShipTierMiddle
	ShipTierHigh
	ShipTierNone TierGroup = 999
)
