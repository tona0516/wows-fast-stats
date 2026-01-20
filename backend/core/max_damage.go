package core

type MaxDamage struct {
	ShipID   ShipID `json:"shipID"`
	ShipName string `json:"shipName"`
	ShipTier uint   `json:"shipTier"`
	Value    uint   `json:"value"`
}
