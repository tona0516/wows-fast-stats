package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"slices"
	"time"
	"wfs/backend2/infra/clans"
	"wfs/backend2/infra/numbers"
	"wfs/backend2/infra/wargaming"
	"wfs/backend2/usecase"

	"github.com/imroc/req/v3"
	"github.com/mitchellh/go-ps"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"go.uber.org/ratelimit"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var Base64ConfigYml string

func main() {
	if isAlreadyRunning() {
		os.Exit(0)

		return
	}

	config, err := getConfig()
	if err != nil {
		log.Fatalln(err)
	}

	injector := do.New()
	do.ProvideNamed(injector, "WargamingAPIClient", func(i *do.Injector) (*req.Client, error) {
		rateLimiter := ratelimit.New(10)

		return req.C().
			SetBaseURL(config.Wargaming.URL).
			AddCommonQueryParam("application_id", config.Wargaming.AppID).
			SetCommonRetryCount(config.Wargaming.MaxRetry).
			SetTimeout(time.Duration(config.Wargaming.TimeoutSec) * time.Second).
			OnBeforeRequest(func(client *req.Client, req *req.Request) error {
				rateLimiter.Take()

				return nil
			}).
			AddCommonRetryCondition(func(resp *req.Response, err error) bool {
				if err != nil {
					return true
				}

				var rb wargaming.ResponseCommon[any]
				if err := json.Unmarshal(resp.Bytes(), &rb); err != nil {
					return true
				}

				if rb.GetStatus() == "error" {
					// Note:
					// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
					message := rb.GetError().Message
					if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
						return true
					}
				}

				return false
			}).EnableTraceAll(), nil
	})
	do.ProvideNamedValue(injector, "WargamingAppID", config.Wargaming.AppID)
	do.ProvideNamed(injector, "ClansAPIClient", func(i *do.Injector) (*req.Client, error) {
		return req.C().
			SetBaseURL(config.UnofficialWargaming.URL).
			SetCommonRetryCount(config.UnofficialWargaming.MaxRetry).
			SetTimeout(time.Duration(config.UnofficialWargaming.TimeoutSec) * time.Second).EnableTraceAll(), nil
	})

	do.ProvideNamed(injector, "NumbersAPIClient", func(i *do.Injector) (*req.Client, error) {
		return req.C().
			SetBaseURL(config.Numbers.URL).
			SetCommonRetryCount(config.Numbers.MaxRetry).
			SetTimeout(time.Duration(config.Numbers.TimeoutSec) * time.Second).
			EnableInsecureSkipVerify().EnableTraceAll(), nil
	})

	// infra
	do.Provide(injector, wargaming.NewAPI)
	do.Provide(injector, clans.NewAPI)
	do.Provide(injector, numbers.NewAPI)

	// usecase
	do.Provide(injector, usecase.NewGetBattle)

	app := NewApp(config, injector)

	err = wails.Run(&options.App{
		Title:  config.App.Name,
		Width:  config.App.Width,
		Height: config.App.Height,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.onStartup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func isAlreadyRunning() bool {
	ownPid := os.Getpid()

	ownPidInfo, err := ps.FindProcess(ownPid)
	if err != nil {
		// Note: 可用性のためfalseを返す
		return false
	}

	processes, err := ps.Processes()
	if err != nil {
		// Note: 可用性のためfalseを返す
		return false
	}

	isRunning := false

	for _, p := range processes {
		if p.Pid() != ownPid && p.Executable() == ownPidInfo.Executable() {
			isRunning = true

			break
		}
	}

	return isRunning
}

func getConfig() (Config, error) {
	configYml, err := base64.StdEncoding.DecodeString(Base64ConfigYml)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := yaml.Unmarshal(configYml, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
