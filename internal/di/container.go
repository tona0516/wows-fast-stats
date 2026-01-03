package di

import (
	"context"
	"time"
	"wfs/internal/infra"
	"wfs/internal/repository"
	"wfs/internal/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Container struct {
	// Config
	Config Config

	// services
	ConfigService   *service.Config
	BattlePublisher *service.BattlePublisher
	BattleService   *service.BattleFetcher
	UpdaterService  *service.UpdateChecker
	Logger          repository.LoggerInterface
}

func NewContainer(ctx context.Context, config Config) *Container {
	alertDiscord := infra.NewDiscord(
		config.Discord.AlertURL,
		config.Discord.MaxRetry,
		config.Discord.TimeoutSec,
	)
	infoDiscord := infra.NewDiscord(
		config.Discord.InfoURL,
		config.Discord.MaxRetry,
		config.Discord.TimeoutSec,
	)

	persistence := infra.NewPersistence(config.Local.StoragePath)

	ownIGN, _ := persistence.LoadOwnIGN()
	logger := infra.NewLogger(
		config.App.Name,
		config.App.Semver,
		config.Logger.ZerologLogLevel,
		alertDiscord,
		infoDiscord,
	)
	logger.SetOwnIGN(ownIGN)
	logger.Init(ctx)

	wargaming := infra.NewWargaming(
		config.Wargaming.URL,
		config.Wargaming.MaxRetry,
		config.Wargaming.TimeoutSec,
		config.Wargaming.RetryIntervalMs,
		config.Wargaming.RateLimitRPS,
		config.Wargaming.AppID,
	)
	uwargaming := infra.NewUnofficialWargaming(
		config.UnofficialWargaming.URL,
		config.UnofficialWargaming.MaxRetry,
		config.UnofficialWargaming.TimeoutSec,
	)
	numbers := infra.NewNumbers(
		config.Numbers.URL,
		config.Numbers.MaxRetry,
		config.Numbers.TimeoutSec,
	)
	localFile := infra.NewLocalFile()
	github := infra.NewGithub(
		config.Github.URL,
		config.Github.MaxRetry,
		config.Github.TimeoutSec,
	)

	// services
	configService := service.NewSetting(persistence, wargaming, logger)
	battleFetcher := service.NewBattleFetcher(
		ctx,
		wargaming,
		uwargaming,
		numbers,
		persistence,
		logger,
		runtime.EventsEmit,
	)
	battlePublisher := service.NewBattlePublisher(
		ctx,
		time.Duration(config.Watcher.IntervalSec)*time.Second,
		localFile,
		persistence,
		logger,
		runtime.EventsEmit,
	)
	updaterService := service.NewUpdateChecker(config.App.Semver, github)

	return &Container{
		Config:          config,
		ConfigService:   configService,
		BattlePublisher: battlePublisher,
		BattleService:   battleFetcher,
		UpdaterService:  updaterService,
		Logger:          logger,
	}
}
