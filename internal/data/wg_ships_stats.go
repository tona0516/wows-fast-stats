package data

import (
	"reflect"
	"wfs/internal/util"
)

type AllPlayerShipsStats map[int]WGShipsStats

func (w AllPlayerShipsStats) Player(accountID int) []WGShipsStatsData {
	return w[accountID].Data[accountID]
}

type WGShipsStats struct {
	WGResponseCommon[map[int][]WGShipsStatsData]
}

func (w WGShipsStats) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGShipsStatsData]())
}

type WGShipsStatsData struct {
	Pvp     WGShipStatsValues `json:"pvp"`
	PvpSolo WGShipStatsValues `json:"pvp_solo"`
	PvpDiv2 struct {
		Battles uint `json:"battles"`
	} `json:"pvp_div2"`
	PvpDiv3 struct {
		Battles uint `json:"battles"`
	} `json:"pvp_div3"`
	RankSolo WGShipStatsValues `json:"rank_solo"`
	ShipID   int               `json:"ship_id"`
}
