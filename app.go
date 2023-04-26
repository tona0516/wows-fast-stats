package main

import (
	"changeme/backend/repo"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	configService service.ConfigService
	userConfig    vo.UserConfig
	appConfig     vo.AppConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.configService = service.ConfigService{}
	a.userConfig, _ = a.configService.User()
	a.appConfig, _ = a.configService.App()

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
	a.configService.UpdateApp(a.appConfig)

	// Remove old caches
	wargaming := repo.Wargaming{AppID: a.userConfig.Appid}
	encyclopediaInfo, _ := wargaming.EncyclopediaInfo()
	gameVersion := encyclopediaInfo.Data.GameVersion
	caches := repo.NewCaches(gameVersion)
	caches.RemoveOld()

	return false
}

func (a *App) TempArenaInfoHash() (string, error) {
	statsService := service.StatsService{
		Parallels:  5,
		UserConfig: a.userConfig,
	}
	return statsService.TempArenaInfoHash()
}

func (a *App) Battle() (vo.Battle, error) {
	statsService := service.StatsService{
		Parallels:  5,
		UserConfig: a.userConfig,
	}

	return statsService.Battle()
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
}

func (a *App) UserConfig() (vo.UserConfig, error) {
	configService := service.ConfigService{}
	return configService.User()
}

func (a *App) ApplyUserConfig(config vo.UserConfig) (vo.UserConfig, error) {
	configService := service.ConfigService{}

	updatedConfig, err := configService.UpdateUser(config)
	if err != nil {
		return updatedConfig, err
	}

	a.userConfig = updatedConfig
	return updatedConfig, nil
}

func (a *App) SaveScreenshot(filename string, base64Data string, isSelectable bool) error {
	screenshotService := service.ScreenshotService{}
	if isSelectable {
		return screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	}

	return screenshotService.SaveForAuto(filename, base64Data)
}

func (a *App) Cwd() (string, error) {
	return os.Getwd()
}
