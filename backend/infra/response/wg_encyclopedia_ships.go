package response

import (
	"wfs/backend/data"
)

type WGEncycShips struct {
	WGResponse[data.WGEncycShips]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

func (w WGEncycShips) Field() string {
	return wgAPIField(&data.WGEncycShipsData{})
}
