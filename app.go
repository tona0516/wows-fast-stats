package main

import (
	"context"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"
	"wfs/backend/service"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	eventUpdateConfig       = "CONFIG_UPDATE"
	eventUpdateAlertPlayers = "ALERT_PLAYERS_UPDATE"
)

//nolint:containedctx
type App struct {
	ctx            context.Context
	env            data.Env
	logger         repository.Logger
	versionFetcher repository.VersionFetcher
	config         service.Config
	screenshot     service.Screenshot
	watcher        service.Watcher
	battle         service.Battle
	configMigrator service.ConfigMigrator

	cancelWacthFunc context.CancelFunc
}

func NewApp(
	env data.Env,
	logger repository.Logger,
	versionFetcher repository.VersionFetcher,
	config service.Config,
	screenshot service.Screenshot,
	watcher service.Watcher,
	battle service.Battle,
	configMigrator service.ConfigMigrator,
) *App {
	return &App{
		env:            env,
		logger:         logger,
		versionFetcher: versionFetcher,
		config:         config,
		screenshot:     screenshot,
		watcher:        watcher,
		battle:         battle,
		configMigrator: configMigrator,
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogSetLogLevel(ctx, logger.INFO)
	a.logger.Init(ctx)
}

func (a *App) MigrateIfNeeded() error {
	if err := a.configMigrator.ExecuteIfNeeded(); err != nil {
		a.logger.Error(err, nil)
		return apperr.Unwrap(err)
	}

	return nil
}

func (a *App) StartWatching() error {
	if err := a.watcher.Prepare(); err != nil {
		a.logger.Error(err, nil)
		return apperr.Unwrap(err)
	}

	if a.cancelWacthFunc != nil {
		a.cancelWacthFunc()
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	a.cancelWacthFunc = cancel

	go a.watcher.Start(a.ctx, cancelCtx)

	return nil
}

func (a *App) Battle() (model.Battle, error) {
	result := model.Battle{}

	userConfig, err := a.config.User()
	if err != nil {
		a.logger.Error(err, nil)
		return result, apperr.Unwrap(err)
	}

	result, err = a.battle.Get(a.ctx, userConfig)
	if err != nil {
		a.logger.Error(err, nil)
		return result, apperr.Unwrap(err)
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := a.config.SelectDirectory(a.ctx)
	if err != nil {
		a.logger.Error(err, nil)
	}

	return path, apperr.Unwrap(err)
}

func (a *App) OpenDirectory(path string) error {
	err := a.config.OpenDirectory(path)
	if err != nil {
		a.logger.Warn(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) DefaultUserConfig() model.UserConfigV2 {
	return model.DefaultUserConfigV2()
}

func (a *App) UserConfig() (model.UserConfigV2, error) {
	config, err := a.config.User()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) UpdateUserConfig(config model.UserConfigV2) error {
	err := a.config.UpdateOptional(config)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateConfig, config)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateInstallPath(path string) string {
	err := a.config.ValidateInstallPath(path)
	if err != nil {
		a.logger.Error(err, nil)
	}

	if err := apperr.Unwrap(err); err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) UpdateInstallPath(path string) error {
	config, err := a.config.UpdateInstallPath(path)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateConfig, config)
	}

	return apperr.Unwrap(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) (bool, error) {
	saved, err := a.screenshot.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Error(err, nil)
	}
	return saved, apperr.Unwrap(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshot.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Error(err, nil)
	}
	return apperr.Unwrap(err)
}

func (a *App) Semver() string {
	return a.env.Semver
}

func (a *App) AlertPlayers() ([]model.AlertPlayer, error) {
	players, err := a.config.AlertPlayers()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player model.AlertPlayer) error {
	players, err := a.config.UpdateAlertPlayer(player)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	players, err := a.config.RemoveAlertPlayer(accountID)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) map[string]int {
	return a.config.SearchPlayer(prefix)
}

func (a *App) AlertPatterns() []string {
	return model.AlertPatterns()
}

func (a *App) LogError(errString string, contexts map[string]string) {
	err := failure.New(apperr.FrontendError, failure.Message(errString))
	a.logger.Error(err, contexts)
}

func (a *App) LogInfo(message string, contexts map[string]string) {
	a.logger.Info(message, contexts)
}

func (a *App) LatestRelease() (model.LatestRelease, error) {
	latestRelease, err := a.versionFetcher.Fetch(a.env.Semver)
	return latestRelease, apperr.Unwrap(err)
}
