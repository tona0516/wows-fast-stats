package main

import (
	"embed"
	"fmt"
	"os"
	"time"
	"wfs/backend/infra"
	"wfs/backend/service"
	"wfs/backend/vo"

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
	semver            string
	revision          string
	env               string
	discordWebhookURL string
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
	version := vo.Version{Semver: semver, Revision: revision}

	// infra
	var maxRetry uint64 = 2
	logger := infra.NewLogger(vo.Env{Str: env}, version)
	discordRepo := infra.NewDiscord(vo.RequestConfig{
		URL:   discordWebhookURL,
		Retry: maxRetry,
	})
	wargamingRepo := infra.NewWargaming(vo.RequestConfig{
		URL:   "https://api.worldofwarships.asia",
		Retry: maxRetry,
	})
	numbersRepo := infra.NewNumbers(vo.RequestConfig{
		URL:   "https://api.wows-numbers.com/personal/rating/expected/json/",
		Retry: maxRetry,
	})
	tempArenaInfoRepo := infra.NewTempArenaInfo()
	configRepo := infra.NewConfig()
	screenshotRepo := infra.NewScreenshot()
	unregisteredRepo := infra.NewUnregistered()
	githubRepo := infra.NewGithub(vo.RequestConfig{
		URL:   "https://api.github.com",
		Retry: maxRetry,
	})

	// service
	var parallels uint = 5
	interval := 1 * time.Second
	configService := service.NewConfig(configRepo, wargamingRepo)
	screenshotService := service.NewScreenshot(screenshotRepo, runtime.SaveFileDialog)
	battleService := service.NewBattle(parallels, wargamingRepo, tempArenaInfoRepo, numbersRepo, unregisteredRepo)
	watcherService := service.NewWatcher(interval, configRepo, tempArenaInfoRepo, runtime.EventsEmit)
	reportService := service.NewReport(discordRepo, configRepo, tempArenaInfoRepo)
	updaterService := service.NewUpdater(version, githubRepo)

	return NewApp(
		vo.Env{Str: env},
		vo.Version{Semver: semver, Revision: revision},
		*configService,
		*screenshotService,
		*watcherService,
		*battleService,
		*reportService,
		*updaterService,
		logger,
	)
}

func isAlreadyRunning() bool {
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
