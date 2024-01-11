package model

type AllPlayerShipsStats map[int]WGShipsStats

func (w AllPlayerShipsStats) Player(accountID int) []WGShipsStatsData {
	return w[accountID][accountID]
}

type WGShipsStats map[int][]WGShipsStatsData

type WGShipsStatsData struct {
	Pvp     WGShipStatsValues `json:"pvp"`
	PvpSolo WGShipStatsValues `json:"pvp_solo"`
	PvpDiv2 WGShipStatsValues `json:"pvp_div2"`
	PvpDiv3 WGShipStatsValues `json:"pvp_div3"`
	ShipID  int               `json:"ship_id"`
}
