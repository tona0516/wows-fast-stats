package persistence

type AlertBootstrapIcon string

func AlertBootstrapIcons() []AlertBootstrapIcon {
	return []AlertBootstrapIcon{
		"bi-check-circle-fill",
		"bi-exclamation-triangle-fill",
		"bi-patch-question-fill",
		"bi-1-square-fill",
		"bi-2-square-fill",
		"bi-3-square-fill",
	}
}
