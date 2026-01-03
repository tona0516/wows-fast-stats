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
	Players  Players   `json:"players"`
	PvPSolo  TeamStats `json:"pvp_solo"`
	PvPAll   TeamStats `json:"pvp_all"`
	RankSolo TeamStats `json:"rank_solo"`
}

type TeamStats struct {
	TeamAverageStats TeamAverageStats `json:"team_average_stats"`
	TeamThreatLevel  TeamThreatLevel  `json:"team_threat_level"`
}

type TeamAverageStats struct {
	ShipPR         float64 `json:"ship_pr"`
	ShipDamage     float64 `json:"ship_damage"`
	ShipWinRate    float64 `json:"ship_win_rate"`
	ShipBattles    uint    `json:"ship_battles"`
	OverallPR      float64 `json:"overall_pr"`
	OverallDamage  float64 `json:"overall_damage"`
	OverallWinRate float64 `json:"overall_win_rate"`
	OverallBattles uint    `json:"overall_battles"`
}

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociation_degree"`
	Accuracy           float64 `json:"accuracy"`
}
