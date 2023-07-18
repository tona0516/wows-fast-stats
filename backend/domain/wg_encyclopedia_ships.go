package domain

import "reflect"

type WGEncycShips struct {
	Status string `json:"status"`
	Meta   struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
	Data  map[int]WGEncyclopediaShipsData `json:"data"`
	Error WGError                         `json:"error"`
}

func (w WGEncycShips) GetStatus() string {
	return w.Status
}

func (w WGEncycShips) GetError() WGError {
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
