package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGBattleArenas struct {
	WGResponseCommon[domain.WGBattleArenas]
}

func (w WGBattleArenas) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGBattleArenasData{}).Elem())
}
