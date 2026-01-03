package data

import (
	"reflect"
	"wfs/internal/util"
)

type WGAccountInfo struct {
	WGResponseCommon[map[int]WGAccountInfoData]
}

func (w WGAccountInfo) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGAccountInfoData]())
}

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp      WGPlayerStatsValues `json:"pvp"`
		PvpSolo  WGPlayerStatsValues `json:"pvp_solo"`
		PvpDiv2  WGPlayerStatsValues `json:"pvp_div2"`
		PvpDiv3  WGPlayerStatsValues `json:"pvp_div3"`
		RankSolo WGPlayerStatsValues `json:"rank_solo"`
	} `json:"statistics"`
}
