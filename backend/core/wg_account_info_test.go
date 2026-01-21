package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGAccountInfoData_playerStatsValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern StatsPattern
		setup   func() WGAccountInfoData
		expect  WGPlayerStatsValues
	}{
		{
			name:    "StatsPatternPvPAllを指定した場合、Pvpの統計値を返す",
			pattern: StatsPatternPvPAll,
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.Pvp = WGPlayerStatsValues{
					Wins:    100,
					Battles: 500,
					Frags:   250,
				}
				data.Statistics.PvpSolo = WGPlayerStatsValues{
					Wins:    50,
					Battles: 200,
					Frags:   100,
				}
				return data
			},
			expect: WGPlayerStatsValues{
				Wins:    100,
				Battles: 500,
				Frags:   250,
			},
		},
		{
			name:    "StatsPatternPvPSoloを指定した場合、PvpSoloの統計値を返す",
			pattern: StatsPatternPvPSolo,
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.PvpSolo = WGPlayerStatsValues{
					Wins:    50,
					Battles: 200,
					Frags:   100,
				}
				return data
			},
			expect: WGPlayerStatsValues{
				Wins:    50,
				Battles: 200,
				Frags:   100,
			},
		},
		{
			name:    "StatsPatternRankSoloを指定した場合、RankSoloの統計値を返す",
			pattern: StatsPatternRankSolo,
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.RankSolo = WGPlayerStatsValues{
					Wins:    30,
					Battles: 100,
					Frags:   50,
				}
				return data
			},
			expect: WGPlayerStatsValues{
				Wins:    30,
				Battles: 100,
				Frags:   50,
			},
		},
		{
			name:    "不正なパターンを指定した場合、空の構造体を返す",
			pattern: StatsPattern("invalid"),
			setup: func() WGAccountInfoData {
				return WGAccountInfoData{}
			},
			expect: WGPlayerStatsValues{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setup()
			result := data.playerStatsValue(tt.pattern)

			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestWGAccountInfoData_platoonRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func() WGAccountInfoData
		expect float64
	}{
		{
			name: "ソロのみの場合、1.0を返す",
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.Pvp.Battles = 100
				data.Statistics.PvpSolo.Battles = 100
				data.Statistics.PvpDiv2.Battles = 0
				data.Statistics.PvpDiv3.Battles = 0
				return data
			},
			expect: 1.0,
		},
		{
			name: "2人分艦隊で100%の場合、2.0を返す",
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.Pvp.Battles = 100
				data.Statistics.PvpSolo.Battles = 0
				data.Statistics.PvpDiv2.Battles = 100
				data.Statistics.PvpDiv3.Battles = 0
				return data
			},
			expect: 2.0,
		},
		{
			name: "3人分艦隊で100%の場合、3.0を返す",
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.Pvp.Battles = 100
				data.Statistics.PvpSolo.Battles = 0
				data.Statistics.PvpDiv2.Battles = 0
				data.Statistics.PvpDiv3.Battles = 100
				return data
			},
			expect: 3.0,
		},
		{
			name: "混合の場合、加重平均を返す",
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.Pvp.Battles = 100
				data.Statistics.PvpSolo.Battles = 50 // 50%
				data.Statistics.PvpDiv2.Battles = 30 // 30%
				data.Statistics.PvpDiv3.Battles = 20 // 20%
				return data
			},
			expect: 0.5*1 + 0.3*2 + 0.2*3,
		},
		{
			name: "全てが0の場合、0を返す",
			setup: func() WGAccountInfoData {
				data := WGAccountInfoData{}
				data.Statistics.Pvp.Battles = 0
				data.Statistics.PvpSolo.Battles = 0
				data.Statistics.PvpDiv2.Battles = 0
				data.Statistics.PvpDiv3.Battles = 0
				return data
			},
			expect: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setup()
			result := data.platoonRate()

			assert.InDelta(t, tt.expect, result, 0.001)
		})
	}
}

func TestWGPlayerStatsValues_avgDamage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values WGPlayerStatsValues
		expect float64
	}{
		{
			name: "通常の計算が行われる",
			values: WGPlayerStatsValues{
				DamageDealt: 10000,
				Battles:     100,
			},
			expect: 100.0,
		},
		{
			name: "小数点を含む結果になる場合",
			values: WGPlayerStatsValues{
				DamageDealt: 10000,
				Battles:     3,
			},
			expect: 10000.0 / 3,
		},
		{
			name: "Battlesが0の場合、0を返す",
			values: WGPlayerStatsValues{
				DamageDealt: 10000,
				Battles:     0,
			},
			expect: 0.0,
		},
		{
			name: "DamageDealtが0の場合、0を返す",
			values: WGPlayerStatsValues{
				DamageDealt: 0,
				Battles:     100,
			},
			expect: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.avgDamage()

			assert.InDelta(t, tt.expect, result, 0.001)
		})
	}
}

func TestWGPlayerStatsValues_avgKills(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values WGPlayerStatsValues
		expect float64
	}{
		{
			name: "通常の計算が行われる",
			values: WGPlayerStatsValues{
				Frags:   500,
				Battles: 100,
			},
			expect: 5.0,
		},
		{
			name: "小数点を含む結果になる場合",
			values: WGPlayerStatsValues{
				Frags:   100,
				Battles: 3,
			},
			expect: 100.0 / 3,
		},
		{
			name: "Battlesが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Frags:   500,
				Battles: 0,
			},
			expect: 0.0,
		},
		{
			name: "Fragsが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Frags:   0,
				Battles: 100,
			},
			expect: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.avgKills()

			assert.InDelta(t, tt.expect, result, 0.001)
		})
	}
}

