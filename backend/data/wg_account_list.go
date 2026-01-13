package data

import (
	"slices"
)

type WGAccountList struct {
	WGResponseCommon[[]WGAccountListData]
}

func (w WGAccountList) AccountIDs() []AccountID {
	accountIDs := make([]AccountID, 0)
	for _, v := range w.Data {
		if v.AccountID != 0 && !slices.Contains(accountIDs, v.AccountID) {
			accountIDs = append(accountIDs, v.AccountID)
		}
	}

	slices.Sort(accountIDs)
	return accountIDs
}

func (w WGAccountList) AccountID(nickname string) AccountID {
	for _, v := range w.Data {
		if v.NickName == nickname {
			return v.AccountID
		}
	}

	return 0
}

type WGAccountListData struct {
	NickName  string    `json:"nickname"`
	AccountID AccountID `json:"account_id"`
}
