package main

import (
	"context"
	"os"
	"wfs/backend/apperr"
	"wfs/backend/data"

	"github.com/mitchellh/go-ps"
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
	config        Config
	ctx           context.Context
	container     *DependencyContainer
	cancelWatcher context.CancelFunc
}

func NewApp(config Config) *App {
	return &App{config: config}
}

func (a *App) MigrateIfNeeded() error {
	if err := a.container.configMigratorService.ExecuteIfNeeded(); err != nil {
		a.container.logger.Error(err, nil)
		return apperr.Unwrap(err)
	}

	return nil
}

func (a *App) StartWatching() error {
	if err := a.container.watcherService.Prepare(); err != nil {
		a.container.logger.Error(err, nil)
		return apperr.Unwrap(err)
	}

	if a.cancelWatcher != nil {
		a.cancelWatcher()
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	a.cancelWatcher = cancel

	go a.container.watcherService.Start(a.ctx, cancelCtx)

	return nil
}

func (a *App) Battle() (data.Battle, error) {
	result := data.Battle{}

	userConfig, err := a.container.configService.User()
	if err != nil {
		a.container.logger.Error(err, nil)
		return result, apperr.Unwrap(err)
	}

	result, err = a.container.battleService.Get(a.ctx, userConfig)
	if err != nil {
		a.container.logger.Error(err, nil)
		return result, apperr.Unwrap(err)
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := a.container.configService.SelectDirectory(a.ctx)
	if err != nil {
		a.container.logger.Error(err, nil)
	}

	return path, apperr.Unwrap(err)
}

func (a *App) OpenDirectory(path string) error {
	err := a.container.configService.OpenDirectory(path)
	if err != nil {
		a.container.logger.Warn(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) DefaultUserConfig() data.UserConfigV2 {
	return data.DefaultUserConfigV2()
}

func (a *App) UserConfig() (data.UserConfigV2, error) {
	config, err := a.container.configService.User()
	if err != nil {
		a.container.logger.Error(err, nil)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) UpdateUserConfig(config data.UserConfigV2) error {
	err := a.container.configService.UpdateOptional(config)
	if err != nil {
		a.container.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateConfig, config)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateInstallPath(path string) string {
	err := a.container.configService.ValidateInstallPath(path)

	if err := apperr.Unwrap(err); err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) UpdateInstallPath(path string) error {
	config, err := a.container.configService.UpdateInstallPath(path)
	if err != nil {
		a.container.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateConfig, config)
	}

	return apperr.Unwrap(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) (bool, error) {
	saved, err := a.container.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.container.logger.Error(err, nil)
	}
	return saved, apperr.Unwrap(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.container.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.container.logger.Error(err, nil)
	}
	return apperr.Unwrap(err)
}

func (a *App) Semver() string {
	return *a.config.App.Semver
}

func (a *App) AlertPlayers() ([]data.AlertPlayer, error) {
	players, err := a.container.configService.AlertPlayers()
	if err != nil {
		a.container.logger.Error(err, nil)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player data.AlertPlayer) error {
	players, err := a.container.configService.UpdateAlertPlayer(player)
	if err != nil {
		a.container.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	players, err := a.container.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		a.container.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) data.WGAccountList {
	return a.container.configService.SearchPlayer(prefix)
}

func (a *App) AlertPatterns() []string {
	return data.AlertPatterns()
}

func (a *App) LogError(errString string, contexts map[string]string) {
	err := failure.New(apperr.FrontendError, failure.Messagef("%s", errString))
	a.container.logger.Error(err, contexts)
}

func (a *App) LogInfo(message string, contexts map[string]string) {
	a.container.logger.Info(message, contexts)
}

func (a *App) LatestRelease() (data.GHLatestRelease, error) {
	latestRelease, err := a.container.updaterService.IsUpdatable()
	return latestRelease, apperr.Unwrap(err)
}

func (a *App) onStartup(ctx context.Context) {
	runtime.LogSetLogLevel(ctx, logger.INFO)

	if isAlreadyRunning() {
		a.showExistDialog(ctx, "すでに起動しています。", 1)
	}

	container, err := NewDependencyContainer(ctx, a.config)
	if err != nil {
		a.showExistDialog(ctx, "意図しないエラーが発生しました。\n"+err.Error(), 1)
	}

	a.ctx = ctx
	a.container = container
}

func isAlreadyRunning() bool {
	ownPid := os.Getpid()
	ownPidInfo, err := ps.FindProcess(ownPid)
	if err != nil {
		// Note: 可用性のためfalseを返す
		return false
	}

	processes, err := ps.Processes()
	if err != nil {
		// Note: 可用性のためfalseを返す
		return false
	}

	isRunning := false
	for _, p := range processes {
		if p.Pid() != ownPid && p.Executable() == ownPidInfo.Executable() {
			isRunning = true
			break
		}
	}

	return isRunning
}

func (a *App) showExistDialog(ctx context.Context, message string, code int) {
	_, _ = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Title:   *a.config.App.Name,
		Message: message,
	})
	os.Exit(code)
}
