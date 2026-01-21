package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOverallStats_戦闘数が0の場合(t *testing.T) {
	t.Parallel()

	result := NewOverallStats(
		StatsPatternPvPAll,
		WGAccountInfoData{},
		PlayerShipStats{},
		PlayerShipBadges{},
		ShipID(1),
		AccountID(123),
		[]Vehicle{},
		Warships{},
	)

	assert.Equal(t, OverallStats{}, result)
}

func TestNewOverallStats_正常系(t *testing.T) {
	t.Parallel()

	accountInfoData := WGAccountInfoData{
		HiddenProfile: false,
	}
	accountInfoData.Statistics.Pvp = WGPlayerStatsValues{
		Battles:              100,
		Wins:                 60,
		DamageDealt:          500000,
		MaxDamageDealt:       10000,
		MaxDamageDealtShipID: 1,
		Frags:                150,
		SurvivedWins:         30,
		SurvivedBattles:      50,
		Xp:                   800000,
	}
	accountInfoData.Statistics.PvpSolo = WGPlayerStatsValues{
		Battles: 100,
	}

	playerShipStats := PlayerShipStats{
		1: WGShipsStatsData{
			ShipID: 1,
			Pvp: WGShipStatsValues{
				Battles:         50,
				Wins:            30,
				DamageDealt:     250000,
				MaxDamageDealt:  9000,
				Frags:           75,
				SurvivedWins:    15,
				SurvivedBattles: 25,
				Xp:              400000,
			},
		},
	}

	warships := Warships{
		1: {
			ID:   1,
			Name: "TestShip",
			Tier: 8,
			Type: ShipTypeBB,
		},
	}

	result := NewOverallStats(
		StatsPatternPvPAll,
		accountInfoData,
		playerShipStats,
		PlayerShipBadges{},
		ShipID(1),
		AccountID(123),
		[]Vehicle{},
		warships,
	)

	assert.Equal(t, uint(100), result.Battles)
	assert.Greater(t, result.Damage.Value, 0.0)
	assert.Equal(t, RatingNone, result.Damage.Rating)
	assert.Equal(t, uint(10000), result.MaxDamage.Value)
	assert.Equal(t, ShipID(1), result.MaxDamage.ShipID)
	assert.Equal(t, "TestShip", result.MaxDamage.ShipName)
	assert.Equal(t, uint(8), result.MaxDamage.ShipTier)
	assert.Greater(t, result.WinRate.Value, 0.0)
	assert.NotEqual(t, RatingNone, result.WinRate.Rating)
	assert.Greater(t, result.KdRate, 0.0)
	assert.Greater(t, result.Kill, 0.0)
	assert.Greater(t, result.Exp, 0.0)
	assert.Greater(t, result.SurvivedRate.All, 0.0)
}

func TestNewOverallStats_複数のStatsPattern(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		pattern        StatsPattern
		statsBattles   uint
		expectedBattle uint
	}{
		{
			name:           "PvP All",
			pattern:        StatsPatternPvPAll,
			statsBattles:   100,
			expectedBattle: 100,
		},
		{
			name:           "PvP Solo",
			pattern:        StatsPatternPvPSolo,
			statsBattles:   30,
			expectedBattle: 30,
		},
		{
			name:           "Rank Solo - 戦闘なし",
			pattern:        StatsPatternRankSolo,
			statsBattles:   0,
			expectedBattle: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountInfoData := WGAccountInfoData{
				HiddenProfile: false,
			}
			accountInfoData.Statistics.Pvp = WGPlayerStatsValues{
				Battles:              100,
				Wins:                 60,
				DamageDealt:          500000,
				MaxDamageDealt:       10000,
				MaxDamageDealtShipID: 1,
				Frags:                100,
				SurvivedWins:         20,
				SurvivedBattles:      40,
				Xp:                   800000,
			}
			accountInfoData.Statistics.PvpSolo = WGPlayerStatsValues{
				Battles:              30,
				Wins:                 18,
				DamageDealt:          150000,
				MaxDamageDealt:       5000,
				MaxDamageDealtShipID: 1,
				Frags:                45,
				SurvivedWins:         10,
				SurvivedBattles:      15,
				Xp:                   240000,
			}
			accountInfoData.Statistics.RankSolo = WGPlayerStatsValues{
				Battles: 0,
			}

			playerShipStats := PlayerShipStats{
				1: WGShipsStatsData{
					ShipID: 1,
					Pvp: WGShipStatsValues{
						Battles:         50,
						Wins:            30,
						DamageDealt:     250000,
						MaxDamageDealt:  5000,
						Frags:           75,
						SurvivedWins:    20,
						SurvivedBattles: 30,
						Xp:              400000,
					},
					PvpSolo: WGShipStatsValues{
						Battles: 30,
					},
				},
			}

			result := NewOverallStats(
				tt.pattern,
				accountInfoData,
				playerShipStats,
				PlayerShipBadges{},
				ShipID(1),
				AccountID(123),
				[]Vehicle{},
				Warships{
					1: {ID: 1, Name: "Ship1", Tier: 6},
				},
			)

			assert.Equal(t, tt.expectedBattle, result.Battles)
		})
	}
}
