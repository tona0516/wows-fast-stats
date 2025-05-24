package main

import (
	"context"
	"time"
	"wfs/backend/infra"
	"wfs/backend/repository"
	"wfs/backend/service"

	"github.com/dgraph-io/badger/v4"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DependencyContainer struct {
	// config
	config Config

	// services
	configService         *service.Config
	screenshotService     *service.Screenshot
	watcherService        *service.Watcher
	battleService         *service.Battle
	updaterService        *service.Updater
	configMigratorService *service.ConfigMigrator
	logger                repository.LoggerInterface
}

func NewDependencyContainer(ctx context.Context, config Config) (*DependencyContainer, error) {
	alertDiscord := infra.NewDiscord(
		*config.Discord.AlertURL,
		*config.Discord.MaxRetry,
		*config.Discord.TimeoutSec,
	)
	infoDiscord := infra.NewDiscord(
		*config.Discord.InfoURL,
		*config.Discord.MaxRetry,
		*config.Discord.TimeoutSec,
	)

	options := badger.DefaultOptions(*config.Local.StoragePath)
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	storage := infra.NewStorage(db)
	ownIGN, _ := storage.OwnIGN()

	logger := infra.NewLogger(
		*config.App.Name,
		*config.App.Semver,
		*config.Logger.ZerologLogLevel,
		alertDiscord,
		infoDiscord,
	)
	logger.SetOwnIGN(ownIGN)
	logger.Init(ctx)

	wargaming := infra.NewWargaming(
		*config.Wargaming.URL,
		*config.Wargaming.MaxRetry,
		*config.Wargaming.TimeoutSec,
		*config.Wargaming.RetryIntervalMs,
		*config.Wargaming.RateLimitRPS,
		*config.Wargaming.AppID,
	)
	uwargaming := infra.NewUnofficialWargaming(
		*config.UnofficialWargaming.URL,
		*config.UnofficialWargaming.MaxRetry,
		*config.UnofficialWargaming.TimeoutSec,
	)
	numbers := infra.NewNumbers(
		*config.Numbers.URL,
		*config.Numbers.MaxRetry,
		*config.Numbers.TimeoutSec,
	)
	localFile := infra.NewLocalFile()
	configV0 := infra.NewConfigV0()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(
		*config.Github.URL,
		*config.Github.MaxRetry,
		*config.Github.TimeoutSec,
	)

	// services
	configService := service.NewConfig(localFile, wargaming, storage, logger)
	screenshotService := service.NewScreenshot(localFile, logger)
	battleService := service.NewBattle(
		wargaming,
		uwargaming,
		localFile,
		numbers,
		unregistered,
		storage,
		logger,
		runtime.EventsEmit,
	)
	watcherService := service.NewWatcher(
		time.Duration(*config.Watcher.IntervalSec)*time.Second,
		localFile,
		storage,
		logger,
		runtime.EventsEmit,
	)
	updaterService := service.NewUpdater(*config.App.Semver, github, logger)
	configMigratorService := service.NewConfigMigrator(configV0, storage, logger)

	return &DependencyContainer{
		config:                config,
		configService:         configService,
		screenshotService:     screenshotService,
		watcherService:        watcherService,
		battleService:         battleService,
		updaterService:        updaterService,
		configMigratorService: configMigratorService,
		logger:                logger,
	}, nil
}
