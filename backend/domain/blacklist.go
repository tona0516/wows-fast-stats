package domain

type BlackList []BlackListItem

type BlackListItem struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	Pattern   string `json:"pattern"`
	Message   string `json:"message"`
	CreatedAt uint64 `json:"created_at"`
}
