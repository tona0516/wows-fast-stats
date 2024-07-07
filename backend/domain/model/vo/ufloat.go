package vo

import (
	"errors"
	"strconv"
)

type UFloat struct {
	ValueObject[float64]
}

func NewUFloat(value float64) (UFloat, error) {
	if value < 0 {
		return UFloat{}, errors.New("negative_value:" + strconv.FormatFloat(value, 'f', -1, 64))
	}
    return UFloat{ValueObject[float64]{value}}, nil
}
