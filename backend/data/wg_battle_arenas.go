package data

import (
	"reflect"
	"wfs/backend/util"
)

type WGBattleArenas struct {
	WGResponseCommon[map[int]WGBattleArenasData]
}

func (w WGBattleArenas) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGBattleArenasData]())
}

type WGBattleArenasData struct {
	Name string `json:"name"`
}
