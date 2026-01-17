package core

type RatingValue struct {
	Value  float64 `json:"value"`
	Rating Rating  `json:"rating"`
}

func NewRatingValue(value float64, rating Rating) RatingValue {
	return RatingValue{
		Value:  value,
		Rating: rating,
	}
}

func NewRatingValueNone() RatingValue {
	return RatingValue{
		Value:  -1,
		Rating: RatingNone,
	}
}
