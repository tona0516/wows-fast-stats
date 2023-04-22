package main

import (
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"

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
	a.userConfig, _ = a.configService.ReadUserConfig()
	a.appConfig, _ = a.configService.ReadAppConfig()

	window := a.appConfig.Window
	if window.Width != 0 && window.Height != 0 {
		runtime.WindowSetSize(ctx, window.Width, window.Height)
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	width, height := runtime.WindowGetSize(ctx)
	a.appConfig.Window.Width = width
	a.appConfig.Window.Height = height
	a.configService.UpdateAppConfig(a.appConfig)

	return false
}

func (a *App) GetTempArenaInfoHash() (string, error) {
	statsService := service.StatsService{
		Parallels: 5,
	}
	return statsService.GetTempArenaInfoHash(a.userConfig.InstallPath)
}

func (a *App) GetBattle() (vo.Battle, error) {
	statsService := service.StatsService{
		Parallels: 5,
	}

	return statsService.GetsBattle(a.userConfig.InstallPath, a.userConfig.Appid)
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
}

func (a *App) GetConfig() (vo.UserConfig, error) {
	configService := service.ConfigService{}
	return configService.ReadUserConfig()
}

func (a *App) ApplyConfig(config vo.UserConfig) (vo.UserConfig, error) {
	configService := service.ConfigService{}

	updatedConfig, err := configService.UpdateUserConfig(config)
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
