package data

type MaxDamage struct {
	ShipID   int    `json:"ship_id"`
	ShipName string `json:"ship_name"`
	ShipTier uint   `json:"ship_tier"`
	Value    uint   `json:"value"`
}
