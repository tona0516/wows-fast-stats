package data

import "math"

type PRFactor struct {
	damage float64
	frags  float64
	wins   float64
}

func (rs *PRFactor) Valid() bool {
	// All values are not nan or inf.
	return !(math.IsNaN(rs.damage) || math.IsInf(rs.damage, 1) ||
		math.IsNaN(rs.frags) || math.IsInf(rs.frags, 1) ||
		math.IsNaN(rs.wins) || math.IsInf(rs.wins, 1))
}
