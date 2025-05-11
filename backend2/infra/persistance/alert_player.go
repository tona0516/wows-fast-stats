package persistence

type AlertPlayer struct {
	AccountID   int                `json:"account_id"`
	AccountName string             `json:"account_name"`
	Icon        AlertBootstrapIcon `json:"icon"`
	Message     string             `json:"message"`
}

type AlertPlayers []AlertPlayer
