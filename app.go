package main

import (
	"context"
	"fmt"
	"os"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/core"
	"wfs/backend/infra"
	"wfs/backend/usecase"

	"github.com/mitchellh/go-ps"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	prefetchUsecase           *usecase.Prefetch
	fetchBattleUsecase        *usecase.FetchBattle
	pollMatchUsecase          *usecase.PollMatch
	installPathSettingUsecase *usecase.InstallPathSetting
	updateCheckUsecase        *usecase.UpdateCheck
	loadPrefUsecase           *usecase.LoadPref
	savePrefUsecase           *usecase.SavePref
	logger                    adapter.Logger

	ctx                 context.Context
	config              config.Config
	pollMatchCancelFunc context.CancelFunc
	prefetchResult      *core.PrefetchResult
}

func NewApp(config config.Config) *App {
	return &App{config: config}
}

func (a *App) Prefetch() {
	runtime.EventsEmit(a.ctx, usecase.EventOnStartPrefetch, "艦・マップ情報を取得中")

	result, err := a.prefetchUsecase.Invoke(a.ctx)
	if err != nil {
		code, _ := failure.CodeOf(err)
		message := fmt.Sprintf("[%s] 艦・マップ情報に失敗しました\n再起動してください", code.ErrorCode())
		runtime.EventsEmit(a.ctx, usecase.EventOnPrefetchFailure, message)
		// TODO: キャッシュの削除
		return
	}
	a.prefetchResult = result
}

func (a *App) StartPollingMatch() {
	if a.pollMatchCancelFunc != nil {
		a.pollMatchCancelFunc()
	}

	cancelCtx, cancelFunc := context.WithCancel(a.ctx)
	a.pollMatchCancelFunc = cancelFunc

	installPath, err := a.pollMatchUsecase.GetInstallPath()
	if err != nil {
		if failure.Is(err, core.ErrInitialSettingRequired) {
			runtime.EventsEmit(a.ctx, usecase.EventOnPromote, "設定から初期設定をおこなってください")
		} else {
			code, _ := failure.CodeOf(err)
			message := fmt.Sprintf("[%s] 戦闘検知開始に失敗しました\n再起動してください", code.ErrorCode())
			runtime.EventsEmit(a.ctx, usecase.EventOnFetchBattleFailre, message)
		}
		return
	}

	pollingResult := make(chan usecase.PollingResult)

	go a.pollMatchUsecase.Invoke(a.ctx, cancelCtx, installPath, pollingResult)
	runtime.EventsEmit(a.ctx, usecase.EventOnStartPolling, "戦闘開始時に自動的にリロードします")

	for result := range pollingResult {
		if result.Error != nil {
			code, _ := failure.CodeOf(result.Error)
			message := fmt.Sprintf("[%s] 戦闘検知に失敗しました\nリトライしてください", code.ErrorCode())
			runtime.EventsEmit(a.ctx, usecase.EventOnFetchBattleFailre, message)
			continue
		}

		runtime.EventsEmit(a.ctx, usecase.EventOnStartBattle, "戦闘データを読み込み中")
		battle, err := a.fetchBattleUsecase.Invoke(a.ctx, *result.TempArenaInfo, a.prefetchResult)
		if err != nil {
			code, _ := failure.CodeOf(err)
			message := fmt.Sprintf("[%s] 戦闘データの取得に失敗しました\nリトライしてください", code.ErrorCode())
			runtime.EventsEmit(a.ctx, usecase.EventOnFetchBattleFailre, message)
			continue
		}
		runtime.EventsEmit(a.ctx, usecase.EventOnFetchBattleSuccess, battle)
	}
}

func (a *App) LoadPref() (core.Pref, error) {
	return a.loadPrefUsecase.Invoke()
}

func (a *App) SavePref(pref core.Pref) error {
	return a.savePrefUsecase.Invoke(pref)
}

func (a *App) TrySaveInstallPath() (bool, error) {
	return a.installPathSettingUsecase.Invoke(a.ctx)
}

func (a *App) CurrentVersion() string {
	return a.config.Basic.Version
}

func (a *App) NewVersion() *core.NewVersion {
	return a.updateCheckUsecase.Invoke(a.ctx)
}

func (a *App) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   a.config.Basic.Name,
		Message: message,
	})
}

// 構造体のバインド用のメソッド.
func (a *App) EmptyBattle() core.Battle {
	return core.Battle{}
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
	a.installPathSettingUsecase = do.MustInvoke[*usecase.InstallPathSetting](injector)
	a.updateCheckUsecase = do.MustInvoke[*usecase.UpdateCheck](injector)
	a.loadPrefUsecase = do.MustInvoke[*usecase.LoadPref](injector)
	a.savePrefUsecase = do.MustInvoke[*usecase.SavePref](injector)
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

	// service
	do.Provide(injector, usecase.NewStatsService)
	do.Provide(injector, usecase.NewClanService)
	do.Provide(injector, usecase.NewBadgeService)

	// usecase
	do.Provide(injector, usecase.NewPrefetch)
	do.Provide(injector, usecase.NewFetchBattle)
	do.Provide(injector, usecase.NewPollMatch)
	do.Provide(injector, usecase.NewInstallPathSetting)
	do.Provide(injector, usecase.NewUpdateCheck)
	do.Provide(injector, usecase.NewLoadPref)
	do.Provide(injector, usecase.NewSavePref)

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
