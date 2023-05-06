package domain

import (
	"math"
)

type RatingFactor struct {
	Damage float64
	Frags  float64
	Wins   float64
}

func (rs *RatingFactor) Valid() bool {
	// All values are not nan or inf.
	return !(math.IsNaN(rs.Damage) || math.IsInf(rs.Damage, 1) ||
		math.IsNaN(rs.Frags) || math.IsInf(rs.Frags, 1) ||
		math.IsNaN(rs.Wins) || math.IsInf(rs.Wins, 1))
}

type Rating struct{}

func (r *Rating) PersonalRating(
	actual RatingFactor,
	expected RatingFactor,
) float64 {
	ratio := RatingFactor{
		Damage: actual.Damage / expected.Damage,
		Frags:  actual.Frags / expected.Frags,
		Wins:   actual.Wins / expected.Wins,
	}

	if !ratio.Valid() {
		return -1
	}

	norm := RatingFactor{
		Damage: math.Max(0, (ratio.Damage-0.4)/(1-0.4)),
		Frags:  math.Max(0, (ratio.Frags-0.1)/(1-0.1)),
		Wins:   math.Max(0, (ratio.Wins-0.7)/(1-0.7)),
	}

	personalRating := 700*norm.Damage + 300*norm.Frags + 150*norm.Wins

	return personalRating
}
