package domain

type ShipType string

//nolint:gochecknoglobals
var shipTypePriority = map[ShipType]int{
	CV:   0,
	BB:   1,
	CL:   2,
	DD:   3,
	SS:   4,
	AUX:  5,
	NONE: 999,
}

func NewShipType(raw string) ShipType {
	switch raw {
	case "AirCarrier":
		return CV
	case "Battleship":
		return BB
	case "Cruiser":
		return CL
	case "Destroyer":
		return DD
	case "Submarine":
		return SS
	case "Auxiliary":
		return AUX
	}

	return NONE
}

func (s ShipType) Priority() int {
	if priority, ok := shipTypePriority[s]; ok {
		return priority
	}

	return 999
}

const (
	CV   ShipType = "cv"
	BB   ShipType = "bb"
	CL   ShipType = "cl"
	DD   ShipType = "dd"
	SS   ShipType = "ss"
	AUX  ShipType = "aux"
	NONE ShipType = "none"
)
