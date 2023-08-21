package domain

import "reflect"

type WGBattleArenas struct {
	WGResponseCommon[map[int]WGBattleArenasData]
}

func (w WGBattleArenas) Field() string {
	return fieldQuery(reflect.TypeOf(&WGBattleArenasData{}).Elem())
}

type WGBattleArenasData struct {
	Name string `json:"name"`
}
