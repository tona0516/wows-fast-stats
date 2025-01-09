package model

type Warships map[int]Warship

type Warship struct {
	ID            int
	Name          string
	Tier          uint
	Type          ShipType
	Nation        Nation
	IsPremium     bool
	AverageDamage float64
	AverageFrags  float64
	WinRate       float64
}
