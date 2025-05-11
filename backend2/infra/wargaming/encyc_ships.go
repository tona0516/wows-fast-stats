package wargaming

import (
	"reflect"
)

type EncycShips struct {
	ResponseCommon[map[int]EncycShipsData]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

func (w EncycShips) Field() string {
	return fieldQuery(reflect.TypeOf(&EncycShipsData{}).Elem())
}

type EncycShipsData struct {
	Tier      uint   `json:"tier"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Nation    string `json:"nation"`
	IsPremium bool   `json:"is_premium"`
}
