package main

import (
	"embed"
	"log"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/controller"
	"wfs/backend/infra"
	"wfs/backend/service"
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

func main() {
	injector := do.New()

	appConfig := config.NewConfig(appName, version, env)
	do.ProvideValue(injector, appConfig)

	// infra
	do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
		return infra.NewWails(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.CacheStore, error) {
		return infra.NewCacheStore(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.ConfigStore, error) {
		return infra.NewConfigStore(i)
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
			appConfig.DiscordClient.AlertWebhookURL,
			appConfig.DiscordClient.RetryCount,
			appConfig.DiscordClient.Timeout,
		)
	})
	do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return infra.NewDiscordClient(
			appConfig.DiscordClient.InfoWebhookURL,
			appConfig.DiscordClient.RetryCount,
			appConfig.DiscordClient.Timeout,
		)
	})

	// service
	do.Provide(injector, service.NewUserDataFetcher)
	do.Provide(injector, service.NewNonUserDataFetcher)

	// usecase
	do.Provide(injector, usecase.NewFetchBattle)
	do.Provide(injector, usecase.NewPollMatch)
	do.Provide(injector, usecase.NewInstallPathSetting)
	do.Provide(injector, usecase.NewUpdateCheck)

	// controller
	do.Provide(injector, controller.NewController)

	controller := do.MustInvoke[*controller.Controller](injector)

	err := wails.Run(&options.App{
		Title:     appConfig.Basic.Name,
		Width:     appConfig.Basic.Width,
		Height:    appConfig.Basic.Height,
		MinWidth:  appConfig.Basic.MinWidth,
		MinHeight: appConfig.Basic.MinHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: controller.OnStartup,
		Bind:      []any{controller},
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}
