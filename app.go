package main

import (
	"context"
	"os"
	"wfs/internal/data"
	"wfs/internal/di"

	"github.com/mitchellh/go-ps"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	config    *di.Config
	container *di.Container

	ctx                 context.Context
	pollMatchCancelFunc context.CancelFunc
}

func NewApp(
	config *di.Config,
	container *di.Container,
) *App {
	return &App{config: config, container: container}
}

func (a *App) StartPollingMatch() {
	if a.pollMatchCancelFunc != nil {
		a.pollMatchCancelFunc()
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	a.pollMatchCancelFunc = cancelFunc
	channel := make(chan data.TempArenaInfo)

	go a.container.PollMatchUsecase.Invoke(a.ctx, cancelCtx, channel)
	for tempArenaInfo := range channel {
		a.container.FetchBattleUsecase.Invoke(a.ctx, tempArenaInfo)
	}
}

func (a *App) GetUserConfig() (data.UserConfig, error) {
	return a.container.ConfigService.GetUserConfig()
}

func (a *App) SaveUserConfig(config data.UserConfig) error {
	return a.container.ConfigService.SaveUserConfig(config)
}

func (a *App) TrySaveInstallPath() (bool, error) {
	return a.container.InstallPathSettingUsecase.Invoke(a.ctx)
}

func (a *App) CurrentVersion() string {
	return a.config.Basic.Version
}

func (a *App) NewVersion() *data.NewVersion {
	return a.container.UpdateCheckUsecase.Invoke()
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
