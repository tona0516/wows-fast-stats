package core

type RatingThreshold struct {
	Rating          Rating
	PR              float64
	ShipDamageRatio float64
	WinRate         float64
}

var ratingThresholds = []RatingThreshold{
	{Rating: RatingSuperUnicum, PR: 2450, ShipDamageRatio: 1.6, WinRate: 65},
	{Rating: RatingUnicum, PR: 2100, ShipDamageRatio: 1.5, WinRate: 60},
	{Rating: RatingGreat, PR: 1750, ShipDamageRatio: 1.4, WinRate: 56},
	{Rating: RatingVeryGood, PR: 1550, ShipDamageRatio: 1.2, WinRate: 54},
	{Rating: RatingGood, PR: 1350, ShipDamageRatio: 1.0, WinRate: 52},
	{Rating: RatingAvg, PR: 1100, ShipDamageRatio: 0.8, WinRate: 50},
	{Rating: RatingBelowAvg, PR: 750, ShipDamageRatio: 0.6, WinRate: 47},
	{Rating: RatingBad, PR: 0, ShipDamageRatio: 0, WinRate: 0},
}
