package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	allowableDelta = 0.01
)

//nolint:gochecknoglobals
var (
	emptyAccountInfo   = WGAccountInfoData{}
	emptyShipsStats    = []WGShipsStatsData{}
	emptyExpectedStats = AllExpectedStats{}
	emptyWarships      = Warships{}
)

// ref: https://asia.wows-numbers.com/personal/rating
func TestStats_PR_Ship(t *testing.T) {
	t.Parallel()

	params := []struct {
		battles  uint
		expected float64
	}{
		{100, 1875},
		{0, -1},
	}

	useShipID := 0
	for _, p := range params {
		stats := NewStats(
			useShipID,
			emptyAccountInfo,
			[]WGShipsStatsData{
				{
					Pvp: WGShipStatsValues{
						Battles:     p.battles,
						DamageDealt: 1000000,
						Frags:       200,
						Wins:        60,
					},
					ShipID: useShipID,
				},
			},
			AllExpectedStats{
				useShipID: {
					AverageDamageDealt: 8000,
					AverageFrags:       1,
					WinRate:            50,
				},
			},
			emptyWarships,
		)

		assert.InDelta(t, p.expected, stats.PR(StatsCategoryShip, StatsPatternPvPAll), allowableDelta)
	}
}

// ref: https://asia.wows-numbers.com/personal/rating
func TestStats_PR_Overall(t *testing.T) {
	t.Parallel()

	useShipID := 0
	stats := NewStats(
		useShipID,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp: WGShipStatsValues{
					Battles:     2,
					DamageDealt: 54468,
					Wins:        1,
					Frags:       1,
				},
				ShipID: 1,
			},
			{
				Pvp: WGShipStatsValues{
					Battles:     1,
					DamageDealt: 155185,
					Wins:        1,
					Frags:       1,
				},
				ShipID: 2,
			},
			{
				Pvp: WGShipStatsValues{
					Battles:     1,
					DamageDealt: 51576,
					Wins:        1,
					Frags:       2,
				},
				ShipID: 3,
			},
			{
				Pvp: WGShipStatsValues{
					Battles:     1,
					DamageDealt: 117285,
					Wins:        1,
					Frags:       2,
				},
				ShipID: 4,
			},
		},
		AllExpectedStats{
			1: {
				AverageDamageDealt: 53792.23172971,
				WinRate:            50.092406353286,
				AverageFrags:       0.6935181784796,
			},
			2: {
				AverageDamageDealt: 46228.419395466,
				WinRate:            51.202824307302,
				AverageFrags:       0.80128883291351,
			},
			3: {
				AverageDamageDealt: 25864.417248367,
				WinRate:            51.11762215717,
				AverageFrags:       0.69715604593558,
			},
			4: {
				AverageDamageDealt: 77931.580907796,
				WinRate:            50.386342357012,
				AverageFrags:       0.68628943618969,
			},
		},
		emptyWarships,
	)

	assert.InDelta(t, 2215.0243612353, stats.PR(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_AvgDamage_Overall(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 10000, stats.AvgDamage(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_AvgDamage_Overall_Solo(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				PvpSolo: WGPlayerStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 10000, stats.AvgDamage(StatsCategoryOverall, StatsPatternPvPSolo), allowableDelta)
}

func TestStats_AvgDamage_Ship(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp: WGShipStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 10000, stats.AvgDamage(StatsCategoryShip, StatsPatternPvPAll), allowableDelta)
}

func TestStats_AvgDamage_Ship_Solo(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				PvpSolo: WGShipStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 10000, stats.AvgDamage(StatsCategoryShip, StatsPatternPvPSolo), allowableDelta)
}

func TestStats_MaxDamage_Ship(t *testing.T) {
	t.Parallel()

	expected := MaxDamage{
		Damage: 200000,
	}

	useShipID := 100
	stats := NewStats(
		useShipID,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp: WGShipStatsValues{
					MaxDamageDealt: expected.Damage,
				},
				ShipID: useShipID,
			},
		},
		emptyExpectedStats,
		emptyWarships,
	)

	assert.Equal(t, expected, stats.MaxDamage(StatsCategoryShip, StatsPatternPvPAll))
}

