package main

import (
	"changeme/backend/repo"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	config vo.UserConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configAdapter := repo.ConfigAdapter{}
	config, _ := configAdapter.Read()
	a.config = config
}

func (a *App) GetTempArenaInfoHash() (string, error) {
	statsService := service.StatsService{
		InstallPath: a.config.InstallPath,
		AppID:       a.config.Appid,
		Parallels:   5,
	}
	return statsService.GetTempArenaInfoHash()
}

func (a *App) Load() ([]vo.Team, error) {
	statsService := service.StatsService{
		InstallPath: a.config.InstallPath,
		AppID:       a.config.Appid,
		Parallels:   5,
	}

	return statsService.GetsStats()
}

func (a *App) Debug(message string) {
	logger.NewDefaultLogger().Debug(message)
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
}

func (a *App) GetConfig() (vo.UserConfig, error) {
	configService := service.ConfigService{}
	config, err := configService.Read()
	if err != nil {
		return config, err
	}

	return config, nil
}

func (a *App) ApplyConfig(installPath string, appid string, fontSize string) (vo.UserConfig, error) {
	configService := service.ConfigService{}
	config := vo.UserConfig{
		InstallPath: installPath,
		Appid:       appid,
		FontSize:    fontSize,
	}

	updatedConfig, err := configService.Update(config)
	if err != nil {
		return updatedConfig, err
	}

	a.config = updatedConfig
	return updatedConfig, nil
}
