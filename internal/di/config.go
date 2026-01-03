package di

type Config struct {
	App struct {
		Name      string
		Semver    string
		Width     int
		Height    int
		MinWidth  int
		MinHeight int
	}
	Watcher struct {
		IntervalSec int
	}
	Wargaming struct {
		URL             string
		MaxRetry        int
		TimeoutSec      int
		RetryIntervalMs int
		RateLimitRPS    int
		AppID           string
	}
	UnofficialWargaming struct {
		URL        string
		MaxRetry   int
		TimeoutSec int
	}
	Numbers struct {
		URL        string
		MaxRetry   int
		TimeoutSec int
	}
	Github struct {
		URL        string
		MaxRetry   int
		TimeoutSec int
	}
	Discord struct {
		AlertURL   string
		InfoURL    string
		MaxRetry   int
		TimeoutSec int
	}
	Local struct {
		StoragePath string
	}
	Logger struct {
		ZerologLogLevel string
	}
}

func NewConfig(env string) Config {
	if env == "prod" {
		return NewProdConfig()
	} else {
		return NewDevConfig()
	}
}

func NewProdConfig() Config {
	return Config{
		App: struct {
			Name      string
			Semver    string
			Width     int
			Height    int
			MinWidth  int
			MinHeight int
		}{
			Name:      "wows-fast-stats",
			Semver:    "1.0.0-alpha.1",
			Width:     1280,
			Height:    720,
			MinWidth:  640,
			MinHeight: 320,
		},
		Watcher: struct{ IntervalSec int }{
			IntervalSec: 1,
		},
		Wargaming: struct {
			URL             string
			MaxRetry        int
			TimeoutSec      int
			RetryIntervalMs int
			RateLimitRPS    int
			AppID           string
		}{
			URL:             "https://api.worldofwarships.asia",
			MaxRetry:        2,
			TimeoutSec:      10,
			RetryIntervalMs: 500,
			RateLimitRPS:    10,
			AppID:           "e25e1a2af190880c9e33d3be7cc5313d",
		},
		UnofficialWargaming: struct {
			URL        string
			MaxRetry   int
			TimeoutSec int
		}{
			URL:        "https://clans.worldofwarships.asia",
			MaxRetry:   2,
			TimeoutSec: 10,
		},
		Numbers: struct {
			URL        string
			MaxRetry   int
			TimeoutSec int
		}{
			URL:        "https://api.wows-numbers.com",
			MaxRetry:   2,
			TimeoutSec: 10,
		},
		Github: struct {
			URL        string
			MaxRetry   int
			TimeoutSec int
		}{
			URL:        "https://api.github.com",
			MaxRetry:   2,
			TimeoutSec: 10,
		},
		Discord: struct {
			AlertURL   string
			InfoURL    string
			MaxRetry   int
			TimeoutSec int
		}{
			AlertURL:   "https://discord.com/api/webhooks/1206175524372873236/ADQt8o5Expg3Kkh45eFdXve99JMIqv39j4vwcxB77SdseOtItONrDYbhosmW_3N2nRG1",
			InfoURL:    "https://discord.com/api/webhooks/1206175770838442024/_c997qrc_7qsXnKlak60ow_Pmuw7kL2hSn0fK0MnarT6414On6r2ZrR4J4TTa38JrxoU",
			MaxRetry:   2,
			TimeoutSec: 10,
		},
		Local: struct{ StoragePath string }{
			StoragePath: "./user_data",
		},
		Logger: struct{ ZerologLogLevel string }{
			ZerologLogLevel: "info",
		},
	}
}

func NewDevConfig() Config {
	config := NewProdConfig()
	config.Discord.AlertURL = "https://discord.com/api/webhooks/1206175852686082089/3_CEFO1eP-H1mTgUUU3f5JM2NnIkpu9PvLUDsuT9MMW6CxUMdKUJjCJGadAhMZQK28sP"
	config.Discord.InfoURL = "https://discord.com/api/webhooks/1206175883681861633/LOGbuFoOIWcrthI81tKkOEAnUwqks3kl_2bFIjdr8ilRGRTEiIORHESJwmOleQ_1R74F"
	config.Logger.ZerologLogLevel = "debug"

	return config
}
