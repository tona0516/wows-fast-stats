package main

import (
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// App struct
type App struct {
	ctx          context.Context
	statsService service.StatsService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.statsService = service.StatsService{
		InstallPath: "./",
		AppID:       "3bd34ff346625bf01cc8ba6a9204dd16",
		Parallels:   5,
	}
}

func (a *App) GetTempArenaInfoHash() string {
	return a.statsService.GetTempArenaInfoHash()
}

func (a *App) Load() vo.Team {
	team, err := a.statsService.GetsStats()
	if err != nil {
		logger.NewDefaultLogger().Fatal(err.Error())
	}

	// logger.NewDefaultLogger().Debug(fmt.Sprintf("%#v", team))

	return *team
}

func (a *App) Debug(message string) {
	logger.NewDefaultLogger().Debug(message)
}
