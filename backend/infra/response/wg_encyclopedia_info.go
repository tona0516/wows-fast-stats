package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGEncycInfo struct {
	WGResponseCommon[data.WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGEncycInfoData{}).Elem())
}
