package core

import (
	"math"

	"github.com/shopspring/decimal"
)

func limitedValue(value float64, max float64, min float64) float64 {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

func floorU4(value float64) float64 {
	pow := math.Pow(10, 4)
	return floor(value*pow) / pow
}

func floor(value float64) float64 {
	result, _ := decimal.NewFromFloat(value).Floor().Float64()
	return result
}

func round(value float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	return math.Round(value*pow) / pow
}

type Number interface {
	float64 | int | uint
}

func safeDivide[T, U Number](numerator T, denominator U) float64 {
	if denominator == 0 {
		return 0
	}

	return float64(numerator) / float64(denominator)
}
