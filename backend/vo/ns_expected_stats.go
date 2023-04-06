package vo

type NSExpectedStats struct {
	Time int                         `json:"time"`
	Data map[int]NSExpectedStatsData `json:"data"`
}

type NSExpectedStatsData struct {
	AverageDamageDealt float64 `json:"average_damage_dealt"`
	AverageFrags       float64 `json:"average_frags"`
	WinRate            float64 `json:"win_rate"`
}
