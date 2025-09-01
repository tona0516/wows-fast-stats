package data

type Battle struct {
	Meta  BattleMetaData `json:"metadata"`
	Teams []Team         `json:"teams"`
}

type BattleMetaData struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
}

type Team struct {
	Players Players `json:"players"`
}

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociation_degree"`
	Accuracy           float64 `json:"asccuracy"`
}