func TestStats_MaxDamage_Overall(t *testing.T) {
	t.Parallel()

	expected := MaxDamage{
		ShipID:   100,
		ShipName: "yamato",
		ShipTier: 10,
		Damage:   200000,
	}

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					MaxDamageDealt:       expected.Damage,
					MaxDamageDealtShipID: expected.ShipID,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		Warships{
			expected.ShipID: {
				Name: expected.ShipName,
				Tier: expected.ShipTier,
			},
		},
	)

	assert.Equal(t, expected, stats.MaxDamage(StatsCategoryOverall, StatsPatternPvPAll))
}

func TestStats_Battles(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 100, stats.Battles(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_KdRate(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles:         100,
					SurvivedBattles: 60,
					Frags:           20,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 0.5, stats.KdRate(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_AvgKill(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
					Frags:   30,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 0.3, stats.AvgKill(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_AvgExp(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
					Xp:      150000,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 1500, stats.AvgExp(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_WinRate(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles: 100,
					Wins:    60,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 60, stats.WinRate(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_WinSurvivedRate(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Wins:         100,
					SurvivedWins: 20,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 20, stats.WinSurvivedRate(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_LoseSurvivedRate(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		WGAccountInfoData{
			Statistics: struct {
				Pvp     WGPlayerStatsValues `json:"pvp"`
				PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
			}{
				Pvp: WGPlayerStatsValues{
					Battles:         100,
					SurvivedBattles: 40,
					Wins:            60,
					SurvivedWins:    20,
				},
			},
		},
		emptyShipsStats,
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 50, stats.LoseSurvivedRate(StatsCategoryOverall, StatsPatternPvPAll), allowableDelta)
}

func TestStats_MainBatteryHitRate(t *testing.T) {
	t.Parallel()

	useShipID := 0
	stats := NewStats(
		useShipID,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp: WGShipStatsValues{
					MainBattery: struct {
						Hits  uint `json:"hits"`
						Shots uint `json:"shots"`
					}{
						Hits:  100,
						Shots: 200,
					},
				},
				ShipID: useShipID,
			},
		},
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 50, stats.MainBatteryHitRate(StatsPatternPvPAll), allowableDelta)
}

func TestStats_TorpedoesHitRate(t *testing.T) {
	t.Parallel()

	useShipID := 0
	stats := NewStats(
		useShipID,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp: WGShipStatsValues{
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
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 25, stats.TorpedoesHitRate(StatsPatternPvPAll), allowableDelta)
}

func TestStats_PlanesKilled(t *testing.T) {
	t.Parallel()

	useShipID := 0
	stats := NewStats(
		useShipID,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp: WGShipStatsValues{
					PlanesKilled: 334,
					Battles:      10,
				},
				ShipID: useShipID,
			},
		},
		emptyExpectedStats,
		emptyWarships,
	)

	assert.InDelta(t, 33.4, stats.PlanesKilled(StatsPatternPvPAll), allowableDelta)
}

func TestStats_AvgTier(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp:    WGShipStatsValues{Battles: 20},
				ShipID: 100,
			},
			{
				Pvp:    WGShipStatsValues{Battles: 50},
				ShipID: 200,
			},
		},
		emptyExpectedStats,
		Warships{
			100: {Tier: 5},
			200: {Tier: 8},
		},
	)

	assert.InDelta(t, 7.14, stats.AvgTier(StatsPatternPvPAll), allowableDelta)
}

func TestStats_UsingTierRate(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp:    WGShipStatsValues{Battles: 30},
				ShipID: 100,
			},
			{
				Pvp:    WGShipStatsValues{Battles: 50},
				ShipID: 200,
			},
			{
				Pvp:    WGShipStatsValues{Battles: 20},
				ShipID: 300,
			},
		},
		emptyExpectedStats,
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

func TestStats_UsingShipTypeRate(t *testing.T) {
	t.Parallel()

	stats := NewStats(
		0,
		emptyAccountInfo,
		[]WGShipsStatsData{
			{
				Pvp:    WGShipStatsValues{Battles: 30},
				ShipID: 100,
			},
			{
				Pvp:    WGShipStatsValues{Battles: 50},
				ShipID: 200,
			},
			{
				Pvp:    WGShipStatsValues{Battles: 20},
				ShipID: 300,
			},
			{
				Pvp:    WGShipStatsValues{Battles: 100},
				ShipID: 400,
			},
		},
		emptyExpectedStats,
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
