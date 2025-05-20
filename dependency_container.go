package main

import (
	"context"
	"time"
	"wfs/backend/infra"
	"wfs/backend/repository"
	"wfs/backend/service"

	"github.com/dgraph-io/badger/v4"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/ratelimit"
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
	alertDiscord := infra.NewDiscord(infra.RequestConfig{
		URL:     config.Discord.AlertURL,
		Retry:   config.Discord.MaxRetry,
		Timeout: time.Duration(config.Discord.TimeoutSec) * time.Second,
	})
	infoDiscord := infra.NewDiscord(infra.RequestConfig{
		URL:     config.Discord.InfoURL,
		Retry:   config.Discord.MaxRetry,
		Timeout: time.Duration(config.Discord.TimeoutSec) * time.Second,
	})

	options := badger.DefaultOptions(config.Local.StoragePath)
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	storage := infra.NewStorage(db)
	ownIGN, _ := storage.OwnIGN()

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
		infra.RequestConfig{
			URL:     config.Wargaming.URL,
			Retry:   config.Wargaming.MaxRetry,
			Timeout: time.Duration(config.Wargaming.TimeoutSec) * time.Second,
		},
		ratelimit.New(config.Wargaming.RateLimitRPS),
		config.Wargaming.AppID,
	)
	uwargaming := infra.NewUnofficialWargaming(infra.RequestConfig{
		URL:     config.UnofficialWargaming.URL,
		Retry:   config.UnofficialWargaming.MaxRetry,
		Timeout: time.Duration(config.UnofficialWargaming.TimeoutSec) * time.Second,
	})
	numbers := infra.NewNumbers(infra.RequestConfig{
		URL:     config.Numbers.URL,
		Retry:   config.Numbers.MaxRetry,
		Timeout: time.Duration(config.Numbers.TimeoutSec) * time.Second,
	})
	localFile := infra.NewLocalFile()
	configV0 := infra.NewConfigV0()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(infra.RequestConfig{
		URL:     config.Github.URL,
		Retry:   config.Github.MaxRetry,
		Timeout: time.Duration(config.Github.TimeoutSec) * time.Second,
	})

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
		time.Duration(config.Watcher.IntervalSec)*time.Second,
		localFile,
		storage,
		logger,
		runtime.EventsEmit,
	)
	updaterService := service.NewUpdater(config.App.Semver, github, logger)
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
