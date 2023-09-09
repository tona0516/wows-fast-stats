package domain

type AllPlayerShipsStats map[int]WGShipsStats

func (w AllPlayerShipsStats) Player(accountID int) []WGShipsStatsData {
	return w[accountID][accountID]
}

type WGShipsStats map[int][]WGShipsStatsData

type WGShipsStatsData struct {
	Pvp     WGShipStatsValues `json:"pvp"`
	PvpSolo WGShipStatsValues `json:"pvp_solo"`
	ShipID  int               `json:"ship_id"`
}
