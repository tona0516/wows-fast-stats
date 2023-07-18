package domain

import "reflect"

type WGAccountInfo struct {
	Status string                    `json:"status"`
	Data   map[int]WGAccountInfoData `json:"data"`
	Error  WGError                   `json:"error"`
}

type WGAccountInfoData struct {
	HiddenProfile bool `json:"hidden_profile"`
	Statistics    struct {
		Pvp     WGStatsValues `json:"pvp"`
		PvpSolo WGStatsValues `json:"pvp_solo"`
	} `json:"statistics"`
}

func (w WGAccountInfo) GetStatus() string {
	return w.Status
}

func (w WGAccountInfo) GetError() WGError {
	return w.Error
}

func (w WGAccountInfoData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
