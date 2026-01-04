package main

import (
	"context"
	"os"
	"wfs/internal/apperr"
	"wfs/internal/data"
	"wfs/internal/di"

	"github.com/mitchellh/go-ps"
	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	config            di.Config
	ctx               context.Context
	container         *di.Container
	unsubscribeBattle context.CancelFunc
}

func NewApp(config di.Config) *App {
	return &App{config: config}
}

func (a *App) SubscribeBattle() {
	if !a.container.BattlePublisher.CanSubcribe() {
		return
	}

	if a.unsubscribeBattle != nil {
		a.unsubscribeBattle()
	}

	cancelCtx, unsubscribeBattle := context.WithCancel(context.Background())
	a.unsubscribeBattle = unsubscribeBattle
	channel := make(chan data.TempArenaInfo)

	go a.container.BattlePublisher.Subcribe(cancelCtx, channel)
	for tempArenaInfo := range channel {
		a.container.BattleService.Invoke(tempArenaInfo)
	}
}

func (a *App) GetUserConfig() (data.UserConfig, error) {
	return a.container.ConfigService.GetUserConfig()
}

func (a *App) SaveUserConfig(config data.UserConfig) error {
	return a.container.ConfigService.SaveUserConfig(config)
}

func (a *App) TrySaveInstallPath() (bool, error) {
	return a.container.ConfigService.TrySaveInstallPath(a.ctx)
}

func (a *App) OpenDirectory(path string) error {
	err := a.container.ConfigService.OpenDirectory(path)
	if err != nil {
		a.container.Logger.Error(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateInstallPath(path string) string {
	err := a.container.ConfigService.ValidateInstallPath(path)

	if err := apperr.Unwrap(err); err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) Semver() string {
	return a.config.Basic.Version
}

func (a *App) SearchPlayer(prefix string) ([]data.WGAccountListData, error) {
	result, err := a.container.ConfigService.SearchPlayer(prefix)

	return result.Data, apperr.Unwrap(err)
}

func (a *App) LogError(errString string, contexts map[string]string) {
	err := failure.New(apperr.FrontendError, failure.Messagef("%s", errString))
	a.container.Logger.Error(err, contexts)
}

func (a *App) LogInfo(message string, contexts map[string]string) {
	a.container.Logger.Info(message, contexts)
}

func (a *App) NewVersion() *data.NewVersion {
	return a.container.UpdaterService.Invoke()
}

func (a *App) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   a.config.Basic.Name,
		Message: message,
	})
}

// 構造体のバインド用のメソッド.
func (a *App) EmptyBattle() data.Battle {
	return data.Battle{}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogSetLogLevel(ctx, logger.INFO)

	if isAlreadyRunning() {
		a.showExistDialog("すでに起動しています。", 1)
	}

	a.container = di.NewContainer(ctx, a.config)
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
