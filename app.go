package main

import (
	"context"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"
	"wfs/backend/usecase"

	"github.com/morikuni/failure"
)

type volatileData struct {
	cancelWatcher  context.CancelFunc
	excludePlayers model.ExcludedPlayers
}

func newVolatileData() volatileData {
	return volatileData{
		cancelWatcher:  nil,
		excludePlayers: make(model.ExcludedPlayers),
	}
}

//nolint:containedctx
type App struct {
	ctx            context.Context
	env            model.Env
	logger         repository.LoggerInterface
	config         usecase.Config
	screenshot     usecase.Screenshot
	watcher        usecase.Watcher
	battle         usecase.Battle
	updater        usecase.Updater
	configMigrator usecase.ConfigMigrator
	volatileData   volatileData
}

func NewApp(
	env model.Env,
	logger repository.LoggerInterface,
	config usecase.Config,
	screenshot usecase.Screenshot,
	watcher usecase.Watcher,
	battle usecase.Battle,
	updater usecase.Updater,
	configMigrator usecase.ConfigMigrator,
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
	return model.DefaultUserConfigV2
}

func (a *App) UserConfig() (model.UserConfigV2, error) {
	config, err := a.config.User()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) ApplyUserConfig(config model.UserConfigV2) error {
	err := a.config.UpdateOptional(config)
	if err != nil {
		a.logger.Error(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateRequiredConfig(
	installPath string,
	appid string,
) model.RequiredConfigError {
	return a.config.ValidateRequired(installPath, appid)
}

func (a *App) ApplyRequiredUserConfig(
	installPath string,
	appid string,
) (model.RequiredConfigError, error) {
	validatedResult, err := a.config.UpdateRequired(installPath, appid)
	if err != nil {
		a.logger.Error(err, nil)
	}

	return validatedResult, apperr.Unwrap(err)
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

func (a *App) ExcludePlayerIDs() []int {
	return a.volatileData.excludePlayers.IDs()
}

func (a *App) AddExcludePlayerID(playerID int) {
	a.volatileData.excludePlayers.Add(playerID)
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	a.volatileData.excludePlayers.Remove(playerID)
}

func (a *App) AlertPlayers() ([]model.AlertPlayer, error) {
	players, err := a.config.AlertPlayers()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player model.AlertPlayer) error {
	err := a.config.UpdateAlertPlayer(player)
	if err != nil {
		a.logger.Error(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	err := a.config.RemoveAlertPlayer(accountID)
	if err != nil {
		a.logger.Error(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) model.WGAccountList {
	return a.config.SearchPlayer(prefix)
}

func (a *App) AlertPatterns() []string {
	return model.AlertPatterns
}

func (a *App) LogError(errString string, contexts map[string]string) {
	err := failure.New(apperr.FrontendError, failure.Messagef("%s", errString))
	a.logger.Error(err, contexts)
}

func (a *App) LogInfo(message string, contexts map[string]string) {
	a.logger.Info(message, contexts)
}

func (a *App) LatestRelease() (model.GHLatestRelease, error) {
	latestRelease, err := a.updater.IsUpdatable()
	return latestRelease, apperr.Unwrap(err)
}
