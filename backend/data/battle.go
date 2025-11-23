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
	Players       Players           `json:"players"`
	PvPSolo       TeamStats         `json:"pvp_solo"`
	PvPAll        TeamStats         `json:"pvp_all"`
	RankSolo      TeamStats         `json:"rank_solo"`
	ShipTypeStats TeamShipTypeStats `json:"ship_type_stats"`
}

type TeamStats struct {
	TeamThreatLevel TeamThreatLevel `json:"team_threat_level"`
}

type TeamAverageStats struct {
	ShipAvgPR        float64 `json:"ship_avg_pr"`
	ShipAvgDamage    float64 `json:"ship_avg_damage"`
	ShipWinRate      float64 `json:"ship_win_rate"`
	OverallAvgPR     float64 `json:"overall_avg_pr"`
	OverallAvgDamage float64 `json:"overall_avg_damage"`
	OverallWinRate   float64 `json:"overall_win_rate"`
}

type TeamShipTypeStats struct {
	CV TeamAverageStats `json:"cv"`
	BB TeamAverageStats `json:"bb"`
	CL TeamAverageStats `json:"cl"`
	DD TeamAverageStats `json:"dd"`
	SS TeamAverageStats `json:"ss"`
}

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociation_degree"`
	Accuracy           float64 `json:"accuracy"`
}
