package domain

import (
	"math"
)

type RatingFactor struct {
	AvgDamage float64
	AvgFrags  float64
	WinRate   float64
}

func (rs *RatingFactor) Valid() bool {
	// All values are not nan or inf.
	return !(math.IsNaN(rs.AvgDamage) || math.IsInf(rs.AvgDamage, 1) ||
		math.IsNaN(rs.AvgFrags) || math.IsInf(rs.AvgFrags, 1) ||
		math.IsNaN(rs.WinRate) || math.IsInf(rs.WinRate, 1))
}

type Rating struct{}

func (r *Rating) PersonalRating(
	actual RatingFactor,
	expected RatingFactor,
) float64 {
	ratio := RatingFactor{
		AvgDamage: actual.AvgDamage / expected.AvgDamage,
		AvgFrags:  actual.AvgFrags / expected.AvgFrags,
		WinRate:   actual.WinRate / expected.WinRate,
	}

	if !ratio.Valid() {
		return -1
	}

	norm := RatingFactor{
		AvgDamage: math.Max(0, (ratio.AvgDamage-0.4)/(1-0.4)),
		AvgFrags:  math.Max(0, (ratio.AvgFrags-0.1)/(1-0.1)),
		WinRate:   math.Max(0, (ratio.WinRate-0.7)/(1-0.7)),
	}

	personalRating := 700*norm.AvgDamage + 300*norm.AvgFrags + 150*norm.WinRate

	return personalRating
}
