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

	assert.InDelta(t, 1875, stats.ShipPR(ModeShip), 0.1)
}
func TestStats_AvgDamage_Overall(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
	}

	assert.InDelta(t, 10000, stats.AvgDamage(ModeOverall), 0.1)
}

func TestStats_AvgDamage_OverallSolo(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				PvpSolo: vo.WGStatsValues{
					Battles:     100,
					DamageDealt: 1000000,
				},
			},
		},
	}

	assert.InDelta(t, 10000, stats.AvgDamage(ModeOverallSolo), 0.1)
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

	assert.InDelta(t, 10000, stats.AvgDamage(ModeShip), 0.1)
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

	assert.InDelta(t, 10000, stats.AvgDamage(ModeShipSolo), 0.1)
}

func TestStats_Battles(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles: 100,
				},
			},
		},
	}

	assert.InDelta(t, 100, stats.Battles(ModeOverall), 0.1)
}

func TestStats_KdRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles:         100,
					SurvivedBattles: 60,
					Frags:           20,
				},
			},
		},
	}

	assert.InDelta(t, 0.5, stats.KdRate(ModeOverall), 0.1)
}

func TestStats_AvgKill(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles: 100,
					Frags:   30,
				},
			},
		},
	}

	assert.InDelta(t, 0.3, stats.AvgKill(ModeOverall), 0.1)
}

func TestStats_AvgDeath(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles:         100,
					SurvivedBattles: 30,
				},
			},
		},
	}

	assert.InDelta(t, 0.7, stats.AvgDeath(ModeOverall), 0.1)
}

func TestStats_AvgExp(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles: 100,
					Xp:      150000,
				},
			},
		},
	}

	assert.InDelta(t, 1500, stats.AvgExp(ModeOverall), 0.1)
}

func TestStats_WinRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Battles: 100,
					Wins:    60,
				},
			},
		},
	}

	assert.InDelta(t, 60, stats.WinRate(ModeOverall), 0.1)
}

func TestStats_WinSurvivedRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
			}{
				Pvp: vo.WGStatsValues{
					Wins:         100,
					SurvivedWins: 20,
				},
			},
		},
	}

	assert.InDelta(t, 20, stats.WinSurvivedRate(ModeOverall), 0.1)
}

func TestStats_LoseSurvivedRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		AccountInfo: vo.WGAccountInfoData{
			Statistics: struct {
				Pvp     vo.WGStatsValues `json:"pvp"`
				PvpSolo vo.WGStatsValues `json:"pvp_solo"`
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

	assert.InDelta(t, 50, stats.LoseSurvivedRate(ModeOverall), 0.1)
}

func TestStats_MainBatteryHitRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			Pvp: vo.WGStatsValues{
				MainBattery: struct {
					Hits  uint `json:"hits"`
					Shots uint `json:"shots"`
				}{
					Hits:  100,
					Shots: 200,
				},
			},
		},
	}

	assert.InDelta(t, 50, stats.MainBatteryHitRate(ModeShip), 0.1)
}

func TestStats_TorpedoesHitRate(t *testing.T) {
	t.Parallel()

	stats := Stats{
		ShipsStats: vo.WGShipsStatsData{
			Pvp: vo.WGStatsValues{
				Torpedoes: struct {
					Hits  uint `json:"hits"`
					Shots uint `json:"shots"`
				}{
					Hits:  10,
					Shots: 40,
				},
			},
		},
	}

	assert.InDelta(t, 25, stats.TorpedoesHitRate(ModeShip), 0.1)
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
	assert.InDelta(t, 7.14, stats.AvgTier(ModeOverall, accountID, shipInfo, shipStats), 0.1)
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
	tierGroup := stats.UsingTierRate(ModeOverall, accountID, shipInfo, shipStats)
	assert.InDelta(t, 20, tierGroup.Low, 0.1)
	assert.InDelta(t, 30, tierGroup.Middle, 0.1)
	assert.InDelta(t, 50, tierGroup.High, 0.1)
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
	shipTypeGroup := stats.UsingShipTypeRate(ModeOverall, accountID, shipInfo, shipStats)
	assert.InDelta(t, 0, shipTypeGroup.SS, 0.1)
	assert.InDelta(t, 15, shipTypeGroup.DD, 0.1)
	assert.InDelta(t, 25, shipTypeGroup.CL, 0.1)
	assert.InDelta(t, 10, shipTypeGroup.BB, 0.1)
	assert.InDelta(t, 50, shipTypeGroup.CV, 0.1)
}
