package vo

import "errors"

type TierGroupRate struct {
	values map[TierGroup]Percentage
}

func NewTierGroupRate(
	battles uint,
	battlesHigh uint,
	battlesMiddle uint,
	battlesLow uint,
) (TierGroupRate, error) {
	if battles < 1 {
		return TierGroupRate{}, errors.New("no_battles")
	}

	calcRate := func(numerator uint, denominator uint) (Percentage, error) {
		return NewPercentage(float64(numerator) / float64(denominator) * 100)
	}

	battlesNone := battles - (battlesHigh + battlesMiddle + battlesMiddle)
	eachBattles := []struct {
		tierGroup TierGroup
		battles   uint
	}{
		{
			tierGroup: ShipTierHigh,
			battles:   battlesHigh,
		},
		{
			tierGroup: ShipTierMiddle,
			battles:   battlesMiddle,
		},
		{
			tierGroup: ShipTierLow,
			battles:   battlesLow,
		},
		{
			tierGroup: ShipTierNone,
			battles:   battlesNone,
		},
	}

	values := make(map[TierGroup]Percentage)
	for _, eb := range eachBattles {
		r, err := calcRate(eb.battles, battles)
		if err != nil {
			return TierGroupRate{}, err
		}

		values[eb.tierGroup] = r
	}

	return TierGroupRate{values}, nil
}
