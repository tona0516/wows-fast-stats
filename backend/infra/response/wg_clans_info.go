package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGClansInfo struct {
	WGResponseCommon[domain.WGClansInfo]
}

func (w WGClansInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGClansInfoData{}).Elem())
}
