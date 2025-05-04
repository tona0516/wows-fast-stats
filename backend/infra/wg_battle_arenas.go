package infra

import (
	"reflect"
)

type WGBattleArenasResponse struct {
	WGResponseCommon[WGBattleArenas]
}

func (w WGBattleArenasResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&WGBattleArenasData{}).Elem())
}

type WGBattleArenas map[int]WGBattleArenasData

type WGBattleArenasData struct {
	Name string `json:"name"`
}
