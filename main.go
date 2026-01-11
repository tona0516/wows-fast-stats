package main

import (
	"embed"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/infra"
	"wfs/backend/usecase"

	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var (
	appName string
	version string
	env     string
)

func injectDependency() *options.App {
	injector := do.New()

	config := config.NewConfig(appName, version, env)
	do.ProvideValue(injector, config)

	// infra
	do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
		return infra.NewWails(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.CacheStore, error) {
		return infra.NewCacheStore(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
		return infra.NewPrefStore(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.ReplayReader, error) {
		return infra.NewReplayReader(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.ClanClient, error) {
		return infra.NewClanClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.GithubClient, error) {
		return infra.NewGithubClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.NumbersClient, error) {
		return infra.NewNumbersClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
		return infra.NewWargamingClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.Logger, error) {
		return infra.NewLogger(i)
	})
	do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return infra.NewDiscordClient(
			config.DiscordClient.AlertWebhookURL,
			config.DiscordClient.RetryCount,
			config.DiscordClient.Timeout,
		)
	})
	do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return infra.NewDiscordClient(
			config.DiscordClient.InfoWebhookURL,
			config.DiscordClient.RetryCount,
			config.DiscordClient.Timeout,
		)
	})

	// usecase
	do.Provide(injector, usecase.NewPrefetch)
	do.Provide(injector, usecase.NewFetchBattle)
	do.Provide(injector, usecase.NewPollMatch)
	do.Provide(injector, usecase.NewInstallPathSetting)
	do.Provide(injector, usecase.NewUpdateCheck)
	do.Provide(injector, usecase.NewLoadPref)
	do.Provide(injector, usecase.NewSavePref)

	// controller
	do.Provide(injector, NewApp)

	app := do.MustInvoke[*App](injector)

	return &options.App{
		Title:     config.Basic.Name,
		Width:     config.Basic.Width,
		Height:    config.Basic.Height,
		MinWidth:  config.Basic.MinWidth,
		MinHeight: config.Basic.MinHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.OnStartup,
		Bind:      []any{app},
	}
}

func main() {
	app := injectDependency()
	if err := wails.Run(app); err != nil {
		panic(err)
	}
}
