package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGBattleArenas struct {
	WGResponseCommon[data.WGBattleArenas]
}

func (w WGBattleArenas) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGBattleArenasData{}).Elem())
}
