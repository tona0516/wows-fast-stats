package data

type AlertPattern string

func AlertPatterns() []string {
	return []string{
		"bi-check-circle-fill",
		"bi-exclamation-triangle-fill",
		"bi-patch-question-fill",
		"bi-1-square-fill",
		"bi-2-square-fill",
		"bi-3-square-fill",
	}
}

type AlertPlayer struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	Pattern   string `json:"pattern"`
	Message   string `json:"message"`
}
