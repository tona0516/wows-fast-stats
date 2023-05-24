package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersonalRating(t *testing.T) {
	t.Parallel()

	rating := Rating{}

	// Test with valid inputs.
	actual := RatingFactor{
		AvgDamage: 10000,
		AvgFrags:  2,
		WinRate:   60,
	}

	expected := RatingFactor{
		AvgDamage: 8000,
		AvgFrags:  1,
		WinRate:   50,
	}
	assert.Equal(t, float64(1875), rating.PersonalRating(actual, expected))

	// Test with invalid inputs.
	expected = RatingFactor{
		AvgDamage: 0,
		AvgFrags:  1,
		WinRate:   50,
	}
	assert.Equal(t, float64(-1), rating.PersonalRating(actual, expected))

	expected = RatingFactor{
		AvgDamage: 8000,
		AvgFrags:  0,
		WinRate:   50,
	}
	assert.Equal(t, float64(-1), rating.PersonalRating(actual, expected))

	expected = RatingFactor{
		AvgDamage: 8000,
		AvgFrags:  1,
		WinRate:   0,
	}
	assert.Equal(t, float64(-1), rating.PersonalRating(actual, expected))
}
