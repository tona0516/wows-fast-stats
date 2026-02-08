package main

import (
	"context"
	"errors"
	"os"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/core"
	"wfs/backend/infra"
	"wfs/backend/usecase"
	"wfs/backend/usecase/service"

	"github.com/mitchellh/go-ps"
	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	prefetchUsecase             *usecase.Prefetch
	fetchBattleUsecase          *usecase.FetchBattle
	pollMatchUsecase            *usecase.PollMatch
	selectGameClientPathUsecase *usecase.SelectGameClientPath
	checkUpdateUsecase          *usecase.CheckUpdate
	prefStore                   adapter.PrefStore
	logger                      adapter.Logger

	ctx                 context.Context
	config              config.Config
	pollMatchCancelFunc context.CancelFunc
	prefetchResult      *core.PrefetchResult
}

func NewApp(config config.Config) *App {
	return &App{config: config}
}

func (a *App) Prefetch() error {
	result, err := a.prefetchUsecase.Invoke(a.ctx)
	if err != nil {
		return core.ErrorForDisplay(err)
	}

	a.prefetchResult = result
	return nil
}

func (a *App) StartPollingMatch() {
	if a.pollMatchCancelFunc != nil {
		a.pollMatchCancelFunc()
	}

	cancelCtx, cancelFunc := context.WithCancel(a.ctx)
	a.pollMatchCancelFunc = cancelFunc

	go a.pollMatchUsecase.Invoke(a.ctx, cancelCtx)
}

func (a *App) FetchBattle(tempArenaInfo core.TempArenaInfo) (*core.Battle, error) {
	battle, err := a.fetchBattleUsecase.Invoke(a.ctx, tempArenaInfo, a.prefetchResult)
	if err != nil {
		a.pollMatchCancelFunc()
		errForDisplay := core.ErrorForDisplay(err)
		return nil, errForDisplay
	}

	return battle, nil
}

func (a *App) LoadGameClientPath() (string, error) {
	pathString, err := a.prefStore.GameClientPath()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", core.ErrorForDisplay(err)
	}

	return pathString, nil
}

func (a *App) LoadDisplayPref() (string, error) {
	jsonString, err := a.prefStore.DisplayPref()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", core.ErrorForDisplay(err)
	}

	return jsonString, nil
}

func (a *App) SaveDisplayPref(pref string) error {
	return a.prefStore.SetDisplayPref(pref)
}

func (a *App) SelectGameClientPath() (string, error) {
	selectedPath, err := a.selectGameClientPathUsecase.Invoke(a.ctx)
	if err != nil {
		return "", core.ErrorForDisplay(err)
	}

	a.StartPollingMatch()
	return selectedPath, nil
}

func (a *App) CurrentVersion() string {
	return a.config.Basic.Version
}

func (a *App) NewVersion() *core.NewVersion {
	return a.checkUpdateUsecase.Invoke(a.ctx)
}

func (a *App) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   a.config.Basic.Name,
		Message: message,
	})
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	if a.isAlreadyRunning() {
		a.ShowMessageDialog("すでに起動しています")
		os.Exit(1)
		return
	}

	injector := a.getInjector()
	a.prefetchUsecase = do.MustInvoke[*usecase.Prefetch](injector)
	a.fetchBattleUsecase = do.MustInvoke[*usecase.FetchBattle](injector)
	a.pollMatchUsecase = do.MustInvoke[*usecase.PollMatch](injector)
	a.selectGameClientPathUsecase = do.MustInvoke[*usecase.SelectGameClientPath](injector)
	a.checkUpdateUsecase = do.MustInvoke[*usecase.CheckUpdate](injector)
	a.prefStore = do.MustInvoke[adapter.PrefStore](injector)
	a.logger = do.MustInvoke[adapter.Logger](injector)
}

func (a *App) getInjector() do.Injector {
	injector := do.New()

	do.ProvideValue(injector, a.config)

	// infra
	do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
		return infra.NewWails(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.CacheStore, error) {
		return infra.NewCacheStore(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
		return infra.NewPrefStore(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.ReplayReader, error) {
		return infra.NewReplayReader(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.ClanClient, error) {
		return infra.NewClanClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.GithubClient, error) {
		return infra.NewGithubClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.NumbersClient, error) {
		return infra.NewNumbersClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.WargamingClient, error) {
		return infra.NewWargamingClient(i)
	})
	do.Provide(injector, func(i do.Injector) (adapter.Logger, error) {
		return infra.NewLogger(i)
	})
	do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return infra.NewDiscordClient(
			a.config.DiscordClient.AlertWebhookURL,
			a.config.DiscordClient.RetryCount,
			a.config.DiscordClient.Timeout,
		)
	})
	do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return infra.NewDiscordClient(
			a.config.DiscordClient.InfoWebhookURL,
			a.config.DiscordClient.RetryCount,
			a.config.DiscordClient.Timeout,
		)
	})

	// application service
	do.Provide(injector, service.NewBadgeFetcher)
	do.Provide(injector, service.NewClanFetcher)
	do.Provide(injector, service.NewGameClientPathValidator)
	do.Provide(injector, service.NewStatsFetcher)

	// usecase
	do.Provide(injector, usecase.NewPrefetch)
	do.Provide(injector, usecase.NewFetchBattle)
	do.Provide(injector, usecase.NewPollMatch)
	do.Provide(injector, usecase.NewSelectGameClientPath)
	do.Provide(injector, usecase.NewCheckUpdate)

	return injector
}

func (a *App) isAlreadyRunning() bool {
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
