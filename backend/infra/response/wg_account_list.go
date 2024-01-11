package response

import (
	"reflect"
	"wfs/backend/domain/model"
)

type WGAccountList struct {
	WGResponseCommon[model.WGAccountList]
}

func (w WGAccountList) Field() string {
	return fieldQuery(reflect.TypeOf(&model.WGAccountListData{}).Elem())
}
