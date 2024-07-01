package domain

type ShipStat struct {
	battles         uint
	avgDamage       ufloat
	maxDamage       uint
	winRate         percentage
	survivedRate    SurvivedRate
	kdRate          ufloat
	avgKill         ufloat
	avgExp          ufloat
	pr              float64
	hitRate         HitRate
	avgPlanesKilled ufloat
	platoonRate     float64
}
