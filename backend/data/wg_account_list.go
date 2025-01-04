package data

import (
	"slices"
	"sort"
)

type WGAccountList []WGAccountListData

func (w WGAccountList) AccountIDs() []int {
	accountIDs := make([]int, 0)
	for _, v := range w {
		if v.AccountID != 0 && !slices.Contains(accountIDs, v.AccountID) {
			accountIDs = append(accountIDs, v.AccountID)
		}
	}

	sort.Ints(accountIDs)
	return accountIDs
}

func (w WGAccountList) AccountID(nickname string) int {
	for _, v := range w {
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
