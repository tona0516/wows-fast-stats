package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGEncycShips struct {
	WGResponseCommon[domain.WGEncycShips]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

func (w WGEncycShips) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGEncycShipsData{}).Elem())
}
