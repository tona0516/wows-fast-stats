package domain

import "reflect"

type WGShipsStats struct {
	Status string                     `json:"status"`
	Data   map[int][]WGShipsStatsData `json:"data"`
	Error  WGError                    `json:"error"`
}

type WGShipsStatsData struct {
	Pvp     WGStatsValues `json:"pvp"`
	PvpSolo WGStatsValues `json:"pvp_solo"`
	ShipID  int           `json:"ship_id"`
}

func (w WGShipsStats) GetStatus() string {
	return w.Status
}

func (w WGShipsStats) GetError() WGError {
	return w.Error
}

func (w WGShipsStatsData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
