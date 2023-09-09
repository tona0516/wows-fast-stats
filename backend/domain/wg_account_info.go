package domain

type WGAccountInfo map[int]WGAccountInfoData

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp     WGPlayerStatsValues `json:"pvp"`
		PvpSolo WGPlayerStatsValues `json:"pvp_solo"`
	} `json:"statistics"`
}
