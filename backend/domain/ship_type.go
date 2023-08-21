package domain

type ShipType string

const (
	CV   ShipType = "cv"
	BB   ShipType = "bb"
	CL   ShipType = "cl"
	DD   ShipType = "dd"
	SS   ShipType = "ss"
	AUX  ShipType = "aux"
	NONE ShipType = "none"
)

//nolint:gochecknoglobals
var (
	shipTypeNames = map[string]ShipType{
		"AirCarrier": CV,
		"Battleship": BB,
		"Cruiser":    CL,
		"Destroyer":  DD,
		"Submarine":  SS,
		"Auxiliary":  AUX,
	}
	shipTypePriorities = []ShipType{
		CV,
		BB,
		CL,
		DD,
		SS,
		AUX,
	}
)

func NewShipType(raw string) ShipType {
	shipType, ok := shipTypeNames[raw]
	if !ok {
		return NONE
	}

	return shipType
}

func (s ShipType) Priority() int {
	for i, shipType := range shipTypePriorities {
		if shipType == s {
			return i
		}
	}

	return 999
}
