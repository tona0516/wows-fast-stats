package main

import (
	"context"
	"os"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/service"

	"github.com/mitchellh/go-ps"
	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//nolint:containedctx
type App struct {
	config            Config
	ctx               context.Context
	container         *DependencyContainer
	unsubscribeBattle context.CancelFunc
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

func (a *App) SubscribeBattle() {
	if !a.container.battlePublisher.CanSubcribe() {
		return
	}

	if a.unsubscribeBattle != nil {
		a.unsubscribeBattle()
	}

	cancelCtx, unsubscribeBattle := context.WithCancel(context.Background())
	a.unsubscribeBattle = unsubscribeBattle
	channel := make(chan data.TempArenaInfo)

	go a.container.battlePublisher.Subcribe(cancelCtx, channel)
	for tempArenaInfo := range channel {
		a.container.battleService.Invoke(tempArenaInfo)
	}
}

func (a *App) TrySaveInstallPath() (bool, error) {
	path, err := a.container.configService.SelectDirectory(a.ctx)
	if err != nil {
		return false, apperr.Unwrap(err)
	}

	if path == "" {
		return false, nil
	}

	config, err := a.container.configService.UpdateInstallPath(path)
	if err != nil {
		return false, apperr.Unwrap(err)
	}

	runtime.EventsEmit(a.ctx, service.EventUpdateConfig, config)

	return true, nil
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
		runtime.EventsEmit(a.ctx, service.EventUpdateConfig, config)
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

func (a *App) Semver() string {
	return a.config.App.Semver
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
		runtime.EventsEmit(a.ctx, service.EventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	players, err := a.container.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		a.container.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, service.EventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) ([]data.WGAccountListData, error) {
	result, err := a.container.configService.SearchPlayer(prefix)

	return result, apperr.Unwrap(err)
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

func (a *App) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   a.config.App.Name,
		Message: message,
	})
}

// 構造体のバインド用のメソッド
func (a *App) EmptyBattle() data.Battle {
	return data.Battle{}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogSetLogLevel(ctx, logger.INFO)

	if isAlreadyRunning() {
		a.showExistDialog("すでに起動しています。", 1)
	}

	container, err := NewDependencyContainer(ctx, a.config)
	if err != nil {
		a.showExistDialog("意図しないエラーが発生しました。\n"+err.Error(), 1)
	}
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

func (a *App) showExistDialog(message string, code int) {
	a.ShowMessageDialog(message)
	os.Exit(code)
}
