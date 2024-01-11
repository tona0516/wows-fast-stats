package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGAccountInfo struct {
	WGResponseCommon[model.WGAccountInfo]
}

func (w WGAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGAccountInfoData{}).Elem())
}
