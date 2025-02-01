package response

import (
	"reflect"
)

type WGAccountListResponse struct {
	WGResponseCommon[WGAccountList]
}

func (w WGAccountListResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&WGAccountListData{}).Elem())
}

type WGAccountList []WGAccountListData

type WGAccountListData struct {
	NickName  string `json:"nickname"`
	AccountID int    `json:"account_id"`
}
