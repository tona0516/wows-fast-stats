package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGClansInfo struct {
	WGResponseCommon[data.WGClansInfo]
}

func (w WGClansInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGClansInfoData{}).Elem())
}
