package main

import (
	"changeme/backend/apperr"
	"changeme/backend/domain"
	"changeme/backend/infra"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PARALLELS = 5

type App struct {
	version             vo.Version
	env                 vo.Env
	ctx                 context.Context
	userConfig          vo.UserConfig
	appConfig           vo.AppConfig
	excludePlayer       domain.ExcludePlayer
	isFinishedPrepare   bool
	logger              Logger
	cancelReplayWatcher context.CancelFunc
	configService       service.Config
	screenshotService   service.Screenshot
}

func NewApp(env vo.Env, version vo.Version) *App {
	return &App{
		env:               env,
		version:           version,
		excludePlayer:     *domain.NewExcludePlayer(),
		logger:            *NewLogger(env, version),
		configService:     *service.NewConfig(infra.Config{}),
		screenshotService: *service.NewScreenshot(infra.Screenshot{}),
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.logger.Info("start app.")
	a.ctx = ctx

	// Read configs
	userConfig, err := a.configService.User()
	if err != nil {
		a.logger.Info("No user config.")
	}
	a.userConfig = userConfig

	appConfig, err := a.configService.App()
	if err != nil {
		a.logger.Info("No app config.")
	}
	a.appConfig = appConfig

	// Set window size
	window := a.appConfig.Window
	if window.Width > 0 && window.Height > 0 {
		runtime.WindowSetSize(ctx, window.Width, window.Height)
	}
}

func (a *App) onShutdown(ctx context.Context) {
	a.logger.Info("shutdown app.")

	// Save windows size
	width, height := runtime.WindowGetSize(ctx)
	a.appConfig.Window.Width = width
	a.appConfig.Window.Height = height
	err := a.configService.UpdateApp(a.appConfig)
	if err != nil {
		a.logger.Warn("Failed to update app config.", err)
	}
}

func (a *App) Ready() {
	if a.cancelReplayWatcher != nil {
		a.cancelReplayWatcher()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelReplayWatcher = cancel

	rw := service.NewReplayWatcher(a.ctx, infra.Config{}, infra.TempArenaInfo{})
	go rw.Start(ctx)
}

func (a *App) IsFinishedPrepare() bool {
	return a.isFinishedPrepare
}

func (a *App) Prepare() error {
	// Fetch cachable
	prepare := service.NewPrepare(
		PARALLELS,
		infra.Wargaming{AppID: a.userConfig.Appid},
		infra.Numbers{},
		infra.Unregistered{},
	)

	if err := prepare.FetchCachable(); err != nil {
		return err
	}
	a.isFinishedPrepare = true

	return nil
}

func (a *App) Battle() (vo.Battle, error) {
	if !a.isFinishedPrepare {
		if err := a.Prepare(); err != nil {
			a.logger.Error("Failed to prepare", err)

			return vo.Battle{}, err
		}
	}

	battle := service.NewBattle(
		PARALLELS,
		a.userConfig,
		infra.Wargaming{AppID: a.userConfig.Appid},
		infra.TempArenaInfo{},
	)

	result, err := battle.Battle()
	if err != nil {
		a.logger.Error("Failed to get battle.", err)
	}

	return result, err
}

func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		a.logger.Error("Failed to get client installing path.", err)
	}

	return path, err
}

func (a *App) UserConfig() (vo.UserConfig, error) {
	// Note: no logging because this method is called looper
	return a.configService.User()
}

func (a *App) ApplyUserConfig(config vo.UserConfig) error {
	if err := a.configService.UpdateUser(config); err != nil {
		a.logger.Error("Failed to apply user config.", err)

		return err
	}
	a.userConfig = config

	return nil
}

func (a *App) ManualScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Error("Failed to save screenshot by manual.", err)
	}

	return err
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Error("Failed to save screenshot by auto.", err)
	}

	return err
}

func (a *App) Cwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		a.logger.Warn("Failed to get cwd.", err)
	}

	return cwd, errors.WithStack(apperr.App.Cwd.WithRaw(err))
}

func (a *App) AppVersion() vo.Version {
	return a.version
}

func (a *App) OpenDirectory(path string) error {
	err := open.Run(path)
	if err != nil {
		a.logger.Warn("Failed to open directory.", err)
	}

	return errors.WithStack(apperr.App.OpenDir.WithRaw(err))
}

func (a *App) ExcludePlayerIDs() []int {
	return a.excludePlayer.Get()
}

func (a *App) AddExcludePlayerID(playerID int) {
	a.excludePlayer.Add(playerID)
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	a.excludePlayer.Remove(playerID)
}
