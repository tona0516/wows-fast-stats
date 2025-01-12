package model

import (
	"errors"
	"math"
)

var errInvalidExpectedValues = errors.New("invalid expected values for pr")

// https://asia.wows-numbers.com/personal/rating
type PR struct {
	value float64
}

func NewPR(
	averageDamage float64,
	averageFrags float64,
	winRate float64,
	expectedAverageDamage float64,
	expectedAverageFrags float64,
	expectedWinRate float64,
) (*PR, error) {
	if expectedAverageDamage == 0 || expectedAverageFrags == 0 || expectedWinRate == 0 {
		return nil, errInvalidExpectedValues
	}

	damageRatio := averageDamage / expectedAverageDamage
	fragsRatio := averageFrags / expectedAverageFrags
	winRateRatio := winRate / expectedWinRate

	normedDamage := math.Max(0, (damageRatio-0.4)/(1-0.4))
	normedfFrags := math.Max(0, (fragsRatio-0.1)/(1-0.1))
	normedWinRate := math.Max(0, (winRateRatio-0.7)/(1-0.7))

	value := 700*normedDamage + 300*normedfFrags + 150*normedWinRate

	return &PR{value: value}, nil
}

func (p *PR) Value() float64 {
	return p.value
}
