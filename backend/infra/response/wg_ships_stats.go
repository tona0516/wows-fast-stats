package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGShipsStats struct {
	WGResponseCommon[data.WGShipsStats]
}

func (w WGShipsStats) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGShipsStatsData{}).Elem())
}
