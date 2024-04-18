package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGAccountInfo struct {
	WGResponseCommon[data.WGAccountInfo]
}

func (w WGAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGAccountInfoData{}).Elem())
}
