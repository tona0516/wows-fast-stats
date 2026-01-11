package main

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/data"
	"wfs/backend/usecase"

	"github.com/mitchellh/go-ps"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	appConfig                 config.Config
	configStore               adapter.ConfigStore
	prefetchUsecase           *usecase.Prefetch
	fetchBattleUsecase        *usecase.FetchBattle
	pollMatchUsecase          *usecase.PollMatch
	installPathSettingUsecase *usecase.InstallPathSetting
	updateCheckUsecase        *usecase.UpdateCheck

	ctx                 context.Context
	pollMatchCancelFunc context.CancelFunc
	prefetchResult      *data.PrefetchResult
}

func NewApp(i do.Injector) (*App, error) {
	return &App{
		appConfig:                 do.MustInvoke[config.Config](i),
		configStore:               do.MustInvoke[adapter.ConfigStore](i),
		prefetchUsecase:           do.MustInvoke[*usecase.Prefetch](i),
		fetchBattleUsecase:        do.MustInvoke[*usecase.FetchBattle](i),
		pollMatchUsecase:          do.MustInvoke[*usecase.PollMatch](i),
		installPathSettingUsecase: do.MustInvoke[*usecase.InstallPathSetting](i),
		updateCheckUsecase:        do.MustInvoke[*usecase.UpdateCheck](i),
	}, nil
}

func (a *App) Prefetch() {
	result, err := a.prefetchUsecase.Invoke()
	if err != nil {
		return
	}
	a.prefetchResult = result
}

func (a *App) StartPollingMatch() {
	if a.pollMatchCancelFunc != nil {
		a.pollMatchCancelFunc()
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	a.pollMatchCancelFunc = cancelFunc
	channel := make(chan data.TempArenaInfo)

	go a.pollMatchUsecase.Invoke(a.ctx, cancelCtx, channel)
	for tempArenaInfo := range channel {
		a.fetchBattleUsecase.Invoke(a.ctx, tempArenaInfo, a.prefetchResult)
	}
}

func (a *App) GetUserConfig() (data.UserConfig, error) {
	config, err := a.configStore.UserConfig()
	if err == nil {
		return config, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return data.DefaultUserConfig(), nil
	}

	return data.UserConfig{}, failure.Wrap(err)
}

func (a *App) SaveUserConfig(config data.UserConfig) error {
	if err := a.configStore.SetUserConfig(config); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (a *App) TrySaveInstallPath() (bool, error) {
	return a.installPathSettingUsecase.Invoke(a.ctx)
}

func (a *App) CurrentVersion() string {
	return a.appConfig.Basic.Version
}

func (a *App) NewVersion() *data.NewVersion {
	return a.updateCheckUsecase.Invoke()
}

func (a *App) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   a.appConfig.Basic.Name,
		Message: message,
	})
}

// 構造体のバインド用のメソッド.
func (a *App) EmptyBattle() data.Battle {
	return data.Battle{}
}

func (a *App) OnStartup(ctx context.Context) {
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
