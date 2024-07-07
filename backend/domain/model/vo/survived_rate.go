package vo

type SurvivedRate struct {
	ValueObject2[Percentage, Percentage]
}

func NewSurvivedRate(win Percentage, lose Percentage) SurvivedRate {
	vo2 := ValueObject2[Percentage, Percentage]{win, lose}
	return SurvivedRate{vo2}
}
