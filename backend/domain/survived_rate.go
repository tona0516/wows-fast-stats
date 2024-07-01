package domain

type SurvivedRate struct {
	ValueObject2[percentage, percentage]
}

func NewSurvivedRate(win percentage, lose percentage) SurvivedRate {
	vo2 := ValueObject2[percentage, percentage]{win, lose}
	return SurvivedRate{vo2}
}
