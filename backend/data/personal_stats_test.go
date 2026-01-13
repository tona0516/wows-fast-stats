package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	allowableDelta = 0.01
)

//nolint:gochecknoglobals
var (
	emptyAccountInfo = WGAccountInfoData{}
	emptyShipStats   = PlayerShipStats{}
	emptyWarships    = Warships{}
	emptyShipBadges  = PlayerShipBadges{}
)

// ref: https://asia.wows-numbers.com/personal/rating
func TestPersonalStats_PR_Ship(t *testing.T) {
	t.Parallel()

	params := []struct {
		battles  uint
		expected RatingValue
	}{
		{100, NewRatingValue(1875, RatingGreat)},
		{0, NewRatingValue(-1, RatingNone)},
	}

	useShipID := ShipID(0)
	for _, p := range params {
		stats := NewPersonalStats(
			useShipID,
			emptyAccountInfo,
			PlayerShipStats{
				useShipID: {
					Pvp: WGShipStatsValues{
						Battles:     p.battles,
						DamageDealt: 1000000,
						Frags:       200,
						Wins:        60,
					},
					ShipID: useShipID,
				},
			},
			emptyShipBadges,
			Warships{
				useShipID: {
					ServerAverage: ServerAverage{
						Damage:  8000,
						Frags:   1,
						WinRate: 50,
					},
				},
			},
		)

		actual := stats.PR(StatsCategoryShip, StatsPatternPvPAll)
		assert.InDelta(t, p.expected.Value, actual.Value, allowableDelta)
		assert.Equal(t, p.expected.Rating, actual.Rating, allowableDelta)
	}
}

// ref: https://asia.wows-numbers.com/personal/rating
func TestPersonalStats_PR_Overall(t *testing.T) {
	t.Parallel()

	useShipID := ShipID(0)
	stats := NewPersonalStats(
		useShipID,
		emptyAccountInfo,
		PlayerShipStats{
			1: {
				Pvp: WGShipStatsValues{
					Battles:     2,
					DamageDealt: 54468,
					Wins:        1,
					Frags:       1,
				},
				ShipID: 1,
			},
			2: {
				Pvp: WGShipStatsValues{
					Battles:     1,
					DamageDealt: 155185,
					Wins:        1,
					Frags:       1,
				},
				ShipID: 2,
			},
			3: {
				Pvp: WGShipStatsValues{
					Battles:     1,
					DamageDealt: 51576,
					Wins:        1,
					Frags:       2,
				},
				ShipID: 3,
			},
			4: {
				Pvp: WGShipStatsValues{
					Battles:     1,
					DamageDealt: 117285,
					Wins:        1,
					Frags:       2,
				},
				ShipID: 4,
			},
		},
		emptyShipBadges,
		Warships{
			1: {ServerAverage: ServerAverage{Damage: 53792.23172971, Frags: 0.6935181784796, WinRate: 50.092406353286}},
			2: {ServerAverage: ServerAverage{Damage: 46228.419395466, Frags: 0.80128883291351, WinRate: 51.202824307302}},
			3: {ServerAverage: ServerAverage{Damage: 25864.417248367, Frags: 0.69715604593558, WinRate: 51.11762215717}},
			4: {ServerAverage: ServerAverage{Damage: 77931.580907796, Frags: 0.68628943618969, WinRate: 50.386342357012}},
		},
	)

	actual := stats.PR(StatsCategoryOverall, StatsPatternPvPAll)
	assert.InDelta(t, 2215.0243612353, actual.Value, allowableDelta)
	assert.Equal(t, RatingUnicum, actual.Rating)
}

func TestPersonalStats_AvgDamage_Overall(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	actual := stats.AvgDamage(StatsCategoryOverall, StatsPatternPvPAll)
	assert.InDelta(t, 10000, actual.Value, allowableDelta)
	assert.Equal(t, RatingNone, actual.Rating)
}

