package domain

type Battle struct {
	Meta  Meta   `json:"meta"`
	Teams []Team `json:"teams"`
}

type Team struct {
	Players Players `json:"players"`
}

type Meta struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
	OwnShip  string `json:"own_ship"`
}
