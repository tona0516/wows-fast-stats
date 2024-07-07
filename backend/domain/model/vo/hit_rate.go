package vo

type HitRate struct {
	ValueObject2[Percentage, Percentage]
}

func NewHitRate(mainBattery Percentage, torps Percentage) HitRate {
	vo2 := ValueObject2[Percentage, Percentage]{mainBattery, torps}
	return HitRate{vo2}
}
