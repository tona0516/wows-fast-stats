package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGEncycShips struct {
	WGResponseCommon[data.WGEncycShips]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

func (w WGEncycShips) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGEncycShipsData{}).Elem())
}
