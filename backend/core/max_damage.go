package core

type MaxDamage struct {
	ShipID   ShipID `json:"ship_id"`
	ShipName string `json:"ship_name"`
	ShipTier uint   `json:"ship_tier"`
	Value    uint   `json:"value"`
}
