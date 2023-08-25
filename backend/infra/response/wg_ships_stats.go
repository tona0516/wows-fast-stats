package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGShipsStats struct {
	WGResponseCommon[domain.WGShipsStats]
}

func (w WGShipsStats) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGShipsStatsData{}).Elem())
}
