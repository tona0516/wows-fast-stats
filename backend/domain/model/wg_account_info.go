package model

type WGAccountInfo map[int]WGAccountInfoData

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp     WGPlayerStatsValues `json:"pvp"`
		PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
		PvpDiv2 WGPlayerStatsValues `json:"pvp_div2"`
		PvpDiv3 WGPlayerStatsValues `json:"pvp_div3"`
	} `json:"statistics"`
}
