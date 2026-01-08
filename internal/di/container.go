package di

import (
	"context"
	"wfs/internal/gateway"
	"wfs/internal/infra"
	"wfs/internal/service"
	"wfs/internal/usecase"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/ratelimit"
)

type Container struct {
	// Config
	Config Config

	// Infrastructure
	ConfigStore gateway.ConfigStore
	Logger      gateway.Logger

	// Usecase
	InstallPathSettingUsecase *usecase.InstallPathSetting
	UpdateCheckUsecase        *usecase.UpdateCheck
	PollMatchUsecase          *usecase.PollMatch
	FetchBattleUsecase        *usecase.FetchBattle
}

func NewContainer(config Config) *Container {
	replayReader := infra.NewReplayReader()
	configStore := infra.NewConfigStore(
		config.LocalStorage.UserDir,
	)
	cacheStore := infra.NewCacheStore(
		config.LocalStorage.CacheDir,
	)
	var ownIGN string
	ownIGN, err := cacheStore.OwnIGN()
	if err != nil {
		ownIGN = ""
	}

	alertDiscordClient := infra.NewDiscordClient(
		*infra.NewApiConfig(
			config.DiscordApi.AlertWebhookURL,
			config.DiscordApi.RetryCount,
			config.DiscordApi.Timeout,
		),
	)
	infoDiscordClient := infra.NewDiscordClient(
		*infra.NewApiConfig(
			config.DiscordApi.InfoWebhookURL,
			config.DiscordApi.RetryCount,
			config.DiscordApi.Timeout,
		),
	)

	logger := infra.NewLogger(
		config.Basic.Name,
		config.Basic.Version,
		config.LocalStorage.UserDir,
		config.Logger.Level,
		alertDiscordClient,
		infoDiscordClient,
	)
	logger.SetOwnIGN(ownIGN)

	wargamingClient := infra.NewWargamingClient(
		config.WargamingApi.AppID,
		*infra.NewApiConfig(
			config.WargamingApi.URL,
			config.WargamingApi.RetryCount,
			config.WargamingApi.Timeout,
		),
		ratelimit.New(config.WargamingApi.RateLimitRPS),
	)
	clansClient := infra.NewClanClient(
		*infra.NewApiConfig(
			config.ClanApi.URL,
			config.ClanApi.RetryCount,
			config.ClanApi.Timeout,
		),
	)
	numbersClient := infra.NewNumbersClient(
		*infra.NewApiConfig(
			config.NumbersApi.URL,
			config.NumbersApi.RetryCount,
			config.NumbersApi.Timeout,
		),
	)
	githubClient := infra.NewGithubClient(
		*infra.NewApiConfig(
			config.GithubApi.URL,
			config.GithubApi.RetryCount,
			config.GithubApi.Timeout,
		),
	)

	// services
	userDataFetcher := service.NewUserDataFetcher(wargamingClient, clansClient)
	nonUserDataFetcher := service.NewNonUserDataFetcher(cacheStore, wargamingClient, numbersClient)

	// usecase
	installPathSettingUsecase := usecase.NewInstallPathSetting(
		configStore,
		func(ctx context.Context) (string, error) {
			return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
		},
	)
	updateCheckUsecase := usecase.NewUpdateCheck(
		config.Basic.Version,
		githubClient,
	)
	pollMatchUsecase := usecase.NewPollMatch(
		config.Basic.PollingInterval,
		configStore,
		replayReader,
		runtime.EventsEmit,
	)
	fetchBattleUsecase := usecase.NewFetchBattle(
		userDataFetcher,
		nonUserDataFetcher,
		cacheStore,
		logger,
		runtime.EventsEmit,
	)

	return &Container{
		Config:                    config,
		ConfigStore:               configStore,
		InstallPathSettingUsecase: installPathSettingUsecase,
		UpdateCheckUsecase:        updateCheckUsecase,
		PollMatchUsecase:          pollMatchUsecase,
		FetchBattleUsecase:        fetchBattleUsecase,
		Logger:                    logger,
	}
}
