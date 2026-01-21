package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testTolerance = 0.1

func TestCalculatePR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		avgDamage       float64
		avgFrags        float64
		winRate         float64
		expectedDamage  float64
		expectedFrags   float64
		expectedWinRate float64
		expectedPR      float64
	}{
		{
			name:            "全て期待値と同じ場合",
			avgDamage:       50000,
			avgFrags:        1.0,
			winRate:         50,
			expectedDamage:  50000,
			expectedFrags:   1.0,
			expectedWinRate: 50,
			expectedPR:      1150,
		},
		{
			name:            "ダメージが期待値の1.5倍の場合",
			avgDamage:       75000,
			avgFrags:        1.0,
			winRate:         50,
			expectedDamage:  50000,
			expectedFrags:   1.0,
			expectedWinRate: 50,
			expectedPR:      1733.33,
		},
		{
			name:            "フラグが期待値の2倍の場合",
			avgDamage:       50000,
			avgFrags:        2.0,
			winRate:         50,
			expectedDamage:  50000,
			expectedFrags:   1.0,
			expectedWinRate: 50,
			expectedPR:      1483.33,
		},
		{
			name:            "勝率が期待値より高い場合",
			avgDamage:       50000,
			avgFrags:        1.0,
			winRate:         60,
			expectedDamage:  50000,
			expectedFrags:   1.0,
			expectedWinRate: 50,
			expectedPR:      1250,
		},
		{
			name:            "全ての指標が低い場合（0以下の値は0へ）",
			avgDamage:       20000,
			avgFrags:        0.1,
			winRate:         30,
			expectedDamage:  50000,
			expectedFrags:   1.0,
			expectedWinRate: 50,
			expectedPR:      0,
		},
		{
			name:            "期待値がゼロの場合",
			avgDamage:       50000,
			avgFrags:        1.0,
			winRate:         50,
			expectedDamage:  0,
			expectedFrags:   0,
			expectedWinRate: 0,
			expectedPR:      0,
		},
		{
			name:            "混合パターン：ダメージと勝率が高い",
			avgDamage:       70000,
			avgFrags:        1.5,
			winRate:         58,
			expectedDamage:  50000,
			expectedFrags:   1.0,
			expectedWinRate: 50,
			expectedPR:      1863.33,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePR(
				tt.avgDamage,
				tt.avgFrags,
				tt.winRate,
				tt.expectedDamage,
				tt.expectedFrags,
				tt.expectedWinRate,
			)

			assert.InDelta(t, tt.expectedPR, result, testTolerance)
		})
	}
}

func TestCalculateOverallPR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		pattern         StatsPattern
		playerShipStats PlayerShipStats
		warships        Warships
		expectedPR      float64
	}{
		{
			name:            "空のプレイヤー船舶統計の場合",
			pattern:         StatsPatternPvPAll,
			playerShipStats: PlayerShipStats{},
			warships:        Warships{},
			expectedPR:      0,
		},
		{
			name:    "戦闘数ゼロの船舶を無視する場合",
			pattern: StatsPatternPvPAll,
			playerShipStats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles: 0,
					},
				},
			},
			warships:   Warships{},
			expectedPR: 0,
		},
		{
			name:    "サーバー平均が存在しない船舶を無視する場合",
			pattern: StatsPatternPvPAll,
			playerShipStats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles:     10,
						DamageDealt: 500000,
						Frags:       10,
						Wins:        5,
					},
				},
			},
			warships: Warships{
				1: *NewUnknownWarship(),
			},
			expectedPR: 0,
		},
		{
			name:    "対象のWarshipが存在しない場合",
			pattern: StatsPatternPvPAll,
			playerShipStats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles:     10,
						DamageDealt: 500000,
						Frags:       10,
						Wins:        5,
					},
				},
			},
			warships:   Warships{},
			expectedPR: 0,
		},
		{
			name:    "単一の船舶でPRを計算する場合",
			pattern: StatsPatternPvPAll,
			playerShipStats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles:     10,
						DamageDealt: 500000,
						Frags:       10,
						Wins:        5,
					},
				},
			},
			warships: Warships{
				1: *NewWarship(
					1,
					"TestShip",
					8,
					ShipTypeBB,
					Nation("usa"),
					false,
					&ServerAverage{
						Damage:  50000,
						Frags:   1.0,
						WinRate: 50,
					},
				),
			},
			expectedPR: 1150,
		},
		{
			name:    "複数の船舶で合計PRを計算する場合",
			pattern: StatsPatternPvPAll,
			playerShipStats: PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles:     10,
						DamageDealt: 500000,
						Frags:       10,
						Wins:        5,
					},
				},
				2: WGShipsStatsData{
					ShipID: 2,
					Pvp: WGShipStatsValues{
						Battles:     5,
						DamageDealt: 250000,
						Frags:       5,
						Wins:        2,
					},
				},
			},
			warships: Warships{
				1: *NewWarship(
					1,
					"TestShip1",
					8,
					ShipTypeBB,
					Nation("usa"),
					false,
					&ServerAverage{
						Damage:  50000,
						Frags:   1.0,
						WinRate: 50,
					},
				),
				2: *NewWarship(
					2,
					"TestShip2",
					8,
					ShipTypeBB,
					Nation("usa"),
					false,
					&ServerAverage{
						Damage:  50000,
						Frags:   1.0,
						WinRate: 50,
					},
				),
			},
			expectedPR: 1116.67,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateOverallPR(tt.pattern, tt.playerShipStats, tt.warships)

			assert.InDelta(t, tt.expectedPR, result, testTolerance)
		})
	}
}
