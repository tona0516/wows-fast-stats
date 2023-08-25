package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGClansAccountInfo struct {
	WGResponseCommon[domain.WGClansAccountInfo]
}

func (w WGClansAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGClansAccountInfoData{}).Elem())
}
