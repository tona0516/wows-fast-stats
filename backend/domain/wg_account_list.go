package domain

import (
	"reflect"
	"sort"

	"golang.org/x/exp/slices"
)

type WGAccountList struct {
	Status string              `json:"status"`
	Data   []WGAccountListData `json:"data"`
	Error  WGError             `json:"error"`
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

func (w WGAccountList) GetStatus() string {
	return w.Status
}

func (w WGAccountList) GetError() WGError {
	return w.Error
}

type WGAccountListData struct {
	NickName  string `json:"nickname"`
	AccountID int    `json:"account_id"`
}

func (w WGAccountListData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
