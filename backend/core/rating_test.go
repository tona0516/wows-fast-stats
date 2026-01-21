package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDamageRatings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		expect float64
		want   int
	}{
		{
			name:   "期待値が正数の場合、すべてのレーティング閾値分の要素が返される",
			expect: 1000,
			want:   len(ratingThresholds),
		},
		{
			name:   "期待値が0の場合、すべてのレーティング閾値分の要素が返される",
			expect: 0,
			want:   len(ratingThresholds),
		},
		{
			name:   "期待値が負数の場合、すべてのレーティング閾値分の要素が返される",
			expect: -1000,
			want:   len(ratingThresholds),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDamageRatings(tt.expect)

			assert.Equal(t, tt.want, len(result))

			// 各要素がRatingValue構造体で正しく構成されていることを確認
			for i, ratingValue := range result {
				assert.Equal(t, tt.expect*ratingThresholds[i].ShipDamageRatio, ratingValue.Value)
				assert.Equal(t, ratingThresholds[i].Rating, ratingValue.Rating)
			}
		})
	}
}

func TestNewDamageRatings_計算結果の妥当性(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		expect float64
		index  int
		want   float64
	}{
		{
			name:   "SuperUnicumの計算が正しい",
			expect: 1000,
			index:  0,
			want:   1600, // 1000 * 1.6
		},
		{
			name:   "Unicumの計算が正しい",
			expect: 2000,
			index:  1,
			want:   3000, // 2000 * 1.5
		},
		{
			name:   "Badの計算が正しい",
			expect: 500,
			index:  len(ratingThresholds) - 1,
			want:   0, // 500 * 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDamageRatings(tt.expect)

			assert.Equal(t, tt.want, result[tt.index].Value)
		})
	}
}

func TestNewRatingFromPR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float64
		want  Rating
	}{
		{
			name:  "2450以上の場合、RatingSuperUnicumが返される",
			value: 2500,
			want:  RatingSuperUnicum,
		},
		{
			name:  "2450ちょうどの場合、RatingSuperUnicumが返される",
			value: 2450,
			want:  RatingSuperUnicum,
		},
		{
			name:  "2100以上2450未満の場合、RatingUnicumが返される",
			value: 2200,
			want:  RatingUnicum,
		},
		{
			name:  "1750以上2100未満の場合、RatingGreatが返される",
			value: 1800,
			want:  RatingGreat,
		},
		{
			name:  "1550以上1750未満の場合、RatingVeryGoodが返される",
			value: 1600,
			want:  RatingVeryGood,
		},
		{
			name:  "1350以上1550未満の場合、RatingGoodが返される",
			value: 1400,
			want:  RatingGood,
		},
		{
			name:  "1100以上1350未満の場合、RatingAvgが返される",
			value: 1200,
			want:  RatingAvg,
		},
		{
			name:  "750以上1100未満の場合、RatingBelowAvgが返される",
			value: 800,
			want:  RatingBelowAvg,
		},
		{
			name:  "0以上750未満の場合、RatingBadが返される",
			value: 100,
			want:  RatingBad,
		},
		{
			name:  "0ちょうどの場合、RatingBadが返される",
			value: 0,
			want:  RatingBad,
		},
		{
			name:  "負数の場合、RatingNoneが返される",
			value: -100,
			want:  RatingNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewRatingFromPR(tt.value)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNewRatingFromWinRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float64
		want  Rating
	}{
		{
			name:  "65以上の場合、RatingSuperUnicumが返される",
			value: 70,
			want:  RatingSuperUnicum,
		},
		{
			name:  "65ちょうどの場合、RatingSuperUnicumが返される",
			value: 65,
			want:  RatingSuperUnicum,
		},
		{
			name:  "60以上65未満の場合、RatingUnicumが返される",
			value: 62,
			want:  RatingUnicum,
		},
		{
			name:  "56以上60未満の場合、RatingGreatが返される",
			value: 58,
			want:  RatingGreat,
		},
		{
			name:  "54以上56未満の場合、RatingVeryGoodが返される",
			value: 55,
			want:  RatingVeryGood,
		},
		{
			name:  "52以上54未満の場合、RatingGoodが返される",
			value: 53,
			want:  RatingGood,
		},
		{
			name:  "50以上52未満の場合、RatingAvgが返される",
			value: 51,
			want:  RatingAvg,
		},
		{
			name:  "47以上50未満の場合、RatingBelowAvgが返される",
			value: 48,
			want:  RatingBelowAvg,
		},
		{
			name:  "0以上47未満の場合、RatingBadが返される",
			value: 30,
			want:  RatingBad,
		},
		{
			name:  "0ちょうどの場合、RatingBadが返される",
			value: 0,
			want:  RatingBad,
		},
		{
			name:  "負数の場合、RatingNoneが返される",
			value: -10,
			want:  RatingNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewRatingFromWinRate(tt.value)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNewRatingFromShipDamage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    float64
		expected float64
		want     Rating
	}{
		{
			name:     "比率が1.6以上の場合、RatingSuperUnicumが返される",
			value:    1600,
			expected: 1000,
			want:     RatingSuperUnicum,
		},
		{
			name:     "比率が1.6ちょうどの場合、RatingSuperUnicumが返される",
			value:    1600,
			expected: 1000,
			want:     RatingSuperUnicum,
		},
		{
			name:     "比率が1.5以上1.6未満の場合、RatingUnicumが返される",
			value:    1500,
			expected: 1000,
			want:     RatingUnicum,
		},
		{
			name:     "比率が1.4以上1.5未満の場合、RatingGreatが返される",
			value:    1400,
			expected: 1000,
			want:     RatingGreat,
		},
		{
			name:     "比率が1.2以上1.4未満の場合、RatingVeryGoodが返される",
			value:    1200,
			expected: 1000,
			want:     RatingVeryGood,
		},
		{
			name:     "比率が1.0以上1.2未満の場合、RatingGoodが返される",
			value:    1000,
			expected: 1000,
			want:     RatingGood,
		},
		{
			name:     "比率が0.8以上1.0未満の場合、RatingAvgが返される",
			value:    900,
			expected: 1000,
			want:     RatingAvg,
		},
		{
			name:     "比率が0.6以上0.8未満の場合、RatingBelowAvgが返される",
			value:    700,
			expected: 1000,
			want:     RatingBelowAvg,
		},
		{
			name:     "比率が0以上0.6未満の場合、RatingBadが返される",
			value:    500,
			expected: 1000,
			want:     RatingBad,
		},
		{
			name:     "比率が0ちょうどの場合、RatingBadが返される",
			value:    0,
			expected: 1000,
			want:     RatingBad,
		},
		{
			name:     "expectedが0の場合、safeDivideが0を返しRatingBadが返される",
			value:    100,
			expected: 0,
			want:     RatingBad,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewRatingFromShipDamage(tt.value, tt.expected)

			assert.Equal(t, tt.want, result)
		})
	}
}
