package data

type TierGroup struct {
	Low    float64 `json:"low"`    // tier 1~4
	Middle float64 `json:"middle"` // tier 5~7
	High   float64 `json:"high"`   // tier 8~★
}
