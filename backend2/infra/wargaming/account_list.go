package wargaming

import (
	"reflect"
)

type AccountListResponse struct {
	ResponseCommon[AccountList]
}

func (w AccountListResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&AccountListData{}).Elem())
}

type AccountList []AccountListData

type AccountListData struct {
	NickName  string `json:"nickname"`
	AccountID int    `json:"account_id"`
}
