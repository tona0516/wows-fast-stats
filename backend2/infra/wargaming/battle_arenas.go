package wargaming

import (
	"reflect"
)

type BattleArenasResponse struct {
	ResponseCommon[BattleArenas]
}

func (w BattleArenasResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&BattleArenasData{}).Elem())
}

type BattleArenas map[int]BattleArenasData

type BattleArenasData struct {
	Name string `json:"name"`
}
