package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShipStats_shipIDがPlayerShipStatsに存在しない場合(t *testing.T) {
	t.Parallel()

	result := NewShipStats(
		StatsPatternPvPAll,
		PlayerShipStats{},
		PlayerShipBadges{},
		ShipID(1),
		Warships{},
	)

	assert.Equal(t, ShipStats{}, result)
}

func TestNewShipStats_戦闘数が0の場合(t *testing.T) {
	t.Parallel()

	playerShipStats := PlayerShipStats{
		1: WGShipsStatsData{
			ShipID: 1,
			Pvp: WGShipStatsValues{
				Battles: 0,
			},
		},
	}

	result := NewShipStats(
		StatsPatternPvPAll,
		playerShipStats,
		PlayerShipBadges{},
		ShipID(1),
		Warships{},
	)

	assert.Equal(t, ShipStats{}, result)
}

func TestNewShipStats_正常系_Warship情報あり(t *testing.T) {
	t.Parallel()

	shipID := ShipID(1)
	playerShipStats := PlayerShipStats{
		shipID: WGShipsStatsData{
			ShipID: shipID,
			Pvp: WGShipStatsValues{
				Battles:         100,
				Wins:            60,
				DamageDealt:     500000,
				MaxDamageDealt:  10000,
				Frags:           150,
				SurvivedWins:    30,
				SurvivedBattles: 50,
				Xp:              800000,
				PlanesKilled:    10,
			},
			PvpSolo: WGShipStatsValues{
				Battles: 50,
			},
			PvpDiv2: struct {
				Battles uint `json:"battles"`
			}{Battles: 25},
			PvpDiv3: struct {
				Battles uint `json:"battles"`
			}{Battles: 25},
		},
	}

	warships := Warships{
		shipID: {
			ID:   shipID,
			Name: "TestShip",
			Tier: 8,
			Type: ShipTypeBB,
			ServerAverage: &ServerAverage{
				Damage:  50000,
				Frags:   1.5,
				WinRate: 50.0,
			},
		},
	}

	playerShipBadges := PlayerShipBadges{
		shipID: WGShipsBadgesData{
			ShipID:        shipID,
			TopGradeClass: 1,
		},
	}

	result := NewShipStats(
		StatsPatternPvPAll,
		playerShipStats,
		playerShipBadges,
		shipID,
		warships,
	)

	assert.Equal(t, uint(100), result.Battles)
	assert.Equal(t, 5000.0, result.Damage.Value)         // 500000 / 100
	assert.NotEqual(t, RatingNone, result.Damage.Rating) // ServerAverage があるため Rating が設定される
	assert.Equal(t, uint(10000), result.MaxDamage.Value)
	assert.Equal(t, 60.0, result.WinRate.Value)                  // (60 / 100) * 100
	assert.Equal(t, Rating(RatingUnicum), result.WinRate.Rating) // WinRate から計算される
	assert.NotEqual(t, RatingNone, result.PR.Rating)             // ServerAverage があるため Rating が設定される
	assert.Equal(t, 1.5, result.Kill)                            // 150 / 100
	assert.Equal(t, 8000.0, result.Exp)                          // 800000 / 100
	assert.Equal(t, 1.75, result.PlatoonRate)                    // (50/100)*1 + (25/100)*2 + (25/100)*3
}

func TestNewShipStats_正常系_Warship情報なし(t *testing.T) {
	t.Parallel()

	shipID := ShipID(1)
	playerShipStats := PlayerShipStats{
		shipID: WGShipsStatsData{
			ShipID: shipID,
			Pvp: WGShipStatsValues{
				Battles:         50,
				Wins:            30,
				DamageDealt:     250000,
				MaxDamageDealt:  9000,
				Frags:           75,
				SurvivedWins:    15,
				SurvivedBattles: 25,
				Xp:              400000,
				PlanesKilled:    5,
			},
			PvpSolo: WGShipStatsValues{
				Battles: 50,
			},
		},
	}

	result := NewShipStats(
		StatsPatternPvPAll,
		playerShipStats,
		PlayerShipBadges{},
		shipID,
		Warships{},
	)

	assert.Equal(t, uint(50), result.Battles)
	assert.Equal(t, 5000.0, result.Damage.Value)      // 250000 / 50
	assert.Equal(t, RatingNone, result.Damage.Rating) // Warship ServerAverage がないため RatingNone
	assert.Equal(t, uint(9000), result.MaxDamage.Value)
	assert.Equal(t, 60.0, result.WinRate.Value) // (30 / 50) * 100
	assert.Equal(t, Rating(RatingUnicum), result.WinRate.Rating)
	assert.Equal(t, RatingNone, result.PR.Rating) // Warship ServerAverage がないため RatingNone
	assert.Equal(t, 1.5, result.Kill)             // 75 / 50
	assert.Equal(t, 8000.0, result.Exp)           // 400000 / 50
	assert.Equal(t, 1.0, result.PlatoonRate)      // (50/50)*1 = 1
}

func TestNewShipStats_複数パターンの統計(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		pattern     StatsPattern
		statsValues WGShipStatsValues
		expectWins  uint
	}{
		{
			name:    "PvPAll パターン",
			pattern: StatsPatternPvPAll,
			statsValues: WGShipStatsValues{
				Battles: 100,
				Wins:    60,
			},
			expectWins: 60,
		},
		{
			name:    "PvPSolo パターン",
			pattern: StatsPatternPvPSolo,
			statsValues: WGShipStatsValues{
				Battles: 50,
				Wins:    30,
			},
			expectWins: 30,
		},
		{
			name:    "RankSolo パターン",
			pattern: StatsPatternRankSolo,
			statsValues: WGShipStatsValues{
				Battles: 30,
				Wins:    20,
			},
			expectWins: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shipID := ShipID(1)
			playerShipStats := PlayerShipStats{
				shipID: WGShipsStatsData{
					ShipID:   shipID,
					Pvp:      WGShipStatsValues{Battles: 100, Wins: 60},
					PvpSolo:  WGShipStatsValues{Battles: 50, Wins: 30},
					RankSolo: tt.statsValues,
				},
			}

			result := NewShipStats(
				tt.pattern,
				playerShipStats,
				PlayerShipBadges{},
				shipID,
				Warships{},
			)

			expectedWinRate := (float64(tt.expectWins) / float64(tt.statsValues.Battles)) * 100
			assert.Equal(t, expectedWinRate, result.WinRate.Value)
		})
	}
}
