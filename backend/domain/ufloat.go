package domain

import (
	"errors"
	"strconv"
)

type ufloat struct {
	ValueObject[float64]
}

func NewUFloat(value float64) (ufloat, error) {
	if value < 0 {
		return ufloat{}, errors.New("negative_value:" + strconv.FormatFloat(value, 'f', -1, 64))
	}
    return ufloat{ValueObject[float64]{value}}, nil
}
