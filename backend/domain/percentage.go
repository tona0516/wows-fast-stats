package domain

import (
	"errors"
	"strconv"
)

type percentage struct {
	ValueObject[float64]
}

func NewPercentage(value float64) (percentage, error) {
	if value < 0 || value > 1 {
		return percentage{}, errors.New("out_of_percentage:" + strconv.FormatFloat(value, 'f', -1, 64))
	}
	return percentage{ValueObject[float64]{value}}, nil
}
