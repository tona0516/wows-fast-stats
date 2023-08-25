package domain

type WGAccountInfo map[int]WGAccountInfoData

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp     WGStatsValues `json:"pvp"`
		PvpSolo WGStatsValues `json:"pvp_solo"`
	} `json:"statistics"`
}
