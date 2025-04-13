package wargaming

import (
	"reflect"
)

type BattleTypesResponse struct {
	ResponseCommon[BattleTypes]
}

func (w BattleTypesResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&BattleTypesData{}).Elem())
}

type BattleTypes map[string]BattleTypesData

type BattleTypesData struct {
	Name string `json:"name"`
}
