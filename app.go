package main

import (
	"context"
	"os"
	"wfs/internal/apperr"
	"wfs/internal/data"
	"wfs/internal/di"

	"github.com/mitchellh/go-ps"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx               context.Context
	config            di.Config
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

	if isAlreadyRunning() {
		a.ShowMessageDialog("すでに起動しています。")
		os.Exit(1)
		return
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
