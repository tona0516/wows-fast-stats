package data

import (
	"reflect"
	"wfs/backend/util"
)

type WGBattleTypes struct {
	WGResponseCommon[map[string]WGBattleTypesData]
}

func (w WGBattleTypes) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGBattleTypesData]())
}

type WGBattleTypesData struct {
	Name string `json:"name"`
}
