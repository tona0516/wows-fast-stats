package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGShipsStatsData_shipStatsValues(t *testing.T) {
	t.Parallel()

	data := WGShipsStatsData{
		Pvp: WGShipStatsValues{
			Battles:     100,
			Wins:        50,
			DamageDealt: 5000,
		},
		PvpSolo: WGShipStatsValues{
			Battles:     50,
			Wins:        25,
			DamageDealt: 2500,
		},
		RankSolo: WGShipStatsValues{
			Battles:     30,
			Wins:        15,
			DamageDealt: 1500,
		},
	}

	tests := []struct {
		name            string
		pattern         StatsPattern
		expectedBattles uint
		expectedWins    uint
		expectedDamage  uint
	}{
		{
			name:            "PvPAll",
			pattern:         StatsPatternPvPAll,
			expectedBattles: 100,
			expectedWins:    50,
			expectedDamage:  5000,
		},
		{
			name:            "PvPSolo",
			pattern:         StatsPatternPvPSolo,
			expectedBattles: 50,
			expectedWins:    25,
			expectedDamage:  2500,
		},
		{
			name:            "RankSolo",
			pattern:         StatsPatternRankSolo,
			expectedBattles: 30,
			expectedWins:    15,
			expectedDamage:  1500,
		},
		{
			name:            "不明なパターン",
			pattern:         StatsPattern("unknown"),
			expectedBattles: 0,
			expectedWins:    0,
			expectedDamage:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.shipStatsValues(tt.pattern)

			assert.Equal(t, tt.expectedBattles, result.Battles)
			assert.Equal(t, tt.expectedWins, result.Wins)
			assert.Equal(t, tt.expectedDamage, result.DamageDealt)
		})
	}
}

func TestWGShipsStatsData_platoonRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		soloBattles uint
		div2Battles uint
		div3Battles uint
		allBattles  uint
		expected    float64
	}{
		{
			name:        "正常系",
			allBattles:  100,
			soloBattles: 40,
			div2Battles: 30,
			div3Battles: 30,
			expected:    1.9,
		},
		{
			name:        "ゼロ戦闘",
			allBattles:  0,
			soloBattles: 0,
			div2Battles: 0,
			div3Battles: 0,
			expected:    0.0,
		},
		{
			name:        "ソロのみ",
			allBattles:  100,
			soloBattles: 100,
			div2Battles: 0,
			div3Battles: 0,
			expected:    1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := WGShipsStatsData{
				Pvp: WGShipStatsValues{
					Battles: tt.allBattles,
				},
				PvpSolo: WGShipStatsValues{
					Battles: tt.soloBattles,
				},
				PvpDiv2: struct {
					Battles uint `json:"battles"`
				}{
					Battles: tt.div2Battles,
				},
				PvpDiv3: struct {
					Battles uint `json:"battles"`
				}{
					Battles: tt.div3Battles,
				},
			}

			result := data.platoonRate()

			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

func TestWGShipStatsValues_avgDamage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		damageDealt uint
		battles     uint
		expected    float64
	}{
		{
			name:        "正常系",
			damageDealt: 5000,
			battles:     100,
			expected:    50.0,
		},
		{
			name:        "ゼロ戦闘",
			damageDealt: 5000,
			battles:     0,
			expected:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				DamageDealt: tt.damageDealt,
				Battles:     tt.battles,
			}

			result := values.avgDamage()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWGShipStatsValues_avgKills(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		frags    uint
		battles  uint
		expected float64
	}{
		{
			name:     "正常系",
			frags:    50,
			battles:  100,
			expected: 0.5,
		},
		{
			name:     "ゼロ戦闘",
			frags:    50,
			battles:  0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				Frags:   tt.frags,
				Battles: tt.battles,
			}

			result := values.avgKills()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWGShipStatsValues_winRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		wins     uint
		battles  uint
		expected float64
	}{
		{
			name:     "正常系",
			wins:     60,
			battles:  100,
			expected: 60.0,
		},
		{
			name:     "ゼロ戦闘",
			wins:     60,
			battles:  0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				Wins:    tt.wins,
				Battles: tt.battles,
			}

			result := values.winRate()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWGShipStatsValues_kdRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		frags           uint
		battles         uint
		survivedBattles uint
		expected        float64
	}{
		{
			name:            "正常系",
			frags:           100,
			battles:         50,
			survivedBattles: 30,
			expected:        5.0,
		},
		{
			name:            "すべて生存",
			frags:           100,
			battles:         50,
			survivedBattles: 50,
			expected:        100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				Frags:           tt.frags,
				Battles:         tt.battles,
				SurvivedBattles: tt.survivedBattles,
			}

			result := values.kdRate()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWGShipStatsValues_exp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		xp       uint
		battles  uint
		expected float64
	}{
		{
			name:     "正常系",
			xp:       500000,
			battles:  100,
			expected: 5000.0,
		},
		{
			name:     "ゼロ戦闘",
			xp:       500000,
			battles:  0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				Xp:      tt.xp,
				Battles: tt.battles,
			}

			result := values.exp()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWGShipStatsValues_survivedRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		battles         uint
		wins            uint
		survivedBattles uint
		survivedWins    uint
		expectedAll     float64
		expectedWin     float64
		expectedLose    float64
	}{
		{
			name:            "正常系",
			battles:         100,
			wins:            60,
			survivedBattles: 50,
			survivedWins:    40,
			expectedAll:     50.0,
			expectedWin:     66.667,
			expectedLose:    25.0,
		},
		{
			name:            "ゼロ戦闘",
			battles:         0,
			wins:            0,
			survivedBattles: 0,
			survivedWins:    0,
			expectedAll:     0.0,
			expectedWin:     0.0,
			expectedLose:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				Battles:         tt.battles,
				Wins:            tt.wins,
				SurvivedBattles: tt.survivedBattles,
				SurvivedWins:    tt.survivedWins,
			}

			result := values.survivedRate()

			assert.Equal(t, tt.expectedAll, result.All)
			assert.InDelta(t, tt.expectedWin, result.Win, 0.001)
			assert.Equal(t, tt.expectedLose, result.Lose)
		})
	}
}

