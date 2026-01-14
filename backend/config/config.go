//nolint:lll
package config

import (
	"time"

	"github.com/rs/zerolog"
)

//nolint:gochecknoglobals
var discordWebhookURL = "https://discord.com/api/webhooks"

type Config struct {
	Basic           BasicConfig
	WargamingClient WargamingConfig
	ClanClient      ClanConfig
	NumbersClient   NumbersClientConfig
	GithubClient    GithubClientConfig
	DiscordClient   DiscordClientConfig
	LocalFile       LocalFileConfig
	Logger          LoggerConfig
}

type BasicConfig struct {
	Name            string
	Version         string
	Width           int
	Height          int
	MinWidth        int
	MinHeight       int
	PollingInterval time.Duration
}

type WargamingConfig struct {
	URL          string
	RetryCount   int
	Timeout      time.Duration
	RateLimitRPS int
	AppID        string
}

type ClanConfig struct {
	URL        string
	RetryCount int
	Timeout    time.Duration
}

type NumbersClientConfig struct {
	URL        string
	RetryCount int
	Timeout    time.Duration
}

type GithubClientConfig struct {
	URL        string
	RetryCount int
	Timeout    time.Duration
}

type DiscordClientConfig struct {
	AlertWebhookURL string
	InfoWebhookURL  string
	RetryCount      int
	Timeout         time.Duration
}

type LocalFileConfig struct {
	RootDir string
}

type LoggerConfig struct {
	Level zerolog.Level
}

func NewConfig(appName, version, env string) Config {
	if env == "prod" {
		return newProdConfig(appName, version)
	} else {
		return NewDevConfig(appName, version)
	}
}

func newProdConfig(appName, version string) Config {
	return Config{
		Basic: BasicConfig{
			Name:            appName,
			Version:         version,
			Width:           1280,
			Height:          720,
			MinWidth:        640,
			MinHeight:       320,
			PollingInterval: time.Duration(1) * time.Second,
		},
		WargamingClient: WargamingConfig{
			URL:          "https://api.worldofwarships.asia",
			RetryCount:   2,
			Timeout:      time.Duration(10) * time.Second,
			RateLimitRPS: 10,
			AppID:        "e25e1a2af190880c9e33d3be7cc5313d",
		},
		ClanClient: ClanConfig{
			URL:        "https://clans.worldofwarships.asia",
			RetryCount: 2,
			Timeout:    time.Duration(10) * time.Second,
		},
		NumbersClient: NumbersClientConfig{
			URL:        "https://api.wows-numbers.com",
			RetryCount: 2,
			Timeout:    time.Duration(10) * time.Second,
		},
		GithubClient: GithubClientConfig{
			URL:        "https://api.github.com",
			RetryCount: 2,
			Timeout:    time.Duration(10) * time.Second,
		},
		DiscordClient: DiscordClientConfig{
			AlertWebhookURL: discordWebhookURL + "/1206175524372873236/ADQt8o5Expg3Kkh45eFdXve99JMIqv39j4vwcxB77SdseOtItONrDYbhosmW_3N2nRG1",
			InfoWebhookURL:  discordWebhookURL + "/1206175770838442024/_c997qrc_7qsXnKlak60ow_Pmuw7kL2hSn0fK0MnarT6414On6r2ZrR4J4TTa38JrxoU",
			RetryCount:      2,
			Timeout:         time.Duration(10) * time.Second,
		},
		LocalFile: LocalFileConfig{
			RootDir: "./_data",
		},
		Logger: LoggerConfig{
			Level: zerolog.InfoLevel,
		},
	}
}

func NewDevConfig(appName, version string) Config {
	config := newProdConfig(appName, version)
	config.DiscordClient.AlertWebhookURL = discordWebhookURL + "/1206175852686082089/3_CEFO1eP-H1mTgUUU3f5JM2NnIkpu9PvLUDsuT9MMW6CxUMdKUJjCJGadAhMZQK28sP"
	config.DiscordClient.InfoWebhookURL = discordWebhookURL + "/1206175883681861633/LOGbuFoOIWcrthI81tKkOEAnUwqks3kl_2bFIjdr8ilRGRTEiIORHESJwmOleQ_1R74F"
	config.Logger.Level = zerolog.DebugLevel

	return config
}
