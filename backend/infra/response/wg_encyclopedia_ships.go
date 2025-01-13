package response

import (
	"reflect"
)

type WGEncycShips struct {
	WGResponseCommon[map[int]WGEncycShipsData]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

func (w WGEncycShips) Field() string {
	return fieldQuery(reflect.TypeOf(&WGEncycShipsData{}).Elem())
}

type WGEncycShipsData struct {
	Tier      uint   `json:"tier"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Nation    string `json:"nation"`
	IsPremium bool   `json:"is_premium"`
}
