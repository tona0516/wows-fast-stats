package vo

import (
	"errors"
	"strconv"
)

const (
	ShipTypeCV = iota + 1
	ShipTypeBB
	ShipTypeCL
	ShipTypeDD
	ShipTypeSS
	ShipTypeAux
	ShipTypeNone = 999
)

var _shipTypes = []int{
	ShipTypeCV,
	ShipTypeBB,
	ShipTypeCL,
	ShipTypeDD,
	ShipTypeSS,
	ShipTypeAux,
	ShipTypeNone,
}

type ShipType struct {
	ValueObject[int]
}

func NewShipType(value int) (ShipType, error) {
	for _, shipType := range _shipTypes {
		if shipType == value {
			return ShipType{ValueObject[int]{shipType}}, nil
		}
	}

	return ShipType{}, errors.New("invalid_ship_type: " + strconv.Itoa(value))
}

func (t ShipType) Priority() int {
	return t.value
}
