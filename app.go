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

type App struct {
	Version         vo.Version
	ctx             context.Context
	userConfig      vo.UserConfig
	appConfig       vo.AppConfig
	excludePlayerID []int
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configService := service.Config{}
	a.userConfig, _ = configService.User()
	a.appConfig, _ = configService.App()
	a.excludePlayerID = make([]int, 0)

	window := a.appConfig.Window
	if window.Width != 0 && window.Height != 0 {
		runtime.WindowSetSize(ctx, window.Width, window.Height)
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// save window size
	width, height := runtime.WindowGetSize(ctx)
	a.appConfig.Window.Width = width
	a.appConfig.Window.Height = height
	configService := service.Config{}
	configService.UpdateApp(a.appConfig)

	// Remove old caches
	wargaming := infra.Wargaming{AppID: a.userConfig.Appid}
	encyclopediaInfo, _ := wargaming.EncyclopediaInfo()
	gameVersion := encyclopediaInfo.Data.GameVersion
	caches := infra.NewCaches(gameVersion)
	caches.RemoveOld()

	return false
}

func (a *App) TempArenaInfoHash() (string, error) {
	battle := service.Battle{Parallels: 5, UserConfig: a.userConfig}
	return battle.TempArenaInfoHash()
}

func (a *App) Battle() (vo.Battle, error) {
	battle := service.Battle{Parallels: 5, UserConfig: a.userConfig}
	return battle.Battle()
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
}

func (a *App) UserConfig() (vo.UserConfig, error) {
	configService := service.Config{}
	return configService.User()
}

func (a *App) ApplyUserConfig(config vo.UserConfig) error {
	configService := service.Config{}

	if err := configService.UpdateUser(config); err != nil {
		return err
	}

	a.userConfig = config
	return nil
}

func (a *App) SaveScreenshot(filename string, base64Data string, isSelectable bool) error {
	screenshotService := service.Screenshot{}
	if isSelectable {
		return screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	}

	return screenshotService.SaveForAuto(filename, base64Data)
}

func (a *App) Cwd() (string, error) {
	return os.Getwd()
}

func (a *App) AppVersion() vo.Version {
	return a.Version
}

func (a *App) OpenDirectory(path string) error {
	return open.Run(path)
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
