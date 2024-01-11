package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGShipsStats struct {
	WGResponseCommon[model.WGShipsStats]
}

func (w WGShipsStats) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGShipsStatsData{}).Elem())
}
