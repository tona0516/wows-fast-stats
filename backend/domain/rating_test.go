package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRating_PersonalRating(t *testing.T) {
	t.Parallel()

	rating := Rating{}
	actual := RatingFactor{
		AvgDamage: 10000,
		AvgFrags:  2,
		WinRate:   60,
	}

	// 正常系
	assert.InDelta(t, 1875, rating.PersonalRating(actual, RatingFactor{
		AvgDamage: 8000,
		AvgFrags:  1,
		WinRate:   50,
	}), 0.1)
	assert.InDelta(t, 1150, rating.PersonalRating(actual, RatingFactor{
		AvgDamage: 10000,
		AvgFrags:  2,
		WinRate:   60,
	}), 0.1)

	// 異常系
	assert.Equal(t, float64(-1), rating.PersonalRating(actual, RatingFactor{
		AvgDamage: 0,
		AvgFrags:  1,
		WinRate:   50,
	}))

	assert.Equal(t, float64(-1), rating.PersonalRating(actual, RatingFactor{
		AvgDamage: 8000,
		AvgFrags:  0,
		WinRate:   50,
	}))

	assert.Equal(t, float64(-1), rating.PersonalRating(actual, RatingFactor{
		AvgDamage: 8000,
		AvgFrags:  1,
		WinRate:   0,
	}))
}
