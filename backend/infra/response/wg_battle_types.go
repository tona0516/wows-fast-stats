package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGBattleTypes struct {
	WGResponseCommon[data.WGBattleTypes]
}

func (w WGBattleTypes) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGBattleTypesData{}).Elem())
}
