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

type App struct {
	config            Config
	ctx               context.Context
	container         *DependencyContainer
	unsubscribeBattle context.CancelFunc
}

func NewApp(config Config) *App {
	return &App{config: config}
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

func (a *App) GetUserConfig() (data.UserConfig, error) {
	return a.container.configService.GetUserConfig()
}

func (a *App) SaveUserConfig(config data.UserConfig) error {
	return a.container.configService.SaveUserConfig(config)
}

func (a *App) TrySaveInstallPath() (bool, error) {
	return a.container.configService.TrySaveInstallPath(a.ctx)
}

func (a *App) OpenDirectory(path string) error {
	err := a.container.configService.OpenDirectory(path)
	if err != nil {
		a.container.logger.Warn(err, nil)
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

func (a *App) NewVersion() *data.NewVersion {
	return a.container.updaterService.Invoke()
}

func (a *App) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   a.config.App.Name,
		Message: message,
	})
}

func (a *App) GetBlackList() (data.BlackList, error) {
	return a.container.blackListService.Get()
}

func (a *App) UpdateBlackList(item data.BlackListItem) error {
	return a.container.blackListService.Update(item)
}

func (a *App) RemoveFromBlackList(accountID int) error {
	return a.container.blackListService.Remove(accountID)
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

	a.container = NewDependencyContainer(ctx, a.config)
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
