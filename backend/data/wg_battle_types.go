package data

import (
	"reflect"
	"wfs/backend/util"
)

type WGBattleTypes struct {
	WGResponseCommon[map[string]WGBattleTypesData]
}

func (w WGBattleTypes) Field() string {
	return util.FieldQuery(reflect.TypeOf(&WGBattleTypesData{}).Elem())
}

type WGBattleTypesData struct {
	Name string `json:"name"`
}
