package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGEncycShips struct {
	WGResponseCommon[model.WGEncycShips]
	Meta struct {
		PageTotal int `json:"page_total"`
		Page      int `json:"page"`
	} `json:"meta"`
}

func (w WGEncycShips) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGEncycShipsData{}).Elem())
}
