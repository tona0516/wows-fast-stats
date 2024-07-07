package vo

type ShipStat struct {
	battles         uint
	avgDamage       UFloat
	maxDamage       uint
	winRate         Percentage
	survivedRate    SurvivedRate
	kdRate          UFloat
	avgKill         UFloat
	avgExp          UFloat
	pr              float64
	hitRate         HitRate
	avgPlanesKilled UFloat
	platoonRate     float64
}
