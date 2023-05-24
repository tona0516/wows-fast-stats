package vo

import (
	"reflect"
	"sort"

	"github.com/samber/lo"
)

type WGAccountList struct {
	Status string              `json:"status"`
	Data   []WGAccountListData `json:"data"`
	Error  WGError             `json:"error"`
}

func (w WGAccountList) AccountIDs() []int {
	accountIDs := lo.FilterMap(w.Data, func(account WGAccountListData, _ int) (int, bool) {
		return account.AccountID, account.AccountID != 0
	})
	accountIDs = lo.Uniq(accountIDs)
	sort.Ints(accountIDs)
	return accountIDs
}

func (w WGAccountList) AccountID(nickname string) int {
	account, _ := lo.Find(w.Data, func(account WGAccountListData) bool {
		return account.NickName == nickname
	})
	return account.AccountID
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
