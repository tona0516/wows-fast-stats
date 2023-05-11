package vo

import "reflect"

type WGAccountList struct {
	Status string              `json:"status"`
	Data   []WGAccountListData `json:"data"`
	Error  WGError             `json:"error"`
}

func (w WGAccountList) AccountIDs() []int {
	accountIDs := make([]int, 0)
	for i := range w.Data {
		accountID := w.Data[i].AccountID
		if accountID != 0 {
			accountIDs = append(accountIDs, w.Data[i].AccountID)
		}
	}

	return accountIDs
}

func (w WGAccountList) AccountID(nickname string) int {
	for i := range w.Data {
		item := w.Data[i]
		if item.NickName == nickname {
			return item.AccountID
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