func TestWGPlayerStatsValues_winRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values WGPlayerStatsValues
		expect float64
	}{
		{
			name: "通常の計算が行われる",
			values: WGPlayerStatsValues{
				Wins:    50,
				Battles: 100,
			},
			expect: 50.0,
		},
		{
			name: "パーセンテージで返される",
			values: WGPlayerStatsValues{
				Wins:    1,
				Battles: 2,
			},
			expect: 50.0,
		},
		{
			name: "Battlesが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Wins:    50,
				Battles: 0,
			},
			expect: 0.0,
		},
		{
			name: "Winsが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Wins:    0,
				Battles: 100,
			},
			expect: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.winRate()

			assert.InDelta(t, tt.expect, result, 0.001)
		})
	}
}

func TestWGPlayerStatsValues_kdRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values WGPlayerStatsValues
		expect float64
	}{
		{
			name: "通常のK/D計算",
			values: WGPlayerStatsValues{
				Frags:           100,
				Battles:         100,
				SurvivedBattles: 50,
			},
			expect: 100.0 / 50,
		},
		{
			name: "Fragsが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Frags:           0,
				Battles:         100,
				SurvivedBattles: 50,
			},
			expect: 0.0,
		},
		{
			name: "デスが0の場合（全て生き残った場合）、Fragsを1で割る",
			values: WGPlayerStatsValues{
				Frags:           100,
				Battles:         50,
				SurvivedBattles: 50,
			},
			expect: 100.0 / 1,
		},
		{
			name: "デスが非常に少ない場合、最小値1を使用",
			values: WGPlayerStatsValues{
				Frags:           50,
				Battles:         100,
				SurvivedBattles: 99,
			},
			expect: 50.0 / 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.kdRate()

			assert.InDelta(t, tt.expect, result, 0.001)
		})
	}
}

func TestWGPlayerStatsValues_exp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values WGPlayerStatsValues
		expect float64
	}{
		{
			name: "通常の計算が行われる",
			values: WGPlayerStatsValues{
				Xp:      500000,
				Battles: 100,
			},
			expect: 5000.0,
		},
		{
			name: "小数点を含む結果になる場合",
			values: WGPlayerStatsValues{
				Xp:      10000,
				Battles: 3,
			},
			expect: 10000.0 / 3,
		},
		{
			name: "Battlesが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Xp:      500000,
				Battles: 0,
			},
			expect: 0.0,
		},
		{
			name: "Xpが0の場合、0を返す",
			values: WGPlayerStatsValues{
				Xp:      0,
				Battles: 100,
			},
			expect: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.exp()

			assert.InDelta(t, tt.expect, result, 0.001)
		})
	}
}

func TestWGPlayerStatsValues_survivedRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values WGPlayerStatsValues
		expect SurvivedRate
	}{
		{
			name: "通常の生存率計算",
			values: WGPlayerStatsValues{
				Battles:         100,
				SurvivedBattles: 40,
				Wins:            50,
				SurvivedWins:    30,
			},
			expect: SurvivedRate{
				All:  40.0,
				Win:  60.0,
				Lose: 20.0,
			},
		},
		{
			name: "全て生き残った場合",
			values: WGPlayerStatsValues{
				Battles:         100,
				SurvivedBattles: 100,
				Wins:            50,
				SurvivedWins:    50,
			},
			expect: SurvivedRate{
				All:  100.0,
				Win:  100.0,
				Lose: 100.0,
			},
		},
		{
			name: "全て死亡した場合",
			values: WGPlayerStatsValues{
				Battles:         100,
				SurvivedBattles: 0,
				Wins:            50,
				SurvivedWins:    0,
			},
			expect: SurvivedRate{
				All:  0.0,
				Win:  0.0,
				Lose: 0.0,
			},
		},
		{
			name: "敗北数が0の場合（全勝）",
			values: WGPlayerStatsValues{
				Battles:         50,
				SurvivedBattles: 25,
				Wins:            50,
				SurvivedWins:    25,
			},
			expect: SurvivedRate{
				All:  50.0,
				Win:  50.0,
				Lose: 0.0,
			},
		},
		{
			name: "Battlesが0の場合、全て0を返す",
			values: WGPlayerStatsValues{
				Battles:         0,
				SurvivedBattles: 0,
				Wins:            0,
				SurvivedWins:    0,
			},
			expect: SurvivedRate{
				All:  0.0,
				Win:  0.0,
				Lose: 0.0,
			},
		},
		{
			name: "Winsが0の場合（全敗）",
			values: WGPlayerStatsValues{
				Battles:         100,
				SurvivedBattles: 10,
				Wins:            0,
				SurvivedWins:    0,
			},
			expect: SurvivedRate{
				All:  10.0,
				Win:  0.0,
				Lose: 10.0,
			},
		},
		{
			name: "小数点を含む計算",
			values: WGPlayerStatsValues{
				Battles:         3,
				SurvivedBattles: 1,
				Wins:            1,
				SurvivedWins:    1,
			},
			expect: SurvivedRate{
				All:  100.0 / 3,
				Win:  100.0,
				Lose: 0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.survivedRate()

			assert.InDelta(t, tt.expect.All, result.All, 0.001)
			assert.InDelta(t, tt.expect.Win, result.Win, 0.001)
			assert.InDelta(t, tt.expect.Lose, result.Lose, 0.001)
		})
	}
}