func TestPersonalStats_AvgDamage_Overall_Solo(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				PvpSolo: WGPlayerStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	actual := stats.AvgDamage(StatsCategoryOverall, StatsPatternPvPSolo)
	assert.InDelta(t, 10000, actual.Value, allowableDelta)
	assert.Equal(t, RatingNone, actual.Rating)
}

func TestPersonalStats_AvgDamage_Overall_Rank(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				RankSolo: WGPlayerStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	actual := stats.AvgDamage(StatsCategoryOverall, StatsPatternRankSolo)
	assert.InDelta(t, 10000, actual.Value, allowableDelta)
	assert.Equal(t, RatingNone, actual.Rating)
}

func TestPersonalStats_AvgDamage_Ship(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		1,
		emptyAccountInfo,
		PlayerShipStats{
			1: {
				Pvp: WGShipStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
				ShipID: 1,
			},
		},
		emptyShipBadges,
		Warships{
			1: {
				ServerAverage: ServerAverage{
					Damage: 12000,
				},
			},
		},
	)

	actual := stats.AvgDamage(StatsCategoryShip, StatsPatternPvPAll)
	assert.InDelta(t, 10000, actual.Value, allowableDelta)
	assert.Equal(t, RatingAvg, actual.Rating)
}

func TestPersonalStats_AvgDamage_Ship_Solo(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		1,
		emptyAccountInfo,
		PlayerShipStats{
			1: {
				PvpSolo: WGShipStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
				ShipID: 1,
			},
		},
		emptyShipBadges,
		Warships{
			1: {
				ServerAverage: ServerAverage{
					Damage: 12000,
				},
			},
		},
	)

	actual := stats.AvgDamage(StatsCategoryShip, StatsPatternPvPSolo)
	assert.InDelta(t, 10000, actual.Value, allowableDelta)
	assert.Equal(t, RatingAvg, actual.Rating)
}

func TestPersonalStats_AvgDamage_Ship_Rank(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		1,
		emptyAccountInfo,
		PlayerShipStats{
			1: {
				RankSolo: WGShipStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
				ShipID: 1,
			},
		},
		emptyShipBadges,
		Warships{
			1: {
				ServerAverage: ServerAverage{
					Damage: 12000,
				},
			},
		},
	)

	actual := stats.AvgDamage(StatsCategoryShip, StatsPatternRankSolo)
	assert.InDelta(t, 10000, actual.Value, allowableDelta)
	assert.Equal(t, RatingAvg, actual.Rating)
}

func TestPersonalStats_MaxDamage_Ship(t *testing.T) {
	t.Parallel()

	expected := MaxDamage{
		Value: 200000,
	}

	useShipID := ShipID(100)
	stats := NewPersonalStats(
		useShipID,
		emptyAccountInfo,
		PlayerShipStats{
			useShipID: {
				Pvp: WGShipStatsValues{
					MaxDamageDealt: expected.Value,
				},
				ShipID: useShipID,
			},
		},
		emptyShipBadges,
		emptyWarships,
	)

	assert.Equal(t, expected, stats.MaxDamage(StatsCategoryShip, StatsPatternPvPAll))
}

func TestPersonalStats_MaxDamage_Overall(t *testing.T) {
	t.Parallel()

	expected := MaxDamage{
		ShipID:   100,
		ShipName: "yamato",
		ShipTier: 10,
		Value:    200000,
	}

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					MaxDamageDealt:       expected.Value,
					MaxDamageDealtShipID: expected.ShipID,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		Warships{
			expected.ShipID: {
				Name: expected.ShipName,
				Tier: expected.ShipTier,
			},
		},
	)

	assert.Equal(t, expected, stats.MaxDamage(StatsCategoryOverall, StatsPatternPvPAll))
}

