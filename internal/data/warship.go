package data

type Warships map[int]Warship

type Warship struct {
	Name      string
	Tier      uint
	Type      ShipType
	Nation    Nation
	IsPremium bool
}
