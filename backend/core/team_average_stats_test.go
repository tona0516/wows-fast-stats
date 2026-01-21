package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRatingValue(value float64) RatingValue {
	return RatingValue{Value: value}
}

func newShipStats(battles uint, pr, damage, winRate float64) ShipStats {
	return ShipStats{
		Battles: battles,
		PR:      newRatingValue(pr),
		Damage:  newRatingValue(damage),
		WinRate: newRatingValue(winRate),
	}
}

func newOverallStats(battles uint, pr, damage, winRate float64) OverallStats {
	return OverallStats{
		Battles: battles,
		PR:      newRatingValue(pr),
		Damage:  newRatingValue(damage),
		WinRate: newRatingValue(winRate),
	}
}

func newPlayerStats(shipStats ShipStats, overallStats OverallStats) PlayerStats {
	return PlayerStats{
		ShipStats:    shipStats,
		OverallStats: overallStats,
	}
}

func newPlayerWithPvPSolo(shipStats ShipStats, overallStats OverallStats) Player {
	return Player{
		PvPSolo: newPlayerStats(shipStats, overallStats),
	}
}

func newPlayerWithPvPAll(shipStats ShipStats, overallStats OverallStats) Player {
	return Player{
		PvPAll: newPlayerStats(shipStats, overallStats),
	}
}

func newPlayerWithRankSolo(shipStats ShipStats, overallStats OverallStats) Player {
	return Player{
		RankSolo: newPlayerStats(shipStats, overallStats),
	}
}

func TestNewTeamAverageStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		players      Players
		statsPattern StatsPattern
		expected     TeamAverageStats
	}{
		{
			name:         "プレイヤーなし",
			players:      Players{},
			statsPattern: StatsPatternPvPSolo,
			expected: TeamAverageStats{
				ShipPR:         0,
				ShipDamage:     0,
				ShipWinRate:    0,
				ShipBattles:    0,
				OverallPR:      0,
				OverallDamage:  0,
				OverallWinRate: 0,
				OverallBattles: 0,
			},
		},
		{
			name: "PvPSolo単一プレイヤー",
			players: Players{
				newPlayerWithPvPSolo(
					newShipStats(100, 2000, 50000, 55.5),
					newOverallStats(500, 1800, 48000, 52.0),
				),
			},
			statsPattern: StatsPatternPvPSolo,
			expected: TeamAverageStats{
				ShipPR:         2000,
				ShipDamage:     50000,
				ShipWinRate:    55.5,
				ShipBattles:    100,
				OverallPR:      1800,
				OverallDamage:  48000,
				OverallWinRate: 52.0,
				OverallBattles: 500,
			},
		},
		{
			name: "PvPAll複数プレイヤー",
			players: Players{
				newPlayerWithPvPAll(
					newShipStats(100, 2000, 50000, 60),
					newOverallStats(500, 1800, 48000, 55),
				),
				newPlayerWithPvPAll(
					newShipStats(200, 3000, 60000, 50),
					newOverallStats(600, 2200, 55000, 52),
				),
			},
			statsPattern: StatsPatternPvPAll,
			expected: TeamAverageStats{
				ShipPR:         2500,  // (2000 + 3000) / 2
				ShipDamage:     55000, // (50000 + 60000) / 2
				ShipWinRate:    55,    // (60 + 50) / 2
				ShipBattles:    150,   // (100 + 200) / 2
				OverallPR:      2000,  // (1800 + 2200) / 2
				OverallDamage:  51500, // (48000 + 55000) / 2
				OverallWinRate: 53.5,  // (55 + 52) / 2
				OverallBattles: 550,   // (500 + 600) / 2
			},
		},
		{
			name: "RankSolo統計情報なしのプレイヤーを含む",
			players: Players{
				newPlayerWithRankSolo(
					newShipStats(150, 2500, 55000, 58),
					newOverallStats(0, 0, 0, 0),
				),
				newPlayerWithRankSolo(
					newShipStats(0, 0, 0, 0),
					newOverallStats(700, 2400, 52000, 54),
				),
			},
			statsPattern: StatsPatternRankSolo,
			expected: TeamAverageStats{
				ShipPR:         2500,  // 1プレイヤーのみ
				ShipDamage:     55000, // 1プレイヤーのみ
				ShipWinRate:    58,    // 1プレイヤーのみ
				ShipBattles:    150,   // 1プレイヤーのみ
				OverallPR:      2400,  // 1プレイヤーのみ
				OverallDamage:  52000, // 1プレイヤーのみ
				OverallWinRate: 54,    // 1プレイヤーのみ
				OverallBattles: 700,   // 1プレイヤーのみ
			},
		},
		{
			name: "全プレイヤーが統計情報なし",
			players: Players{
				newPlayerWithPvPSolo(
					newShipStats(0, 0, 0, 0),
					newOverallStats(0, 0, 0, 0),
				),
				newPlayerWithPvPSolo(
					newShipStats(0, 0, 0, 0),
					newOverallStats(0, 0, 0, 0),
				),
			},
			statsPattern: StatsPatternPvPSolo,
			expected: TeamAverageStats{
				ShipPR:         0,
				ShipDamage:     0,
				ShipWinRate:    0,
				ShipBattles:    0,
				OverallPR:      0,
				OverallDamage:  0,
				OverallWinRate: 0,
				OverallBattles: 0,
			},
		},
		{
			name: "小数値の平均計算",
			players: Players{
				newPlayerWithPvPSolo(
					newShipStats(50, 1500.5, 45000.75, 51.25),
					newOverallStats(300, 1600.25, 46000.5, 49.75),
				),
				newPlayerWithPvPSolo(
					newShipStats(50, 1700.5, 52000.25, 56.75),
					newOverallStats(300, 1800.75, 54000.5, 58.25),
				),
			},
			statsPattern: StatsPatternPvPSolo,
			expected: TeamAverageStats{
				ShipPR:         1600.5,  // (1500.5 + 1700.5) / 2
				ShipDamage:     48500.5, // (45000.75 + 52000.25) / 2
				ShipWinRate:    54,      // (51.25 + 56.75) / 2
				ShipBattles:    50,      // (50 + 50) / 2
				OverallPR:      1700.5,  // (1600.25 + 1800.75) / 2
				OverallDamage:  50000.5, // (46000.5 + 54000.5) / 2
				OverallWinRate: 54,      // (49.75 + 58.25) / 2
				OverallBattles: 300,     // (300 + 300) / 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewTeamAverageStats(tt.players, tt.statsPattern)
			assert.Equal(t, tt.expected, result)
		})
	}
}
