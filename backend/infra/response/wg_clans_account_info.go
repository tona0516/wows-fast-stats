package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGClansAccountInfo struct {
	WGResponseCommon[data.WGClansAccountInfo]
}

func (w WGClansAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGClansAccountInfoData{}).Elem())
}
