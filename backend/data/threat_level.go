package data

type ThreatLevel struct {
	Rank     string  `json:"rank"`
	Raw      float64 `json:"raw"`
	Modified float64 `json:"modified"`
}
