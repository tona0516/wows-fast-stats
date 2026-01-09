package data

import (
	"reflect"
	"slices"
	"sort"
	"wfs/backend/util"
)

type WGAccountList struct {
	WGResponseCommon[[]WGAccountListData]
}

func (w WGAccountList) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGAccountListData]())
}

func (w WGAccountList) AccountIDs() []int {
	accountIDs := make([]int, 0)
	for _, v := range w.Data {
		if v.AccountID != 0 && !slices.Contains(accountIDs, v.AccountID) {
			accountIDs = append(accountIDs, v.AccountID)
		}
	}

	sort.Ints(accountIDs)
	return accountIDs
}

func (w WGAccountList) AccountID(nickname string) int {
	for _, v := range w.Data {
		if v.NickName == nickname {
			return v.AccountID
		}
	}

	return 0
}

type WGAccountListData struct {
	NickName  string `json:"nickname"`
	AccountID int    `json:"account_id"`
}
