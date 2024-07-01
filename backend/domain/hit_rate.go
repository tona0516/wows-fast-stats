package domain

type HitRate struct {
	ValueObject2[percentage, percentage]
}

func NewHitRate(mainBattery percentage, torps percentage) HitRate {
	vo2 := ValueObject2[percentage, percentage]{mainBattery, torps}
	return HitRate{vo2}
}
