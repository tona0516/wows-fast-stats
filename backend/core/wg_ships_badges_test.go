package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerShipBadges_efficiencyBadge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		badges   PlayerShipBadges
		shipID   ShipID
		expected EfficiencyBadge
	}{
		{
			name: "TopGradeClass が1の場合、Expert バッジを返す",
			badges: PlayerShipBadges{
				ShipID(1): WGShipsBadgesData{
					ShipID:        ShipID(1),
					TopGradeClass: 1,
				},
			},
			shipID:   ShipID(1),
			expected: EfficiencyBadgeExpert,
		},
		{
			name: "TopGradeClass が2の場合、First バッジを返す",
			badges: PlayerShipBadges{
				ShipID(2): WGShipsBadgesData{
					ShipID:        ShipID(2),
					TopGradeClass: 2,
				},
			},
			shipID:   ShipID(2),
			expected: EfficiencyBadgeFirst,
		},
		{
			name: "TopGradeClass が3の場合、Second バッジを返す",
			badges: PlayerShipBadges{
				ShipID(3): WGShipsBadgesData{
					ShipID:        ShipID(3),
					TopGradeClass: 3,
				},
			},
			shipID:   ShipID(3),
			expected: EfficiencyBadgeSecond,
		},
		{
			name: "TopGradeClass が4の場合、Third バッジを返す",
			badges: PlayerShipBadges{
				ShipID(4): WGShipsBadgesData{
					ShipID:        ShipID(4),
					TopGradeClass: 4,
				},
			},
			shipID:   ShipID(4),
			expected: EfficiencyBadgeThird,
		},
		{
			name: "TopGradeClass が不正な値の場合、None を返す",
			badges: PlayerShipBadges{
				ShipID(5): WGShipsBadgesData{
					ShipID:        ShipID(5),
					TopGradeClass: 999,
				},
			},
			shipID:   ShipID(5),
			expected: EfficiencyBadgeNone,
		},
		{
			name:     "指定された shipID が存在しない場合、None を返す",
			badges:   PlayerShipBadges{},
			shipID:   ShipID(999),
			expected: EfficiencyBadgeNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.badges.efficiencyBadge(tt.shipID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlayerShipBadges_efficiencyBadges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		badges   PlayerShipBadges
		expected EfficiencyBadgeGroup
	}{
		{
			name:   "空の PlayerShipBadges の場合、すべてのカウントが0",
			badges: PlayerShipBadges{},
			expected: EfficiencyBadgeGroup{
				Expert: 0,
				First:  0,
				Second: 0,
				Third:  0,
			},
		},
		{
			name: "Expert バッジ1個のみの場合、Expert のカウントが1",
			badges: PlayerShipBadges{
				ShipID(1): WGShipsBadgesData{
					ShipID:        ShipID(1),
					TopGradeClass: 1,
				},
			},
			expected: EfficiencyBadgeGroup{
				Expert: 1,
				First:  0,
				Second: 0,
				Third:  0,
			},
		},
		{
			name: "各バッジが複数個ある場合、各カウントが正確に計算される",
			badges: PlayerShipBadges{
				ShipID(1): WGShipsBadgesData{
					ShipID:        ShipID(1),
					TopGradeClass: 1,
				},
				ShipID(2): WGShipsBadgesData{
					ShipID:        ShipID(2),
					TopGradeClass: 1,
				},
				ShipID(3): WGShipsBadgesData{
					ShipID:        ShipID(3),
					TopGradeClass: 2,
				},
				ShipID(4): WGShipsBadgesData{
					ShipID:        ShipID(4),
					TopGradeClass: 3,
				},
				ShipID(5): WGShipsBadgesData{
					ShipID:        ShipID(5),
					TopGradeClass: 3,
				},
				ShipID(6): WGShipsBadgesData{
					ShipID:        ShipID(6),
					TopGradeClass: 4,
				},
			},
			expected: EfficiencyBadgeGroup{
				Expert: 2,
				First:  1,
				Second: 2,
				Third:  1,
			},
		},
		{
			name: "不正な TopGradeClass は無視される",
			badges: PlayerShipBadges{
				ShipID(1): WGShipsBadgesData{
					ShipID:        ShipID(1),
					TopGradeClass: 1,
				},
				ShipID(2): WGShipsBadgesData{
					ShipID:        ShipID(2),
					TopGradeClass: 999,
				},
				ShipID(3): WGShipsBadgesData{
					ShipID:        ShipID(3),
					TopGradeClass: 0,
				},
			},
			expected: EfficiencyBadgeGroup{
				Expert: 1,
				First:  0,
				Second: 0,
				Third:  0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.badges.efficiencyBadges()
			assert.Equal(t, tt.expected, result)
		})
	}
}
