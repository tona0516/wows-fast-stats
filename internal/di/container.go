package di

import (
	"context"
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
	ConfigStore infra.ConfigStore

	// Usecase
	InstallPathSettingUsecase *usecase.InstallPathSetting
	UpdateCheckUsecase        *usecase.UpdateCheck
	PollMatchUsecase          *usecase.PollMatch
	FetchBattleUsecase        *usecase.FetchBattle

	// Services
	Logger infra.Logger
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

	alertDiscordApiClient := infra.NewDiscordApiClient(
		*infra.NewApiConfig(
			config.DiscordApi.AlertWebhookURL,
			config.DiscordApi.RetryCount,
			config.DiscordApi.Timeout,
		),
	)
	infoDiscordApiClient := infra.NewDiscordApiClient(
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
		alertDiscordApiClient,
		infoDiscordApiClient,
	)
	logger.SetOwnIGN(ownIGN)

	wargamingApiClient := infra.NewWargamingApiClient(
		config.WargamingApi.AppID,
		*infra.NewApiConfig(
			config.WargamingApi.URL,
			config.WargamingApi.RetryCount,
			config.WargamingApi.Timeout,
		),
		ratelimit.New(config.WargamingApi.RateLimitRPS),
	)
	clansApiClient := infra.NewClansApiClient(
		*infra.NewApiConfig(
			config.ClanApi.URL,
			config.ClanApi.RetryCount,
			config.ClanApi.Timeout,
		),
	)
	numbersApiClient := infra.NewNumbersApiClient(
		*infra.NewApiConfig(
			config.NumbersApi.URL,
			config.NumbersApi.RetryCount,
			config.NumbersApi.Timeout,
		),
	)
	githubApiClient := infra.NewGithubApiClient(
		*infra.NewApiConfig(
			config.GithubApi.URL,
			config.GithubApi.RetryCount,
			config.GithubApi.Timeout,
		),
	)

	// services
	userDataFetcher := service.NewUserDataFetcher(wargamingApiClient, clansApiClient)
	nonUserDataFetcher := service.NewNonUserDataFetcher(cacheStore, wargamingApiClient, numbersApiClient)

	// usecase
	installPathSettingUsecase := usecase.NewInstallPathSetting(
		configStore,
		func(ctx context.Context) (string, error) {
			return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
		},
	)
	updateCheckUsecase := usecase.NewUpdateCheck(
		config.Basic.Version,
		githubApiClient,
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
