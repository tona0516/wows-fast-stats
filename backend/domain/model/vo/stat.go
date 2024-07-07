package vo

type Stat struct {
	ship struct {
		battles      uint
		damage       UFloat
		maxDamage    UFloat
		winRate      Percentage
		survivedRate SurvivedRate
		kdRate       UFloat
		avgKill      UFloat
		avgExp       UFloat
		pr           UFloat
		hitRate      HitRate
		planesKilled UFloat
		platoonRate  UFloat
	}
	overall struct {
		battles       uint
		damage        UFloat
		maxDamage     MaxDamage
		winRate       Percentage
		survivedRate  SurvivedRate
		kdRate        UFloat
		avgKill       UFloat
		avgExp        UFloat
		pr            UFloat
		threatLevel   ThreatLevel
		avgTier       UFloat
		shipTypeRate  ShipTypeRate
		tierGroupRate TierGroupRate
		platoonRate   UFloat
	}
}

func NewStat() Stat {
	return Stat{}
}
