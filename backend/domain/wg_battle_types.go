package domain

import "reflect"

type WGBattleTypes struct {
	WGResponseCommon[map[string]WGBattleTypesData]
}

func (w WGBattleTypes) Field() string {
	return fieldQuery(reflect.TypeOf(&WGBattleTypesData{}).Elem())
}

type WGBattleTypesData struct {
	Name string `json:"name"`
}
