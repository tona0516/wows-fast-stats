package infra

import "strconv"

type Env struct {
	AppName                string
	AppVer                 string
	AlertDiscordWebhookURL string
	InfoDiscordWebhookURL  string
	WGAppID                string
	IsDev                  bool
}

func NewEnv(
	appName string,
	appVer string,
	alertDiscordWebhookURL string,
	infoDiscordWebhookURL string,
	wgAppID string,
	isDevStr string,
) Env {
	isDev, _ := strconv.ParseBool(isDevStr)

	return Env{
		AppName:                appName,
		AppVer:                 appVer,
		AlertDiscordWebhookURL: alertDiscordWebhookURL,
		InfoDiscordWebhookURL:  infoDiscordWebhookURL,
		WGAppID:                wgAppID,
		IsDev:                  isDev,
	}
}
