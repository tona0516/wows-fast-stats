package vo

type ShipType string

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
	switch s {
	case CV:
		return 0
	case BB:
		return 1
	case CL:
		return 2
	case DD:
		return 3
	case SS:
		return 4
	case AUX:
		return 5
	case NONE:
		return 999
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
