package vo

type WGAccountList struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data []struct {
		NickName  string `json:"nickname"`
		AccountID int    `json:"account_id"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}

func (w *WGAccountList) AccountIDs() []int {
	accountIDs := make([]int, 0)
	for i := range w.Data {
		accountID := w.Data[i].AccountID
		if accountID != 0 {
			accountIDs = append(accountIDs, w.Data[i].AccountID)
		}
	}
	return accountIDs
}

func (w *WGAccountList) AccountID(nickname string) int {
	for i := range w.Data {
		item := w.Data[i]
		if item.NickName == nickname {
			return item.AccountID
		}
	}
	return 0
}
