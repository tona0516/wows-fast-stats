package domain

import (
	"math"
)

type Rating struct{}

func (s *Rating) CombatPower(
	avgDamage float64,
	kdRate float64,
	avgExp float64,
	tier uint,
	shipType string,
) float64 {
	shipTypeCoef := 1.0
	if shipType == "Battleship" {
		shipTypeCoef = 0.7
	}
	if shipType == "AirCarrier" {
		shipTypeCoef = 0.5
	}

	combatPower := ((avgDamage * kdRate * avgExp) / 800.0) * (1.0 - 0.03*float64(tier)) * shipTypeCoef

	return combatPower
}

func (s *Rating) PersonalRating(
	actualDamage float64,
	actualFrags float64,
	actualWins float64,
	expectedDamage float64,
	expectedFrags float64,
	expectedWins float64,
) float64 {
	damageRatio := actualDamage / expectedDamage
    if math.IsNaN(damageRatio) {
        damageRatio = 0
    }
	fragsRatio := actualFrags / expectedFrags
    if math.IsNaN(fragsRatio) {
        fragsRatio = 0
    }
	winsRatio := actualWins / expectedWins
    if math.IsNaN(winsRatio) {
        winsRatio = 0
    }

	damageNorm := math.Max(0, (damageRatio-0.4)/(1-0.4))
	fragsNorm := math.Max(0, (fragsRatio-0.1)/(1-0.1))
	winsNorm := math.Max(0, (winsRatio-0.7)/(1-0.7))

	personalRating := 700*damageNorm + 300*fragsNorm + 150*winsNorm

	return personalRating
}
