package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGClansInfo struct {
	WGResponseCommon[model.WGClansInfo]
}

func (w WGClansInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGClansInfoData{}).Elem())
}
