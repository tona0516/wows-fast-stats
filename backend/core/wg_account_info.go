package core

import "math"

type WGAccountInfo struct {
	WGResponseCommon[map[AccountID]WGAccountInfoData]
}

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp      WGPlayerStatsValues `json:"pvp"`
		PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
		PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
		PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
		RankSolo WGPlayerStatsValues `json:"rank_solo"`
	} `json:"statistics"`
}

func (ad WGAccountInfoData) playerStatsValue(pattern StatsPattern) WGPlayerStatsValues {
	switch pattern {
	case StatsPatternPvPAll:
		return ad.Statistics.Pvp
	case StatsPatternPvPSolo:
		return ad.Statistics.PvpSolo
	case StatsPatternRankSolo:
		return ad.Statistics.RankSolo
	}

	return WGPlayerStatsValues{}
}

func (ad WGAccountInfoData) platoonRate() float64 {
	allBattles := ad.Statistics.Pvp.Battles
	soloRate := safeDivide(ad.Statistics.PvpSolo.Battles, allBattles) * 1
	div2Rate := safeDivide(ad.Statistics.PvpDiv2.Battles, allBattles) * 2
	div3Rate := safeDivide(ad.Statistics.PvpDiv3.Battles, allBattles) * 3

	return soloRate + div2Rate + div3Rate
}

type WGPlayerStatsValues struct {
	Wins                 uint   `json:"wins"`
	Battles              uint   `json:"battles"`
	DamageDealt          uint   `json:"damage_dealt"`
	MaxDamageDealt       uint   `json:"max_damage_dealt"`
	MaxDamageDealtShipID ShipID `json:"max_damage_dealt_ship_id"`
	Frags                uint   `json:"frags"`
	SurvivedWins         uint   `json:"survived_wins"`
	SurvivedBattles      uint   `json:"survived_battles"`
	Xp                   uint   `json:"xp"`
}

func (v WGPlayerStatsValues) avgDamage() float64 {
	return safeDivide(v.DamageDealt, v.Battles)
}

func (v WGPlayerStatsValues) avgKills() float64 {
	return safeDivide(v.Frags, v.Battles)
}

func (v WGPlayerStatsValues) winRate() float64 {
	return safeDivide(v.Wins, v.Battles) * 100
}

func (v WGPlayerStatsValues) kdRate() float64 {
	death := math.Max(float64(v.Battles-v.SurvivedBattles), 1)
	return safeDivide(v.Frags, death)
}

func (v WGPlayerStatsValues) exp() float64 {
	return safeDivide(v.Xp, v.Battles)
}

func (v WGPlayerStatsValues) survivedRate() SurvivedRate {
	survivedLoses := v.SurvivedBattles - v.SurvivedWins
	loses := v.Battles - v.Wins

	return SurvivedRate{
		All:  safeDivide(v.SurvivedBattles, v.Battles) * 100,
		Win:  safeDivide(v.SurvivedWins, v.Wins) * 100,
		Lose: safeDivide(survivedLoses, loses) * 100,
	}
}
