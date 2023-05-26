package main

import (
	"changeme/backend/infra"
	"changeme/backend/service"
	"changeme/backend/vo"
	"embed"
	"fmt"
	"os"

	"github.com/mitchellh/go-ps"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var (
	semver   string
	revision string
	env      string
)

func main() {
	app := initApp()

	if isAlreadyRunning() {
		os.Exit(0)
		return
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wows-fast-stats",
		Width:  1280,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.onStartup,
		OnShutdown:       app.onShutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func initApp() *App {
	// infra
	wargamingRepo := infra.NewWargaming()
	numbersRepo := infra.NewNumbers("https://api.wows-numbers.com/personal/rating/expected/json/")
	tempArenaInfoRepo := infra.NewTempArenaInfo()
	configRepo := infra.NewConfig()
	screenshotRepo := infra.NewScreenshot()
	unregisteredRepo := infra.NewUnregistered()
	caches := infra.NewCaches("cache")
	logger := infra.NewLogger(vo.Env{Str: env}, vo.Version{Semver: semver, Revision: revision})

	// service
	var parallels uint = 5
	configService := service.NewConfig(configRepo, wargamingRepo)
	screenshotService := service.NewScreenshot(screenshotRepo, runtime.SaveFileDialog)
	prepareService := service.NewPrepare(parallels, wargamingRepo, numbersRepo, unregisteredRepo, *caches)
	battleService := service.NewBattle(parallels, wargamingRepo, tempArenaInfoRepo, *caches, *prepareService)
	replayWatcher := service.NewReplayWatcher(configRepo, tempArenaInfoRepo, runtime.EventsEmit)

	return NewApp(
		vo.Env{Str: env},
		vo.Version{Semver: semver, Revision: revision},
		*configService,
		*screenshotService,
		*replayWatcher,
		*prepareService,
		*battleService,
		*logger,
	)
}

func isAlreadyRunning() bool {
	ownPid := os.Getpid()
	ownPidInfo, err := ps.FindProcess(ownPid)
	if err != nil {
		// Note: for availability.
		return false
	}

	processes, err := ps.Processes()
	if err != nil {
		// Note: for availability.
		return false
	}

	for _, p := range processes {
		if p.Pid() != ownPid && p.Executable() == ownPidInfo.Executable() {
			return true
		}
	}

	return false
}
