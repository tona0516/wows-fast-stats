package domain

type TierGroup uint

const (
	ShipTierLow TierGroup = iota + 1
	ShipTierMiddle
	ShipTierHigh
	ShipTierNone TierGroup = 999
)
