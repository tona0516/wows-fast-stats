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

	// Usecase
	InstallPathSettingUsecase *usecase.InstallPathSetting
	UpdateCheckUsecase        *usecase.UpdateCheck
	PollMatchUsecase          *usecase.PollMatch
	FetchBattleUsecase        *usecase.FetchBattle

	// Services
	ConfigService *service.Config
	Logger        infra.Logger
}

func NewContainer(config Config) *Container {
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

	localStorage := infra.NewLocalStorage(
		config.LocalStorage.UserDataDir,
		config.LocalStorage.CacheDir,
	)
	var ownIGN string
	ownIGN, err := localStorage.OwnIGN()
	if err != nil {
		ownIGN = ""
	}

	logger := infra.NewLogger(
		config.Basic.Name,
		config.Basic.Version,
		config.LocalStorage.UserDataDir,
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
	nonUserDataFetcher := service.NewNonUserDataFetcher(localStorage, wargamingApiClient, numbersApiClient)
	configService := service.NewConfig(localStorage)

	// usecase
	installPathSettingUsecase := usecase.NewInstallPathSetting(localStorage, func(ctx context.Context) (string, error) {
		return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	})
	updateCheckUsecase := usecase.NewUpdateCheck(
		config.Basic.Version,
		githubApiClient,
	)
	pollMatchUsecase := usecase.NewPollMatch(
		config.Basic.PollingInterval,
		localStorage,
		runtime.EventsEmit,
	)
	fetchBattleUsecase := usecase.NewFetchBattle(
		nonUserDataFetcher,
		localStorage,
		wargamingApiClient,
		clansApiClient,
		numbersApiClient,
		logger,
		runtime.EventsEmit,
	)

	return &Container{
		Config:                    config,
		InstallPathSettingUsecase: installPathSettingUsecase,
		UpdateCheckUsecase:        updateCheckUsecase,
		PollMatchUsecase:          pollMatchUsecase,
		FetchBattleUsecase:        fetchBattleUsecase,
		ConfigService:             configService,
		Logger:                    logger,
	}
}
