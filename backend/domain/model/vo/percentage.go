package vo

import (
	"errors"
	"strconv"
)

type Percentage struct {
	ValueObject[float64]
}

func NewPercentage(value float64) (Percentage, error) {
	if value < 0 || value > 1 {
		return Percentage{}, errors.New("out_of_percentage:" + strconv.FormatFloat(value, 'f', -1, 64))
	}
	return Percentage{ValueObject[float64]{value}}, nil
}
