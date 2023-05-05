package main

import (
	"changeme/backend/infra"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/slices"
)

const PARALLELS = 5

type App struct {
	Version         vo.Version
	Env             vo.Env
	ctx             context.Context
	userConfig      vo.UserConfig
	appConfig       vo.AppConfig
	excludePlayerID []int
	isFirstBattle   bool
	logger          Logger
}

func NewApp(env vo.Env, version vo.Version) *App {
	return &App{
		Env:     env,
		Version: version,
	}
}

func (a *App) startup(ctx context.Context) {
	a.logger = *NewLogger(a.Env, a.Version)
	a.ctx = ctx

	a.logger.Info("start app.")

	var err error
	configService := service.Config{}
	a.userConfig, err = configService.User()
	if err != nil {
		a.logger.Info("No user config.")
	}

	a.appConfig, err = configService.App()
	if err != nil {
		a.logger.Info("No app config.")
	}

	a.excludePlayerID = make([]int, 0)
	a.isFirstBattle = true

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

func (a *App) TempArenaInfoHash() (string, error) {
	// Note: no logging because this method is called looper
	battle := service.NewBattle(
		PARALLELS,
		a.userConfig,
		infra.Wargaming{AppID: a.userConfig.Appid},
		infra.TempArenaInfo{},
	)
	return battle.TempArenaInfoHash()
}

func (a *App) Battle() (vo.Battle, error) {
	if a.isFirstBattle {
		prepare := service.NewPrepare(
			PARALLELS,
			infra.Wargaming{AppID: a.userConfig.Appid},
			infra.Numbers{},
			infra.Unregistered{},
		)
		if err := prepare.FetchCachable(); err != nil {
			a.logger.Error("Failed to fetch cachable.", err)
			return vo.Battle{}, err
		}
		a.isFirstBattle = false
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

func (a *App) SaveScreenshot(filename string, base64Data string, isSelectable bool) error {
	screenshotService := service.Screenshot{}
	if isSelectable {
		err := screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
		if err != nil {
			a.logger.Error("Failed to save screenshot.", err)
		}
	}

	err := screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Error("Failed to autosave screenshot.", err)
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
	return a.excludePlayerID
}

func (a *App) AddExcludePlayerID(playerID int) {
	if !slices.Contains(a.excludePlayerID, playerID) {
		a.excludePlayerID = append(a.excludePlayerID, playerID)
	}
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	index := slices.Index(a.excludePlayerID, playerID)
	if index != -1 {
		a.excludePlayerID = append(a.excludePlayerID[:index], a.excludePlayerID[index+1:]...)
	}
}
