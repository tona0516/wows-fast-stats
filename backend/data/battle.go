package data

type Battle struct {
	Meta  Meta   `json:"meta"`
	Teams []Team `json:"teams"`
}

type Team struct {
	Players Players `json:"players"`
}

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociation_degree"`
	Accuracy           float64 `json:"asccuracy"`
}

type Meta struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
	OwnShip  string `json:"own_ship"`
}
