package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGEncycInfo struct {
	WGResponseCommon[domain.WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGEncycInfoData{}).Elem())
}
