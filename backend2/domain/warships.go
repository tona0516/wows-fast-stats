package domain

type Warships map[int]Warship

type Warship struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Tier          uint     `json:"tier"`
	Type          ShipType `json:"type"`
	Nation        Nation   `json:"nation"`
	IsPremium     bool     `json:"is_premium"`
	AverageDamage float64  `json:"average_damage"`
	AverageFrags  float64  `json:"average_frags"`
	WinRate       float64  `json:"win_rate"`
}
