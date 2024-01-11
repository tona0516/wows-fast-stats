package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGBattleArenas struct {
	WGResponseCommon[model.WGBattleArenas]
}

func (w WGBattleArenas) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGBattleArenasData{}).Elem())
}
