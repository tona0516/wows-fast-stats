package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuildPlayerStatsCalculationConsistency は BuildPlayerStats 関数が正常に PlayerStats を構築することを検証する。
func TestBuildPlayerStatsCalculationConsistency(t *testing.T) {
	t.Parallel()

	// テスト用のモックデータを作成
	tempArenaInfo := TempArenaInfo{
		Vehicles: []Vehicle{
			{
				Name:   "TestPlayer1",
				ShipID: 3591915536, // T10 DD
			},
		},
		DateTime:   "2025-01-01T00:00:00Z",
		MapID:      1,
		MatchGroup: "unknown",
		PlayerName: "TestPlayer1",
	}

	accountInfo := WGAccountInfoData{
		HiddenProfile: false,
	}

	allShipsStats := []WGShipsStatsData{
		{ShipID: 3591915536},
	}

	shipsBadges := []WGShipsBadgesData{}

	warships := Warships{
		3591915536: *NewUnknownWarship(),
	}

	stats := NewPersonalStats(
		3591915536,
		accountInfo,
		allShipsStats,
		shipsBadges,
		warships,
		tempArenaInfo,
	)

	// BuildPlayerStats が正常に実行され、結果が返されることを検証
	result := BuildPlayerStats(
		StatsPatternPvPSolo,
		stats,
		123456,
		3591915536,
		tempArenaInfo,
		warships,
	)

	// 基本的な検証
	assert.NotNil(t, result)
	assert.NotNil(t, result.ShipStats)
	assert.NotNil(t, result.OverallStats)

	// ShipStats のフィールドが正しく設定されることを確認
	assert.GreaterOrEqual(t, result.ShipStats.Battles, uint(0))
	assert.NotNil(t, result.ShipStats.Damage)
	assert.NotNil(t, result.ShipStats.WinRate)

	// OverallStats のフィールドが正しく設定されることを確認
	assert.GreaterOrEqual(t, result.OverallStats.Battles, uint(0))
	assert.NotNil(t, result.OverallStats.Damage)
	assert.NotNil(t, result.OverallStats.WinRate)
	assert.NotNil(t, result.OverallStats.ThreatLevel)
}

// TestCalculateTeamAverageStats は単一のチーム平均統計計算を検証する。
func TestCalculateTeamAverageStats(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PvPSolo: PlayerStats{
				ShipStats: ShipStats{
					Battles: 10,
					PR:      RatingValue{Value: 2000},
					Damage:  RatingValue{Value: 50000},
					WinRate: RatingValue{Value: 55},
				},
				OverallStats: OverallStats{
					Battles: 100,
					PR:      RatingValue{Value: 2500},
					Damage:  RatingValue{Value: 45000},
					WinRate: RatingValue{Value: 52},
				},
			},
		},
		{
			PvPSolo: PlayerStats{
				ShipStats: ShipStats{
					Battles: 5,
					PR:      RatingValue{Value: 1800},
					Damage:  RatingValue{Value: 48000},
					WinRate: RatingValue{Value: 50},
				},
				OverallStats: OverallStats{
					Battles: 80,
					PR:      RatingValue{Value: 2300},
					Damage:  RatingValue{Value: 44000},
					WinRate: RatingValue{Value: 51},
				},
			},
		},
	}

	result := CalculateTeamAverageStats(players, StatsPatternPvPSolo)

	// 結果の整合性を検証
	assert.NotNil(t, result)
	assert.Greater(t, result.ShipPR, float64(0))
	assert.Greater(t, result.ShipDamage, float64(0))
	assert.Greater(t, result.OverallPR, float64(0))
	assert.Greater(t, result.OverallDamage, float64(0))

	// 期待値の計算
	expectedShipPR := (2000 + 1800) / 2.0
	expectedShipDamage := (50000 + 48000) / 2.0

	assert.InDelta(t, expectedShipPR, result.ShipPR, 0.1, "平均PRが正しく計算されるべき")
	assert.InDelta(t, expectedShipDamage, result.ShipDamage, 0.1, "平均ダメージが正しく計算されるべき")
}

// TestCalculateTeamAverageStatsEmptyPlayers はプレイヤーがいない場合のハンドリングを検証する。
func TestCalculateTeamAverageStatsEmptyPlayers(t *testing.T) {
	t.Parallel()

	players := Players{}

	result := CalculateTeamAverageStats(players, StatsPatternPvPSolo)

	// 結果が安全に処理されることを確認
	assert.NotNil(t, result)
	assert.Equal(t, float64(0), result.ShipPR)
	assert.Equal(t, float64(0), result.ShipDamage)
	assert.Equal(t, uint(0), result.ShipBattles)
}

// TestCalculateTeamAverageStatsZeroBattles は戦闘数がゼロのプレイヤーを無視することを検証する。
func TestCalculateTeamAverageStatsZeroBattles(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PvPSolo: PlayerStats{
				ShipStats: ShipStats{
					Battles: 0, // ゼロバトル (無視されるべき)
					PR:      RatingValue{Value: 2000},
					Damage:  RatingValue{Value: 50000},
					WinRate: RatingValue{Value: 55},
				},
				OverallStats: OverallStats{
					Battles: 0,
					PR:      RatingValue{Value: 2500},
					Damage:  RatingValue{Value: 45000},
					WinRate: RatingValue{Value: 52},
				},
			},
		},
		{
			PvPSolo: PlayerStats{
				ShipStats: ShipStats{
					Battles: 10,
					PR:      RatingValue{Value: 1800},
					Damage:  RatingValue{Value: 48000},
					WinRate: RatingValue{Value: 50},
				},
				OverallStats: OverallStats{
					Battles: 100,
					PR:      RatingValue{Value: 2300},
					Damage:  RatingValue{Value: 44000},
					WinRate: RatingValue{Value: 51},
				},
			},
		},
	}

	result := CalculateTeamAverageStats(players, StatsPatternPvPSolo)

	// ゼロバトルプレイヤーは無視され、2番目のプレイヤーのみで計算されるべき
	assert.NotNil(t, result)
	assert.InDelta(t, 1800, result.ShipPR, 0.1, "1プレイヤーのみで計算されるべき")
	assert.InDelta(t, 48000, result.ShipDamage, 0.1, "1プレイヤーのみで計算されるべき")
}
