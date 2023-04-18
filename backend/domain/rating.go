package domain

import (
	"math"
)

type Rating struct{}

func (s *Rating) PersonalRating(
	actualDamage float64,
	actualFrags float64,
	actualWins float64,
	expectedDamage float64,
	expectedFrags float64,
	expectedWins float64,
) float64 {
	damageRatio := actualDamage / expectedDamage
    if math.IsNaN(damageRatio) || math.IsInf(damageRatio, 1) {
        damageRatio = 0
    }
	fragsRatio := actualFrags / expectedFrags
    if math.IsNaN(fragsRatio) || math.IsInf(fragsRatio, 1){
        fragsRatio = 0
    }
	winsRatio := actualWins / expectedWins
    if math.IsNaN(winsRatio) || math.IsInf(winsRatio, 1){
        winsRatio = 0
    }

	damageNorm := math.Max(0, (damageRatio-0.4)/(1-0.4))
	fragsNorm := math.Max(0, (fragsRatio-0.1)/(1-0.1))
	winsNorm := math.Max(0, (winsRatio-0.7)/(1-0.7))

	personalRating := 700*damageNorm + 300*fragsNorm + 150*winsNorm

	return personalRating
}
