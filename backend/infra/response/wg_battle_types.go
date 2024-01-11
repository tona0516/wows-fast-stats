package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGBattleTypes struct {
	WGResponseCommon[model.WGBattleTypes]
}

func (w WGBattleTypes) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGBattleTypesData{}).Elem())
}
