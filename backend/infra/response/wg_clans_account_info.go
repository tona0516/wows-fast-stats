package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGClansAccountInfo struct {
	WGResponseCommon[model.WGClansAccountInfo]
}

func (w WGClansAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGClansAccountInfoData{}).Elem())
}
