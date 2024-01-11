package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGEncycInfo struct {
	WGResponseCommon[model.WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGEncycInfoData{}).Elem())
}
