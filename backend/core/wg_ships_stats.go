package core

type WGShipsStats struct {
	WGResponseCommon[map[AccountID][]WGShipsStatsData]
}

type WGShipsStatsData struct {
	Pvp     WGShipStatsValues `json:"pvp"`
	PvpSolo WGShipStatsValues `json:"pvp_solo"`
	PvpDiv2 struct {
		Battles uint `json:"battles"`
	} `json:"pvp_div2"`
	PvpDiv3 struct {
		Battles uint `json:"battles"`
	} `json:"pvp_div3"`
	RankSolo WGShipStatsValues `json:"rank_solo"`
	ShipID   ShipID            `json:"ship_id"`
}

func (sd WGShipsStatsData) shipStatsValues(pattern StatsPattern) WGShipStatsValues {
	switch pattern {
	case StatsPatternPvPAll:
		return sd.Pvp
	case StatsPatternPvPSolo:
		return sd.PvpSolo
	case StatsPatternRankSolo:
		return sd.RankSolo
	}

	return WGShipStatsValues{}
}

func (sd WGShipsStatsData) platoonRate() float64 {
	allBattles := sd.Pvp.Battles
	soloRate := safeDivide(sd.PvpSolo.Battles, allBattles) * 1
	div2Rate := safeDivide(sd.PvpDiv2.Battles, allBattles) * 2
	div3Rate := safeDivide(sd.PvpDiv3.Battles, allBattles) * 3

	return soloRate + div2Rate + div3Rate
}

type WGShipStatsValues struct {
	Wins            uint       `json:"wins"`
	Battles         uint       `json:"battles"`
	DamageDealt     uint       `json:"damage_dealt"`
	MaxDamageDealt  uint       `json:"max_damage_dealt"`
	Frags           uint       `json:"frags"`
	SurvivedWins    uint       `json:"survived_wins"`
	SurvivedBattles uint       `json:"survived_battles"`
	Xp              uint       `json:"xp"`
	MainBattery     WGArmament `json:"main_battery"`
	Torpedoes       WGArmament `json:"torpedoes"`
	PlanesKilled    uint       `json:"planes_killed"`
}

func (v WGShipStatsValues) avgDamage() float64 {
	return safeDivide(v.DamageDealt, v.Battles)
}

func (v WGShipStatsValues) avgKills() float64 {
	return safeDivide(v.Frags, v.Battles)
}

func (v WGShipStatsValues) winRate() float64 {
	return safeDivide(v.Wins, v.Battles) * 100
}

func (v WGShipStatsValues) kdRate() float64 {
	deaths := max(v.Battles-v.SurvivedBattles, 1)
	return safeDivide(v.Frags, deaths)
}

func (v WGShipStatsValues) exp() float64 {
	return safeDivide(v.Xp, v.Battles)
}

func (v WGShipStatsValues) survivedRate() SurvivedRate {
	survivedLoses := v.SurvivedBattles - v.SurvivedWins
	loses := v.Battles - v.Wins

	return SurvivedRate{
		All:  safeDivide(v.SurvivedBattles, v.Battles) * 100,
		Win:  safeDivide(v.SurvivedWins, v.Wins) * 100,
		Lose: safeDivide(survivedLoses, loses) * 100,
	}
}

func (v WGShipStatsValues) avgPlanesKilled() float64 {
	return safeDivide(v.PlanesKilled, v.Battles)
}

func (v WGShipStatsValues) hitRate() HitRate {
	mb := v.MainBattery.hitRate()
	torps := v.Torpedoes.hitRate()

	return HitRate{
		MainBattery: mb,
		Torpedoes:   torps,
	}
}

type WGArmament struct {
	Hits  uint `json:"hits"`
	Shots uint `json:"shots"`
}

func (a WGArmament) hitRate() float64 {
	return safeDivide(a.Hits, a.Shots) * 100
}

type AllPlayerShipStats map[AccountID]PlayerShipStats

type PlayerShipStats map[ShipID]WGShipsStatsData

func (pss PlayerShipStats) avgTier(
	pattern StatsPattern,
	warships Warships,
) float64 {
	var (
		sum        uint
		allBattles uint
	)

	for shipID, stats := range pss {
		warship, ok := warships[shipID]
		if !ok {
			continue
		}

		values := stats.shipStatsValues(pattern)
		sum += values.Battles * warship.Tier
		allBattles += values.Battles
	}

	return safeDivide(sum, allBattles)
}

func (pss PlayerShipStats) usingShipTypeRate(
	pattern StatsPattern,
	warships Warships,
) ShipTypeGroup {
	shipTypeMap := make(map[ShipType]uint)
	allBattles := uint(0)

	for shipID, stats := range pss {
		warship, ok := warships[shipID]
		if !ok {
			continue
		}

		values := stats.shipStatsValues(pattern)
		shipTypeMap[warship.Type] += values.Battles
		allBattles += values.Battles
	}

	return ShipTypeGroup{
		SS: safeDivide(shipTypeMap[ShipTypeSS], allBattles) * 100,
		DD: safeDivide(shipTypeMap[ShipTypeDD], allBattles) * 100,
		CL: safeDivide(shipTypeMap[ShipTypeCL], allBattles) * 100,
		BB: safeDivide(shipTypeMap[ShipTypeBB], allBattles) * 100,
		CV: safeDivide(shipTypeMap[ShipTypeCV], allBattles) * 100,
	}
}

func (pss PlayerShipStats) usingTierRate(
	pattern StatsPattern,
	warships Warships,
) TierGroup {
	tierGroupMap := make(map[string]uint)
	allBattles := uint(0)

	for shipID, stats := range pss {
		warship, ok := warships[shipID]
		if !ok {
			continue
		}

		values := stats.shipStatsValues(pattern)
		tier := warship.Tier
		battles := values.Battles
		switch {
		case tier >= 1 && tier <= 4:
			tierGroupMap["low"] += battles
		case tier >= 5 && tier <= 7:
			tierGroupMap["middle"] += battles
		case tier >= 8:
			tierGroupMap["high"] += battles
		}
		allBattles += battles
	}

	return TierGroup{
		Low:    safeDivide(tierGroupMap["low"], allBattles) * 100,
		Middle: safeDivide(tierGroupMap["middle"], allBattles) * 100,
		High:   safeDivide(tierGroupMap["high"], allBattles) * 100,
	}
}