func TestWGShipStatsValues_avgPlanesKilled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		planesKilled uint
		battles      uint
		expected     float64
	}{
		{
			name:         "正常系",
			planesKilled: 200,
			battles:      100,
			expected:     2.0,
		},
		{
			name:         "ゼロ戦闘",
			planesKilled: 200,
			battles:      0,
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				PlanesKilled: tt.planesKilled,
				Battles:      tt.battles,
			}

			result := values.avgPlanesKilled()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWGShipStatsValues_hitRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		mainBatteryHits     uint
		mainBatteryShots    uint
		torpedoesHits       uint
		torpedoesShots      uint
		expectedMainBattery float64
		expectedTorpedoes   float64
	}{
		{
			name:                "正常系",
			mainBatteryHits:     500,
			mainBatteryShots:    1000,
			torpedoesHits:       200,
			torpedoesShots:      400,
			expectedMainBattery: 50.0,
			expectedTorpedoes:   50.0,
		},
		{
			name:                "ゼロ発射",
			mainBatteryHits:     500,
			mainBatteryShots:    0,
			torpedoesHits:       200,
			torpedoesShots:      0,
			expectedMainBattery: 0.0,
			expectedTorpedoes:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := WGShipStatsValues{
				MainBattery: WGArmament{
					Hits:  tt.mainBatteryHits,
					Shots: tt.mainBatteryShots,
				},
				Torpedoes: WGArmament{
					Hits:  tt.torpedoesHits,
					Shots: tt.torpedoesShots,
				},
			}

			result := values.hitRate()

			assert.Equal(t, tt.expectedMainBattery, result.MainBattery)
			assert.Equal(t, tt.expectedTorpedoes, result.Torpedoes)
		})
	}
}

