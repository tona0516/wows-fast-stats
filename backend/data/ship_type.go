package data

type ShipType string

const (
	ShipTypeCV   ShipType = "cv"
	ShipTypeBB   ShipType = "bb"
	ShipTypeCL   ShipType = "cl"
	ShipTypeDD   ShipType = "dd"
	ShipTypeSS   ShipType = "ss"
	ShipTypeAUX  ShipType = "aux"
	ShipTypeNONE ShipType = "none"
)

//nolint:gochecknoglobals
var (
	shipTypeNames = map[string]ShipType{
		"AirCarrier": ShipTypeCV,
		"Battleship": ShipTypeBB,
		"Cruiser":    ShipTypeCL,
		"Destroyer":  ShipTypeDD,
		"Submarine":  ShipTypeSS,
		"Auxiliary":  ShipTypeAUX,
	}
	shipTypePriorities = []ShipType{
		ShipTypeCV,
		ShipTypeBB,
		ShipTypeCL,
		ShipTypeDD,
		ShipTypeSS,
		ShipTypeAUX,
	}
)

func NewShipType(raw string) ShipType {
	shipType, ok := shipTypeNames[raw]
	if !ok {
		return ShipTypeNONE
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
