package main

import (
	"changeme/backend/domain"
	"changeme/backend/infra"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PARALLELS = 5

type App struct {
	Version           vo.Version
	Env               vo.Env
	ctx               context.Context
	userConfig        vo.UserConfig
	appConfig         vo.AppConfig
	excludePlayer     domain.ExcludePlayer
	isFinishedPrepare bool
	logger            Logger
	cancelWatch       context.CancelFunc
}

func NewApp(env vo.Env, version vo.Version) *App {
	return &App{
		Env:           env,
		Version:       version,
		excludePlayer: *domain.NewExcludePlayer(),
		logger:        *NewLogger(env, version),
	}
}

func (a *App) startup(ctx context.Context) {
	a.logger.Info("start app.")
	a.ctx = ctx

	// Set window size
	window := a.appConfig.Window
	if window.Width > 0 && window.Height > 0 {
		runtime.WindowSetSize(ctx, window.Width, window.Height)
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	width, height := runtime.WindowGetSize(ctx)
	a.appConfig.Window.Width = width
	a.appConfig.Window.Height = height
	configService := service.Config{}
	err := configService.UpdateApp(a.appConfig)
	if err != nil {
		a.logger.Warn("Failed to update app config.", err)
	}

	return false
}

func (a *App) Ready() {
	if a.cancelWatch != nil {
		a.cancelWatch()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelWatch = cancel

	rw := service.NewReplayWatcher(a.ctx, infra.Config{}, infra.TempArenaInfo{})
	go rw.Start(ctx)
}

func (a *App) IsFinishedPrepare() bool {
	return a.isFinishedPrepare
}

func (a *App) Prepare() error {
	// Read configs
	configService := service.Config{}
	userConfig, err := configService.User()
	if err != nil {
		a.logger.Info("No user config.")
	}
	a.userConfig = userConfig

	appConfig, err := configService.App()
	if err != nil {
		a.logger.Info("No app config.")
	}
	a.appConfig = appConfig

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

	result, tempArenaInfo, err := battle.Battle()
	if err != nil {
		a.logger.Error("Failed to get battle.", err)
		a.logger.Info(tempArenaInfo.ToBase64())
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
	configService := service.Config{}
	return configService.User()
}

func (a *App) ApplyUserConfig(config vo.UserConfig) error {
	configService := service.Config{}

	if err := configService.UpdateUser(config); err != nil {
		a.logger.Error("Failed to apply user config.", err)
		return err
	}

	a.userConfig = config
	return nil
}

func (a *App) ManualScreenshot(filename string, base64Data string) error {
	screenshotService := service.Screenshot{}
	err := screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Error("Failed to save screenshot by manual.", err)
	}
	return err
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	screenshotService := service.Screenshot{}
	err := screenshotService.SaveForAuto(filename, base64Data)
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
	return cwd, err
}

func (a *App) AppVersion() vo.Version {
	return a.Version
}

func (a *App) OpenDirectory(path string) error {
	err := open.Run(path)
	if err != nil {
		a.logger.Warn("Failed to open directory.", err)
	}
	return err
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