func TestPersonalStats_Battles(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	assert.InDelta(t, 100, stats.Battles(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestPersonalStats_KdRate(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles:         100,
					SurvivedBattles: 60,
					Frags:           20,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	assert.InDelta(t, 0.5, stats.KdRate(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestPersonalStats_AvgKill(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
					Frags:   30,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	assert.InDelta(t, 0.3, stats.AvgKill(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestPersonalStats_AvgExp(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
					Xp:      150000,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	assert.InDelta(t, 1500, stats.AvgExp(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestPersonalStats_WinRate(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		1,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
					Wins:    60,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	actual := stats.WinRate(StatsCategoryOverall, StatsPatternPvPAll)
	assert.InDelta(t, 60, actual.Value, allowableDelta)
	assert.Equal(t, RatingUnicum, actual.Rating)
}

func TestPersonalStats_SurvivedRate(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles:         100,
					SurvivedBattles: 40,
					Wins:            60,
					SurvivedWins:    20,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	actual := stats.SurvivedRate(StatsCategoryOverall, StatsPatternPvPAll)
	assert.InDelta(t, 40, actual.All, allowableDelta)
	assert.InDelta(t, 33.33, actual.Win, allowableDelta)
	assert.InDelta(t, 50, actual.Lose, allowableDelta)
}

func TestPersonalStats_HitRate(t *testing.T) {
	t.Parallel()

	useShipID := ShipID(0)
	stats := NewPersonalStats(
		useShipID,
		emptyAccountInfo,
		PlayerShipStats{
			useShipID: {
				Pvp: WGShipStatsValues{
					MainBattery: struct {
						Hits  uint `json:"hits"`
						Shots uint `json:"shots"`
					}{
						Hits:  100,
						Shots: 200,
					},
					Torpedoes: struct {
						Hits  uint `json:"hits"`
						Shots uint `json:"shots"`
					}{
						Hits:  10,
						Shots: 40,
					},
				},
				ShipID: useShipID,
			},
		},
		emptyShipBadges,
		emptyWarships,
	)

	actual := stats.HitRate(StatsPatternPvPAll)
	assert.InDelta(t, 50, actual.MainBattery, allowableDelta)
	assert.InDelta(t, 25, actual.Torpedoes, allowableDelta)
}

func TestPersonalStats_PlanesKilled(t *testing.T) {
	t.Parallel()

	useShipID := ShipID(0)
	stats := NewPersonalStats(
		useShipID,
		emptyAccountInfo,
		PlayerShipStats{
			useShipID: {
				Pvp: WGShipStatsValues{
					PlanesKilled: 334,
					Battles:      10,
				},
				ShipID: useShipID,
			},
		},
		emptyShipBadges,
		emptyWarships,
	)

	assert.InDelta(t, 33.4, stats.PlanesKilled(StatsPatternPvPAll), allowableDelta)
}

func TestPersonalStats_AvgTier(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		emptyAccountInfo,
		PlayerShipStats{
			100: {
				Pvp:    WGShipStatsValues{Battles: 20},
				ShipID: 100,
			},
			200: {
				Pvp:    WGShipStatsValues{Battles: 50},
				ShipID: 200,
			},
		},
		emptyShipBadges,
		Warships{
			100: {Tier: 5},
			200: {Tier: 8},
		},
	)

	assert.InDelta(t, 7.14, stats.AvgTier(StatsPatternPvPAll), allowableDelta)
}

func TestPersonalStats_UsingTierRate(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		emptyAccountInfo,
		PlayerShipStats{
			100: {
				Pvp:    WGShipStatsValues{Battles: 30},
				ShipID: 100,
			},
			200: {
				Pvp:    WGShipStatsValues{Battles: 50},
				ShipID: 200,
			},
			300: {
				Pvp:    WGShipStatsValues{Battles: 20},
				ShipID: 300,
			},
		},
		emptyShipBadges,
		Warships{
			100: {Tier: 5},
			200: {Tier: 8},
			300: {Tier: 4},
		},
	)

	tierGroup := stats.UsingTierRate(StatsPatternPvPAll)
	assert.InDelta(t, 20, tierGroup.Low, allowableDelta)
	assert.InDelta(t, 30, tierGroup.Middle, allowableDelta)
	assert.InDelta(t, 50, tierGroup.High, allowableDelta)
}

func TestPersonalStats_UsingShipTypeRate(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		emptyAccountInfo,
		PlayerShipStats{
			100: {
				Pvp:    WGShipStatsValues{Battles: 30},
				ShipID: 100,
			},
			200: {
				Pvp:    WGShipStatsValues{Battles: 50},
				ShipID: 200,
			},
			300: {
				Pvp:    WGShipStatsValues{Battles: 20},
				ShipID: 300,
			},
			400: {
				Pvp:    WGShipStatsValues{Battles: 100},
				ShipID: 400,
			},
		},
		emptyShipBadges,
		Warships{
			100: {Type: ShipTypeDD},
			200: {Type: ShipTypeCL},
			300: {Type: ShipTypeBB},
			400: {Type: ShipTypeCV},
		},
	)

	shipTypeGroup := stats.UsingShipTypeRate(StatsPatternPvPAll)
	assert.InDelta(t, 0, shipTypeGroup.SS, allowableDelta)
	assert.InDelta(t, 15, shipTypeGroup.DD, allowableDelta)
	assert.InDelta(t, 25, shipTypeGroup.CL, allowableDelta)
	assert.InDelta(t, 10, shipTypeGroup.BB, allowableDelta)
	assert.InDelta(t, 50, shipTypeGroup.CV, allowableDelta)
}

func TestPersonalStats_PlatoonRate(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp      WGPlayerStatsValues `json:"pvp"`
				PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
				PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
				PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
				RankSolo WGPlayerStatsValues `json:"rank_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
				},
				PvpSolo: WGPlayerStatsValues{
					Battles: 10,
				},
				PvpDiv2: WGPlayerStatsValues{
					Battles: 40,
				},
				PvpDiv3: WGPlayerStatsValues{
					Battles: 50,
				},
			},
		},
		emptyShipStats,
		emptyShipBadges,
		emptyWarships,
	)

	assert.InDelta(t, 2.4, stats.PlatoonRate(StatsCategoryOverall), allowableDelta)
}

func TestPersonalStats_EfficiencyBadge(t *testing.T) {
	t.Parallel()

	useShipID := ShipID(0)
	stats := NewPersonalStats(
		useShipID,
		emptyAccountInfo,
		PlayerShipStats{
			useShipID: {
				Pvp: WGShipStatsValues{
					Battles: 100,
				},
				ShipID: useShipID,
			},
		},
		PlayerShipBadges{
			useShipID: {
				ShipID:        useShipID,
				TopGradeClass: 1,
			},
		},
		emptyWarships,
	)

	assert.Equal(t, EfficiencyBadge("E"), stats.EfficiencyBadge())
}

func TestPersonalStats_EfficiencyBadges(t *testing.T) {
	t.Parallel()

	stats := NewPersonalStats(
		0,
		emptyAccountInfo,
		emptyShipStats,
		PlayerShipBadges{
			100: {
				ShipID:        100,
				TopGradeClass: 1,
			},
			200: {
				ShipID:        200,
				TopGradeClass: 2,
			},
			300: {
				ShipID:        300,
				TopGradeClass: 2,
			},
			400: {
				ShipID:        400,
				TopGradeClass: 3,
			},
			500: {
				ShipID:        500,
				TopGradeClass: 4,
			},
		},
		emptyWarships,
	)

	badges := stats.EfficiencyBadges()
	assert.Equal(t, 1, badges.Expert)
	assert.Equal(t, 2, badges.First)
	assert.Equal(t, 1, badges.Second)
	assert.Equal(t, 1, badges.Third)
}
