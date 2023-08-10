package domain

import "reflect"

type WGAccountInfo struct {
	WGResponseCommon[map[int]WGAccountInfoData]
}

func (w WGAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&WGAccountInfoData{}).Elem())
}

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp     WGStatsValues `json:"pvp"`
		PvpSolo WGStatsValues `json:"pvp_solo"`
	} `json:"statistics"`
}
