package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGBattleTypes struct {
	WGResponseCommon[domain.WGBattleTypes]
}

func (w WGBattleTypes) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGBattleTypesData{}).Elem())
}