func TestWGArmament_hitRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		hits     uint
		shots    uint
		expected float64
	}{
		{
			name:     "命中率50%",
			hits:     500,
			shots:    1000,
			expected: 50.0,
		},
		{
			name:     "命中率100%",
			hits:     1000,
			shots:    1000,
			expected: 100.0,
		},
		{
			name:     "命中率0%",
			hits:     0,
			shots:    1000,
			expected: 0.0,
		},
		{
			name:     "ゼロ発射",
			hits:     100,
			shots:    0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			armament := WGArmament{
				Hits:  tt.hits,
				Shots: tt.shots,
			}

			result := armament.hitRate()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlayerShipStats_avgTier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		stats    PlayerShipStats
		warships Warships
		expected float64
	}{
		{
			name: "正常系",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
				2: WGShipsStatsData{
					ShipID: 2,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 8,
				},
				2: {
					ID:   2,
					Tier: 10,
				},
			},
			expected: 9.0,
		},
		{
			name: "不明な艦船ID",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
				999: WGShipsStatsData{
					ShipID: 999,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 8,
				},
			},
			expected: 8.0,
		},
		{
			name: "ゼロ戦闘",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 0,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 8,
				},
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.avgTier(StatsPatternPvPAll, tt.warships)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlayerShipStats_usingShipTypeRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		stats      PlayerShipStats
		warships   Warships
		expectedBB float64
		expectedDD float64
		expectedCL float64
		expectedSS float64
		expectedCV float64
	}{
		{
			name: "正常系",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
				2: WGShipsStatsData{
					ShipID: 2,
					Pvp: WGShipStatsValues{
						Battles: 30,
					},
				},
				3: WGShipsStatsData{
					ShipID: 3,
					Pvp: WGShipStatsValues{
						Battles: 20,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Type: ShipTypeBB,
				},
				2: {
					ID:   2,
					Type: ShipTypeDD,
				},
				3: {
					ID:   3,
					Type: ShipTypeBB,
				},
			},
			expectedBB: 70.0,
			expectedDD: 30.0,
			expectedCL: 0.0,
			expectedSS: 0.0,
			expectedCV: 0.0,
		},
		{
			name: "不明な艦船ID",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
				999: WGShipsStatsData{
					ShipID: 999,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Type: ShipTypeBB,
				},
			},
			expectedBB: 100.0,
			expectedDD: 0.0,
			expectedCL: 0.0,
			expectedSS: 0.0,
			expectedCV: 0.0,
		},
		{
			name: "ゼロ戦闘",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 0,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Type: ShipTypeBB,
				},
			},
			expectedBB: 0.0,
			expectedDD: 0.0,
			expectedCL: 0.0,
			expectedSS: 0.0,
			expectedCV: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.usingShipTypeRate(StatsPatternPvPAll, tt.warships)

			assert.Equal(t, tt.expectedBB, result.BB)
			assert.Equal(t, tt.expectedDD, result.DD)
			assert.Equal(t, tt.expectedCL, result.CL)
			assert.Equal(t, tt.expectedSS, result.SS)
			assert.Equal(t, tt.expectedCV, result.CV)
		})
	}
}

func TestPlayerShipStats_usingTierRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		stats          PlayerShipStats
		warships       Warships
		expectedLow    float64
		expectedMiddle float64
		expectedHigh   float64
	}{
		{
			name: "正常系",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 25,
					},
				},
				2: WGShipsStatsData{
					ShipID: 2,
					Pvp: WGShipStatsValues{
						Battles: 25,
					},
				},
				3: WGShipsStatsData{
					ShipID: 3,
					Pvp: WGShipStatsValues{
						Battles: 25,
					},
				},
				4: WGShipsStatsData{
					ShipID: 4,
					Pvp: WGShipStatsValues{
						Battles: 25,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 3,
				},
				2: {
					ID:   2,
					Tier: 6,
				},
				3: {
					ID:   3,
					Tier: 9,
				},
				4: {
					ID:   4,
					Tier: 8,
				},
			},
			expectedLow:    25.0,
			expectedMiddle: 25.0,
			expectedHigh:   50.0,
		},
		{
			name: "不明な艦船ID",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
				999: WGShipsStatsData{
					ShipID: 999,
					Pvp: WGShipStatsValues{
						Battles: 50,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 8,
				},
			},
			expectedLow:    0.0,
			expectedMiddle: 0.0,
			expectedHigh:   100.0,
		},
		{
			name: "ゼロ戦闘",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 0,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 8,
				},
			},
			expectedLow:    0.0,
			expectedMiddle: 0.0,
			expectedHigh:   0.0,
		},
		{
			name: "各ティア境界",
			stats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 1,
					},
				},
				2: WGShipsStatsData{
					ShipID: 2,
					Pvp: WGShipStatsValues{
						Battles: 1,
					},
				},
				3: WGShipsStatsData{
					ShipID: 3,
					Pvp: WGShipStatsValues{
						Battles: 1,
					},
				},
				4: WGShipsStatsData{
					ShipID: 4,
					Pvp: WGShipStatsValues{
						Battles: 1,
					},
				},
				5: WGShipsStatsData{
					ShipID: 5,
					Pvp: WGShipStatsValues{
						Battles: 1,
					},
				},
			},
			warships: Warships{
				1: {
					ID:   1,
					Tier: 4,
				},
				2: {
					ID:   2,
					Tier: 5,
				},
				3: {
					ID:   3,
					Tier: 7,
				},
				4: {
					ID:   4,
					Tier: 8,
				},
				5: {
					ID:   5,
					Tier: 10,
				},
			},
			expectedLow:    20.0,
			expectedMiddle: 40.0,
			expectedHigh:   40.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.usingTierRate(StatsPatternPvPAll, tt.warships)

			assert.Equal(t, tt.expectedLow, result.Low)
			assert.Equal(t, tt.expectedMiddle, result.Middle)
			assert.Equal(t, tt.expectedHigh, result.High)
		})
	}
}
