package vo

type TierGroup[T any] struct {
    Low T `json:"low"` // tier 1~4
    Middle T `json:"middle"`// tier 5~7
    High T `json:"high"`// tier 8~★
}
