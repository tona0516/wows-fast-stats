package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTeamThreatLevel_プレイヤーなし(t *testing.T) {
	t.Parallel()

	players := Players{}
	result := NewTeamThreatLevel(players, StatsPatternPvPAll)

	assert.Equal(t, TeamThreatLevel{}, result)
}

func TestNewTeamThreatLevel_全員隠しプロフィール(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(1),
				IsHidden: true,
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(2),
				IsHidden: true,
			},
		},
	}

	result := NewTeamThreatLevel(players, StatsPatternPvPAll)

	assert.Equal(t, TeamThreatLevel{}, result)
}

func TestNewTeamThreatLevel_全員IDが0(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(0),
				IsHidden: false,
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(0),
				IsHidden: false,
			},
		},
	}

	result := NewTeamThreatLevel(players, StatsPatternPvPAll)

	assert.Equal(t, TeamThreatLevel{}, result)
}

func TestNewTeamThreatLevel_PvPSolo_単一プレイヤー(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(1),
				IsHidden: false,
			},
			PvPSolo: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 1000,
					},
				},
			},
		},
	}

	result := NewTeamThreatLevel(players, StatsPatternPvPSolo)

	assert.InDelta(t, 1000.0, result.Average, 0.1)
	assert.InDelta(t, 0.0, result.DissociationDegree, 0.1)
	assert.Equal(t, 100.0, result.Accuracy)
}

func TestNewTeamThreatLevel_RankSolo_複数プレイヤー_異なるスコア(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(1),
				IsHidden: false,
			},
			RankSolo: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 1000,
					},
				},
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(2),
				IsHidden: false,
			},
			RankSolo: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 2000,
					},
				},
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(3),
				IsHidden: false,
			},
			RankSolo: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 1500,
					},
				},
			},
		},
	}

	result := NewTeamThreatLevel(players, StatsPatternRankSolo)

	// 幾何平均: (1000 * 2000 * 1500) ^ (1/3) ≈ 1442.25
	assert.InDelta(t, 1442.25, result.Average, 0.1)
	// dissociation degree: (max / mean - 1) * 100
	// (2000 / 1442.25 - 1) * 100 ≈ 38.65
	assert.InDelta(t, 38.65, result.DissociationDegree, 0.1)
	assert.Equal(t, 100.0, result.Accuracy)
}

func TestNewTeamThreatLevel_混在_隠しプロフィール_ID0_有効(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(0),
				IsHidden: false,
			},
			PvPAll: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 1000,
					},
				},
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(1),
				IsHidden: true,
			},
			PvPAll: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 2000,
					},
				},
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(2),
				IsHidden: false,
			},
			PvPAll: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 3000,
					},
				},
			},
		},
	}

	result := NewTeamThreatLevel(players, StatsPatternPvPAll)

	// 有効なプレイヤーは1人（ID=2のみ）で、ID=0と隠しプロフィールはスキップ
	assert.InDelta(t, 3000.0, result.Average, 0.1)
	assert.InDelta(t, 0.0, result.DissociationDegree, 0.1)
	// accuracy: 1 / 3 * 100 ≈ 33.33
	assert.InDelta(t, 33.33, result.Accuracy, 0.1)
}

func TestNewTeamThreatLevel_複数パターン_PvPSolo(t *testing.T) {
	t.Parallel()

	players := Players{
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(1),
				IsHidden: false,
			},
			PvPSolo: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 500,
					},
				},
			},
			PvPAll: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 2000,
					},
				},
			},
		},
		{
			PlayerInfo: PlayerInfo{
				ID:       AccountID(2),
				IsHidden: false,
			},
			PvPSolo: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 600,
					},
				},
			},
			PvPAll: PlayerStats{
				OverallStats: OverallStats{
					ThreatLevel: ThreatLevel{
						Modified: 1500,
					},
				},
			},
		},
	}

	result := NewTeamThreatLevel(players, StatsPatternPvPSolo)

	// PvPSolo: (500 * 600) ^ (1/2) ≈ 547.72
	assert.InDelta(t, 547.72, result.Average, 0.1)
	// dissociation: (600 / 547.72 - 1) * 100 ≈ 9.53
	assert.InDelta(t, 9.53, result.DissociationDegree, 0.1)
	assert.Equal(t, 100.0, result.Accuracy)
}
