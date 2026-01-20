package core

type Rating string

const (
	RatingSuperUnicum Rating = "super_unicum"
	RatingUnicum      Rating = "unicum"
	RatingGreat       Rating = "great"
	RatingVeryGood    Rating = "very_good"
	RatingGood        Rating = "good"
	RatingAvg         Rating = "avg"
	RatingBelowAvg    Rating = "below_avg"
	RatingBad         Rating = "bad"
	RatingNone        Rating = "none"
)

func NewDamageRatings(expect float64) []RatingValue {
	var result = make([]RatingValue, 0, len(ratingThresholds))
	for _, v := range ratingThresholds {
		value := expect * v.ShipDamageRatio
		result = append(result, RatingValue{
			Value:  value,
			Rating: v.Rating,
		})
	}

	return result
}

func NewRatingFromPR(value float64) Rating {
	for _, t := range ratingThresholds {
		if value >= t.PR {
			return t.Rating
		}
	}

	return RatingNone
}

func NewRatingFromWinRate(value float64) Rating {
	for _, t := range ratingThresholds {
		if value >= t.WinRate {
			return t.Rating
		}
	}

	return RatingNone
}

func NewRatingFromShipDamage(value, expected float64) Rating {
	ratio := safeDivide(value, expected)
	for _, t := range ratingThresholds {
		if ratio >= t.ShipDamageRatio {
			return t.Rating
		}
	}

	return RatingNone
}
