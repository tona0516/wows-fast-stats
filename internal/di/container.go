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

	// usecase
	InstallPathSettingUsecase *usecase.InstallPathSetting

	// services
	ConfigService   *service.Config
	BattlePublisher *service.BattlePublisher
	BattleService   *service.BattleFetcher
	UpdaterService  *service.UpdateChecker
	Logger          infra.Logger
}

func NewContainer(ctx context.Context, config Config) *Container {
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

	// usecase
	installPathSetting := usecase.NewInstallPathSetting(localStorage, func(ctx context.Context) (string, error) {
		return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	})

	// services
	configService := service.NewConfig(localStorage)
	battleFetcher := service.NewBattleFetcher(
		ctx,
		localStorage,
		wargamingApiClient,
		clansApiClient,
		numbersApiClient,
		logger,
		runtime.EventsEmit,
	)
	battlePublisher := service.NewBattlePublisher(
		ctx,
		config.Basic.PollingInterval,
		localStorage,
		logger,
		runtime.EventsEmit,
	)
	updaterService := service.NewUpdateChecker(
		config.Basic.Version,
		githubApiClient,
	)

	return &Container{
		Config:                    config,
		InstallPathSettingUsecase: installPathSetting,
		ConfigService:             configService,
		BattlePublisher:           battlePublisher,
		BattleService:             battleFetcher,
		UpdaterService:            updaterService,
		Logger:                    logger,
	}
}
