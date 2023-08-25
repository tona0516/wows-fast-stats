package response

import (
	"reflect"
	"wfs/backend/domain"
)

type WGAccountList struct {
	WGResponseCommon[domain.WGAccountList]
}

func (w WGAccountList) Field() string {
	return fieldQuery(reflect.TypeOf(&domain.WGAccountListData{}).Elem())
}
