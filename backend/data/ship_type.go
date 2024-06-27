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

func shipTypeNames() map[string]ShipType {
	return map[string]ShipType{
		"AirCarrier": ShipTypeCV,
		"Battleship": ShipTypeBB,
		"Cruiser":    ShipTypeCL,
		"Destroyer":  ShipTypeDD,
		"Submarine":  ShipTypeSS,
		"Auxiliary":  ShipTypeAUX,
	}
}

func shipTypePriorities() []ShipType {
	return []ShipType{
		ShipTypeCV,
		ShipTypeBB,
		ShipTypeCL,
		ShipTypeDD,
		ShipTypeSS,
		ShipTypeAUX,
	}
}

func NewShipType(raw string) ShipType {
	shipTypeNames := shipTypeNames()
	shipType, ok := shipTypeNames[raw]
	if !ok {
		return ShipTypeNONE
	}

	return shipType
}

func (s ShipType) Priority() int {
	shipTypePriorities := shipTypePriorities()
	for i, shipType := range shipTypePriorities {
		if shipType == s {
			return i
		}
	}

	return 999
}
