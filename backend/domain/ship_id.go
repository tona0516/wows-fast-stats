package domain

type ShipID struct {
	ValueObject[uint]
}

func NewShipID(value uint) ShipID {
	return ShipID{ValueObject[uint]{value}}
}
