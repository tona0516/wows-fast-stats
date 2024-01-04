package main

import (
	"context"
	"wfs/backend/adapter"
	"wfs/backend/apperr"
	"wfs/backend/application/usecase"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const EventOnload = "ONLOAD"

//nolint:containedctx
type App struct {
	ctx            context.Context
	env            vo.Env
	cancelWatcher  context.CancelFunc
	logger         adapter.LoggerInterface
	config         usecase.Config
	screenshot     usecase.Screenshot
	watcher        usecase.Watcher
	battle         usecase.Battle
	updater        usecase.Updater
	configMigrator usecase.ConfigMigrator
	excludePlayers domain.ExcludedPlayers
}

func NewApp(
	env vo.Env,
	logger adapter.LoggerInterface,
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
		excludePlayers: domain.ExcludedPlayers{},
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Init(ctx)

	runtime.EventsOn(ctx, EventOnload, func(optionalData ...interface{}) {
		a.logger.Info("application started")
	})
}

func (a *App) Migrate() error {
	if err := a.configMigrator.Execute(); err != nil {
		a.logger.Error(err)
		return apperr.Unwrap(err)
	}

	return nil
}

func (a *App) StartWatching() error {
	if err := a.watcher.Prepare(); err != nil {
		a.logger.Error(err)
		return apperr.Unwrap(err)
	}

	if a.cancelWatcher != nil {
		a.cancelWatcher()
	}
	cancelCtx, cancel := context.WithCancel(context.Background())
	a.cancelWatcher = cancel

	go a.watcher.Start(a.ctx, cancelCtx)

	return nil
}

func (a *App) Battle() (domain.Battle, error) {
	result := domain.Battle{}

	userConfig, err := a.config.User()
	if err != nil {
		a.logger.Error(err)
		return result, apperr.Unwrap(err)
	}

	result, err = a.battle.Get(userConfig)
	if err != nil {
		a.logger.Error(err)
		return result, apperr.Unwrap(err)
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := a.config.SelectDirectory(a.ctx)
	if err != nil {
		a.logger.Error(err)
	}

	return path, apperr.Unwrap(err)
}

func (a *App) OpenDirectory(path string) error {
	err := a.config.OpenDirectory(path)
	if err != nil {
		a.logger.Warn(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) DefaultUserConfig() domain.UserConfig {
	return domain.DefaultUserConfig
}

func (a *App) UserConfig() (domain.UserConfig, error) {
	config, err := a.config.User()
	if err != nil {
		a.logger.Error(err)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) ApplyUserConfig(config domain.UserConfig) error {
	err := a.config.UpdateOptional(config)
	if err != nil {
		a.logger.Error(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateRequiredConfig(
	installPath string,
	appid string,
) vo.RequiredConfigError {
	return a.config.ValidateRequired(installPath, appid)
}

func (a *App) ApplyRequiredUserConfig(
	installPath string,
	appid string,
) (vo.RequiredConfigError, error) {
	validatedResult, err := a.config.UpdateRequired(installPath, appid)
	if err != nil {
		a.logger.Error(err)
	}

	return validatedResult, apperr.Unwrap(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) (bool, error) {
	saved, err := a.screenshot.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Error(err)
	}
	return saved, apperr.Unwrap(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshot.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Error(err)
	}
	return apperr.Unwrap(err)
}

func (a *App) Semver() string {
	return a.env.Semver
}

func (a *App) ExcludePlayerIDs() []int {
	return a.excludePlayers.IDs()
}

func (a *App) AddExcludePlayerID(playerID int) {
	a.excludePlayers.Add(playerID)
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	a.excludePlayers.Remove(playerID)
}

func (a *App) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := a.config.AlertPlayers()
	if err != nil {
		a.logger.Error(err)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player domain.AlertPlayer) error {
	err := a.config.UpdateAlertPlayer(player)
	if err != nil {
		a.logger.Error(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	err := a.config.RemoveAlertPlayer(accountID)
	if err != nil {
		a.logger.Error(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) (domain.WGAccountList, error) {
	accountList, err := a.config.SearchPlayer(prefix)
	if err != nil {
		a.logger.Error(err)
	}

	return accountList, apperr.Unwrap(err)
}

func (a *App) AlertPatterns() []string {
	return domain.AlertPatterns
}

func (a *App) LogError(errString string) {
	err := failure.New(apperr.FrontendError, failure.Messagef("%s", errString))
	a.logger.Error(err)
}

func (a *App) LatestRelease() (domain.GHLatestRelease, error) {
	latestRelease, err := a.updater.IsUpdatable()
	return latestRelease, apperr.Unwrap(err)
}
