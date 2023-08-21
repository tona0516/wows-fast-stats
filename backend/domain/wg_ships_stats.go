package domain

import "reflect"

type WGShipsStats struct {
	WGResponseCommon[map[int][]WGShipsStatsData]
}

func (w WGShipsStats) Field() string {
	return fieldQuery(reflect.TypeOf(&WGShipsStatsData{}).Elem())
}

type WGShipsStatsData struct {
	Pvp     WGStatsValues `json:"pvp"`
	PvpSolo WGStatsValues `json:"pvp_solo"`
	ShipID  int           `json:"ship_id"`
}
