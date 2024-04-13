package data

type WGEncycShips map[int]WGEncycShipsData

type WGEncycShipsData struct {
	Tier      uint   `json:"tier"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Nation    string `json:"nation"`
	IsPremium bool   `json:"is_premium"`
}
