package response

import (
	"wfs/backend/data"
)

type WGShipsStats struct {
	WGResponse[data.WGShipsStats]
}

func (w WGShipsStats) Field() string {
	return wgAPIField(&data.WGShipsStatsData{})
}
