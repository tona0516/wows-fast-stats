package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGAccountInfo struct {
	WGResponseCommon[domain.WGAccountInfo]
}

func (w WGAccountInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGAccountInfoData{}).Elem())
}
