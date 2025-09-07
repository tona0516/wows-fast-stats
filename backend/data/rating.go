package data

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
	var result = make([]RatingValue, 0, len(thresholds))
	for _, v := range thresholds {
		value := expect * v.ShipDamageRatio
		result = append(result, RatingValue{
			Value:  value,
			Rating: v.Rating,
		})
	}

	return result
}

type RatingThreshold struct {
	Rating          Rating
	PR              float64
	ShipDamageRatio float64
	WinRate         float64
}

//nolint:gochecknoglobals
var thresholds = []RatingThreshold{
	{Rating: RatingSuperUnicum, PR: 2450, ShipDamageRatio: 1.6, WinRate: 65},
	{Rating: RatingUnicum, PR: 2100, ShipDamageRatio: 1.5, WinRate: 60},
	{Rating: RatingGreat, PR: 1750, ShipDamageRatio: 1.4, WinRate: 56},
	{Rating: RatingVeryGood, PR: 1550, ShipDamageRatio: 1.2, WinRate: 54},
	{Rating: RatingGood, PR: 1350, ShipDamageRatio: 1.0, WinRate: 52},
	{Rating: RatingAvg, PR: 1100, ShipDamageRatio: 0.8, WinRate: 50},
	{Rating: RatingBelowAvg, PR: 750, ShipDamageRatio: 0.6, WinRate: 47},
	{Rating: RatingBad, PR: 0, ShipDamageRatio: 0, WinRate: 0},
}

func NewRatingFromPR(value float64) Rating {
	for _, t := range thresholds {
		if value >= t.PR {
			return t.Rating
		}
	}

	return RatingNone
}

func NewRatingFromWinRate(value float64) Rating {
	for _, t := range thresholds {
		if value >= t.WinRate {
			return t.Rating
		}
	}

	return RatingNone
}

func NewRatingFromShipDamage(value, expected float64) Rating {
	if expected <= 0 {
		return RatingNone
	}

	ratio := value / expected
	for _, t := range thresholds {
		if ratio >= t.ShipDamageRatio {
			return t.Rating
		}
	}

	return RatingNone
}
