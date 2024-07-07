package vo

import "errors"

type ShipTypeRate struct {
	values map[ShipType]Percentage
}

func NewShipTypeRate(
	battles uint,
	battlesCV uint,
	battlesBB uint,
	battlesCL uint,
	battlesDD uint,
	battlesSS uint,
) (ShipTypeRate, error) {
	if battles < 1 {
		return ShipTypeRate{}, errors.New("no_battles")
	}

	calcRate := func(numerator uint, denominator uint) (Percentage, error) {
		return NewPercentage(float64(numerator) / float64(denominator) * 100)
	}

	battlesNone := battles - (battlesCV + battlesBB + battlesCL + battlesDD + battlesSS)
	eachBattles := []struct {
		shipType int
		battles  uint
	}{
		{
			shipType: ShipTypeCV,
			battles:  battlesCV,
		},
		{
			shipType: ShipTypeBB,
			battles:  battlesBB,
		},
		{
			shipType: ShipTypeCL,
			battles:  battlesCL,
		},
		{
			shipType: ShipTypeDD,
			battles:  battlesDD,
		},
		{
			shipType: ShipTypeSS,
			battles:  battlesSS,
		},
		{
			shipType: ShipTypeNone,
			battles:  battlesNone,
		},
	}

	values := make(map[ShipType]Percentage)
	for _, eb := range eachBattles {
		t, err := NewShipType(eb.shipType)
		if err != nil {
			return ShipTypeRate{}, err
		}

		r, err := calcRate(eb.battles, battles)
		if err != nil {
			return ShipTypeRate{}, err
		}

		values[t] = r
	}

	return ShipTypeRate{values}, nil
}
