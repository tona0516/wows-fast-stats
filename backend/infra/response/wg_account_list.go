package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGAccountList struct {
	WGResponseCommon[data.WGAccountList]
}

func (w WGAccountList) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGAccountListData{}).Elem())
}
