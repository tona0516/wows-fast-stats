package vo

import "reflect"

type WGEncyclopediaShips struct {
	Status string `json:"status"`
	Meta   struct {
		Count     int `json:"count"`
		PageTotal int `json:"page_total"`
		Total     int `json:"total"`
		Limit     int `json:"limit"`
		Page      int `json:"page"`
	} `json:"meta"`
	Data  map[int]WGEncyclopediaShipsData `json:"data"`
	Error WGError                         `json:"error"`
}

func (w WGEncyclopediaShips) GetStatus() string {
	return w.Status
}

func (w WGEncyclopediaShips) GetError() WGError {
	return w.Error
}

type WGEncyclopediaShipsData struct {
	Tier   uint   `json:"tier"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Nation string `json:"nation"`
}

func (w WGEncyclopediaShipsData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
