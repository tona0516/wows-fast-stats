package core

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitedValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		value, max, min float64
		expected        float64
	}{
		{
			name:     "値がmax以下かつmin以上の場合は値そのものを返す",
			value:    50,
			max:      100,
			min:      0,
			expected: 50,
		},
		{
			name:     "値がmaxより大きい場合はmaxを返す",
			value:    150,
			max:      100,
			min:      0,
			expected: 100,
		},
		{
			name:     "値がminより小さい場合はminを返す",
			value:    -50,
			max:      100,
			min:      0,
			expected: 0,
		},
		{
			name:     "値がmaxと等しい場合はmaxを返す",
			value:    100,
			max:      100,
			min:      0,
			expected: 100,
		},
		{
			name:     "値がminと等しい場合はminを返す",
			value:    0,
			max:      100,
			min:      0,
			expected: 0,
		},
		{
			name:     "負の値を扱う場合",
			value:    -50,
			max:      -10,
			min:      -100,
			expected: -50,
		},
		{
			name:     "浮動小数点数を扱う場合",
			value:    50.5,
			max:      100.5,
			min:      0.5,
			expected: 50.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := limitedValue(tt.value, tt.max, tt.min)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFloorU4(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{
			name:     "整数値の場合",
			value:    123,
			expected: 123,
		},
		{
			name:     "4桁未満の小数点の場合は切り捨てられない",
			value:    123.456,
			expected: 123.456,
		},
		{
			name:     "4桁より多い小数点の場合は4桁で切り捨てられる",
			value:    123.456789,
			expected: 123.4567,
		},
		{
			name:     "ちょうど4桁の小数点の場合",
			value:    123.4567,
			expected: 123.4567,
		},
		{
			name:     "0の場合",
			value:    0,
			expected: 0,
		},
		{
			name:     "負の値の場合",
			value:    -123.456789,
			expected: -123.4568,
		},
		{
			name:     "0に近い値の場合",
			value:    0.00001,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := floorU4(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFloor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{
			name:     "整数値の場合",
			value:    123,
			expected: 123,
		},
		{
			name:     "正の小数点値の場合は小数点以下が切り捨てられる",
			value:    123.9999,
			expected: 123,
		},
		{
			name:     "負の小数点値の場合は小数点以下が切り捨てられる",
			value:    -123.1,
			expected: -124,
		},
		{
			name:     "0の場合",
			value:    0,
			expected: 0,
		},
		{
			name:     "0に近い正の値の場合",
			value:    0.5,
			expected: 0,
		},
		{
			name:     "0に近い負の値の場合",
			value:    -0.5,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := floor(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRound(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    float64
		places   int
		expected float64
	}{
		{
			name:     "1桁で四捨五入する場合",
			value:    123.456,
			places:   1,
			expected: 123.5,
		},
		{
			name:     "2桁で四捨五入する場合",
			value:    123.456,
			places:   2,
			expected: 123.46,
		},
		{
			name:     "0桁で四捨五入する場合は整数にされる",
			value:    123.456,
			places:   0,
			expected: 123,
		},
		{
			name:     "3桁で四捨五入する場合",
			value:    123.456,
			places:   3,
			expected: 123.456,
		},
		{
			name:     "4桁で四捨五入する場合",
			value:    123.456,
			places:   4,
			expected: 123.456,
		},
		{
			name:     "負の値の場合",
			value:    -123.456,
			places:   2,
			expected: -123.46,
		},
		{
			name:     "0の場合",
			value:    0,
			places:   2,
			expected: 0,
		},
		{
			name:     "境界値の場合（5で四捨五入）",
			value:    123.45,
			places:   1,
			expected: 123.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := round(tt.value, tt.places)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeDivide(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		numerator   float64
		denominator float64
		expected    float64
	}{
		{
			name:        "通常の除算",
			numerator:   10,
			denominator: 2,
			expected:    5,
		},
		{
			name:        "分母が0の場合は0を返す",
			numerator:   10,
			denominator: 0,
			expected:    0,
		},
		{
			name:        "分子が0の場合",
			numerator:   0,
			denominator: 10,
			expected:    0,
		},
		{
			name:        "浮動小数点数の除算",
			numerator:   10.5,
			denominator: 2.5,
			expected:    4.2,
		},
		{
			name:        "負の値の除算",
			numerator:   -10,
			denominator: 2,
			expected:    -5,
		},
		{
			name:        "分子分母両方負の場合",
			numerator:   -10,
			denominator: -2,
			expected:    5,
		},
		{
			name:        "1で割った場合",
			numerator:   123.456,
			denominator: 1,
			expected:    123.456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeDivide(tt.numerator, tt.denominator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMax(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{
			name:     "複数の値から最大値を取得する",
			values:   []float64{1, 5, 3, 2, 4},
			expected: 5,
		},
		{
			name:     "1つの値のみの場合",
			values:   []float64{42},
			expected: 42,
		},
		{
			name:     "空配列の場合は0を返す",
			values:   []float64{},
			expected: 0,
		},
		{
			name:     "負の値を含む場合",
			values:   []float64{-10, -5, -20, -1},
			expected: -1,
		},
		{
			name:     "0を含む場合",
			values:   []float64{-5, 0, 5},
			expected: 5,
		},
		{
			name:     "浮動小数点数の場合",
			values:   []float64{1.1, 1.5, 1.3, 1.2},
			expected: 1.5,
		},
		{
			name:     "すべて同じ値の場合",
			values:   []float64{5, 5, 5, 5},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := max(tt.values)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeometricMean(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{
			name:     "2つの値の幾何平均",
			values:   []float64{2, 8},
			expected: 4,
		},
		{
			name:     "3つの値の幾何平均",
			values:   []float64{1, 2, 4},
			expected: 2,
		},
		{
			name:     "1つの値の場合はその値を返す",
			values:   []float64{5},
			expected: 5,
		},
		{
			name:     "空配列の場合は0を返す",
			values:   []float64{},
			expected: 0,
		},
		{
			name:     "すべて同じ値の場合",
			values:   []float64{5, 5, 5},
			expected: 5,
		},
		{
			name:     "1を含む場合",
			values:   []float64{1, 1, 1},
			expected: 1,
		},
		{
			name:     "浮動小数点数の場合",
			values:   []float64{2.5, 5.5},
			expected: math.Sqrt(2.5 * 5.5),
		},
		{
			name:     "4つの値の幾何平均",
			values:   []float64{1, 2, 3, 4},
			expected: math.Pow(24, 0.25),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := geometricMean(tt.values)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}
