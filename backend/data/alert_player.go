package data

type AlertPlayer struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	Pattern   string `json:"pattern"`
	Message   string `json:"message"`
}
