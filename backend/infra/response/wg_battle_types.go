package response

import (
	"reflect"
)

type WGBattleTypesResponse struct {
	WGResponseCommon[WGBattleTypes]
}

func (w WGBattleTypesResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&WGBattleTypesData{}).Elem())
}

type WGBattleTypes map[string]WGBattleTypesData

type WGBattleTypesData struct {
	Name string `json:"name"`
}
