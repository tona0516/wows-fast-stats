package main

import (
	"context"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/repository"
	"wfs/backend/service"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const eventUpdateConfig = "CONFIG_UPDATE"
const eventUpdateAlertPlayers = "ALERT_PLAYERS_UPDATE"

type volatileData struct {
	cancelWatcher context.CancelFunc
}

func newVolatileData() volatileData {
	return volatileData{
		cancelWatcher: nil,
	}
}

//nolint:containedctx
type App struct {
	ctx            context.Context
	env            data.Env
	logger         repository.LoggerInterface
	config         service.Config
	screenshot     service.Screenshot
	watcher        service.Watcher
	battle         service.Battle
	updater        service.Updater
	configMigrator service.ConfigMigrator
	volatileData   volatileData
}

func NewApp(
	env data.Env,
	logger repository.LoggerInterface,
	config service.Config,
	screenshot service.Screenshot,
	watcher service.Watcher,
	battle service.Battle,
	updater service.Updater,
	configMigrator service.ConfigMigrator,
) *App {
	return &App{
		env:            env,
		logger:         logger,
		config:         config,
		screenshot:     screenshot,
		watcher:        watcher,
		battle:         battle,
		updater:        updater,
		configMigrator: configMigrator,
		volatileData:   newVolatileData(),
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

	if a.volatileData.cancelWatcher != nil {
		a.volatileData.cancelWatcher()
	}
	cancelCtx, cancel := context.WithCancel(context.Background())
	a.volatileData.cancelWatcher = cancel

	go a.watcher.Start(a.ctx, cancelCtx)

	return nil
}

func (a *App) Battle() (data.Battle, error) {
	result := data.Battle{}

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

func (a *App) DefaultUserConfig() data.UserConfigV2 {
	return data.DefaultUserConfigV2()
}

func (a *App) UserConfig() (data.UserConfigV2, error) {
	config, err := a.config.User()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) UpdateUserConfig(config data.UserConfigV2) error {
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

func (a *App) AlertPlayers() ([]data.AlertPlayer, error) {
	players, err := a.config.AlertPlayers()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player data.AlertPlayer) error {
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

func (a *App) SearchPlayer(prefix string) data.WGAccountList {
	return a.config.SearchPlayer(prefix)
}

func (a *App) AlertPatterns() []string {
	return data.AlertPatterns()
}

func (a *App) LogError(errString string, contexts map[string]string) {
	err := failure.New(apperr.FrontendError, failure.Messagef("%s", errString))
	a.logger.Error(err, contexts)
}

func (a *App) LogInfo(message string, contexts map[string]string) {
	a.logger.Info(message, contexts)
}

func (a *App) LatestRelease() (data.GHLatestRelease, error) {
	latestRelease, err := a.updater.IsUpdatable()
	return latestRelease, apperr.Unwrap(err)
}
