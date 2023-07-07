package vo

type Meta struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
	OwnShip  string `json:"own_ship"`
}

type Battle struct {
	Meta  Meta   `json:"meta"`
	Teams []Team `json:"teams"`
}
