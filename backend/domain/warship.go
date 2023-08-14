package domain

type Warship struct {
	Name      string
	Tier      uint
	Type      ShipType
	Nation    Nation
	IsPremium bool
}
