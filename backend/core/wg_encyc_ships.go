package core

type WGEncycShips struct {
	WGResponseCommon[map[ShipID]WGEncycShipsData]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

type WGEncycShipsData struct {
	Tier      uint   `json:"tier"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Nation    string `json:"nation"`
	IsPremium bool   `json:"is_premium"`
}
