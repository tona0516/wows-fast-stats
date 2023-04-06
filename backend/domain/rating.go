package domain

import (
	"math"
)

type Rating struct{}

func (s *Rating) CombatPower(
	avgDamage float64,
	kdRate float64,
	avgExp float64,
	tier int,
	shipType string,
) int {
	shipTypeCoef := 1.0
	if shipType == "Battleship" {
		shipTypeCoef = 0.7
	}
	if shipType == "AirCarrier" {
		shipTypeCoef = 0.5
	}

	combatPower := ((avgDamage * kdRate * avgExp) / 800.0) * (1.0 - 0.03*float64(tier)) * shipTypeCoef

	return int(combatPower)
}

func (s *Rating) PersonalRating(
	actualDamage float64,
	actualFrags float64,
	actualWins float64,
	expectedDamage float64,
	expectedFrags float64,
	expectedWins float64,
) int {
	damageRatio := actualDamage / expectedDamage
	fragsRatio := actualFrags / expectedFrags
	winsRatio := actualWins / expectedWins

	damageNorm := math.Max(0, (damageRatio-0.4)/(1-0.4))
	fragsNorm := math.Max(0, (fragsRatio-0.1)/(1-0.1))
	winsNorm := math.Max(0, (winsRatio-0.7)/(1-0.7))

	personalRating := 700*damageNorm + 300*fragsNorm + 150*winsNorm

	return int(personalRating)
}
