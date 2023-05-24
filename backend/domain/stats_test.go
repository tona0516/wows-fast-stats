package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"changeme/backend/vo"
)

func TestStats_SetShipStats(t *testing.T) {
	t.Parallel()

	stats := Stats{}
	shipStats := vo.WGShipsStatsData{Pvp: vo.WGStatsValues{Battles: 100}}
	stats.SetShipStats(shipStats)

	assert.Equal(t, shipStats, stats.ShipsStats)
}

func TestStats_ShipPR(t *testing.T) {
	t.Parallel()

	expected := vo.NSExpectedStatsData{
		AverageDamageDealt: 8000,
		AverageFrags:       1,
		WinRate:            50,
	}

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			Pvp: vo.WGStatsValues{
				Battles:     100,
				DamageDealt: 1000000,
				Frags:       200,
				Wins:        60,
			},
		},
		Expected: expected,
	}

	assert.Equal(t, float64(1875), stats.ShipPR())
}
func TestStats_AvgDamage_Overall(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				Pvp: vo.WGStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
	}

	assert.Equal(t, float64(10000), stats.AvgDamage(ModeOverall))
}

func TestStats_AvgDamage_OverallSolo(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				PvpSolo: vo.WGStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
	}

	assert.Equal(t, float64(10000), stats.AvgDamage(ModeOverallSolo))
}

func TestStats_AvgDamage_Ship(t *testing.T) {
	t.Parallel()

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			Pvp: vo.WGStatsValues{
				Battles:     100,
				DamageDealt: 1000000,
			},
		},
	}

	assert.Equal(t, float64(10000), stats.AvgDamage(ModeShip))
}

func TestStats_AvgDamage_ShipSolo(t *testing.T) {
	t.Parallel()

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			PvpSolo: vo.WGStatsValues{
				Battles:     100,
				DamageDealt: 1000000,
			},
		},
	}

	assert.Equal(t, float64(10000), stats.AvgDamage(ModeShipSolo))
}

func TestStats_KdRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				Pvp: vo.WGStatsValues{
					Battles:         100,
					SurvivedBattles: 60,
					Frags:           20,
				},
			},
		},
	}

	assert.Equal(t, float64(0.5), stats.KdRate(ModeOverall))
}

func TestStats_AvgExp(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				Pvp: vo.WGStatsValues{
					Battles: 100,
					Xp:      150000,
				},
			},
		},
	}

	assert.Equal(t, float64(1500), stats.AvgExp(ModeOverall))
}

func TestStats_WinRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				Pvp: vo.WGStatsValues{
					Battles: 100,
					Wins:    60,
				},
			},
		},
	}

	assert.Equal(t, float64(60), stats.WinRate(ModeOverall))
}

func TestStats_WinSurvivedRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				Pvp: vo.WGStatsValues{
					Wins:         100,
					SurvivedWins: 20,
				},
			},
		},
	}

	assert.Equal(t, float64(20), stats.WinSurvivedRate(ModeOverall))
}

func TestStats_LoseSurvivedRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues "json:\"pvp\""
				PvpSolo vo.WGStatsValues "json:\"pvp_solo\""
			}{
				Pvp: vo.WGStatsValues{
					Battles:         100,
					SurvivedBattles: 40,
					Wins:            60,
					SurvivedWins:    20,
				},
			},
		},
	}

	assert.Equal(t, float64(50), stats.LoseSurvivedRate(ModeOverall))
}

func TestStats_MainBatteryHitRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			Pvp: vo.WGStatsValues{
				MainBattery: struct {
					Hits  uint "json:\"hits\""
					Shots uint "json:\"shots\""
				}{
					Hits:  100,
					Shots: 200,
				},
			},
		},
	}

	assert.Equal(t, float64(50), stats.MainBatteryHitRate(ModeShip))
}

func TestStats_TorpedoesHitRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			Pvp: vo.WGStatsValues{
				Torpedoes: struct {
					Hits  uint "json:\"hits\""
					Shots uint "json:\"shots\""
				}{
					Hits:  10,
					Shots: 40,
				},
			},
		},
	}

	assert.Equal(t, float64(25), stats.TorpedoesHitRate(ModeShip))
}

func TestStats_AvgTier(t *testing.T) {
	t.Parallel()

	accountID := 10
	shipInfo := map[int]vo.Warship{
		100: {Tier: 5},
		200: {Tier: 8},
	}
	shipStats := map[int]vo.WGShipsStats{
		10: {
			Data: map[int][]vo.WGShipsStatsData{
				10: {
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 20},
						ShipID: 100,
					},
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 50},
						ShipID: 200,
					},
				},
			},
		},
	}

	stats := Stats{}
	assert.InDelta(t, float64(7.14), stats.AvgTier(accountID, shipInfo, shipStats), 0.01)
}

func TestStats_UsingTierRate(t *testing.T) {
	t.Parallel()

	accountID := 10
	shipInfo := map[int]vo.Warship{
		100: {Tier: 5},
		200: {Tier: 8},
		300: {Tier: 4},
	}
	shipStats := map[int]vo.WGShipsStats{
		10: {
			Data: map[int][]vo.WGShipsStatsData{
				10: {
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 30},
						ShipID: 100,
					},
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 50},
						ShipID: 200,
					},
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 20},
						ShipID: 300,
					},
				},
			},
		},
	}

	stats := Stats{}
	tierGroup := stats.UsingTierRate(accountID, shipInfo, shipStats)
	assert.Equal(t, float64(20), tierGroup.Low)
	assert.Equal(t, float64(30), tierGroup.Middle)
	assert.Equal(t, float64(50), tierGroup.High)
}

func TestStats_UsingShipTypeRate(t *testing.T) {
	t.Parallel()

	accountID := 10
	shipInfo := map[int]vo.Warship{
		100: {Type: vo.DD},
		200: {Type: vo.CL},
		300: {Type: vo.BB},
		400: {Type: vo.CV},
	}
	shipStats := map[int]vo.WGShipsStats{
		10: {
			Data: map[int][]vo.WGShipsStatsData{
				10: {
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 30},
						ShipID: 100,
					},
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 50},
						ShipID: 200,
					},
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 20},
						ShipID: 300,
					},
					vo.WGShipsStatsData{
						Pvp:    vo.WGStatsValues{Battles: 100},
						ShipID: 400,
					},
				},
			},
		},
	}

	stats := Stats{}
	shipTypeGroup := stats.UsingShipTypeRate(accountID, shipInfo, shipStats)
	assert.Equal(t, float64(0), shipTypeGroup.SS)
	assert.Equal(t, float64(15), shipTypeGroup.DD)
	assert.Equal(t, float64(25), shipTypeGroup.CL)
	assert.Equal(t, float64(10), shipTypeGroup.BB)
	assert.Equal(t, float64(50), shipTypeGroup.CV)
}
