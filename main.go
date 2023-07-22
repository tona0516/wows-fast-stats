package main

import (
	"embed"
	"fmt"
	"os"
	"time"
	"wfs/backend/application/service"
	"wfs/backend/application/vo"
	"wfs/backend/infra"

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
	version := vo.Version{Semver: semver}
	env := vo.Env{Str: env}

	// infra
	var maxRetry uint64 = 2
	discord := infra.NewDiscord(infra.RequestConfig{
		URL:   discordWebhookURL,
		Retry: maxRetry,
	})
	wargaming := infra.NewWargaming(infra.RequestConfig{
		URL:   "https://api.worldofwarships.asia",
		Retry: maxRetry,
	})
	numbers := infra.NewNumbers(infra.RequestConfig{
		URL:   "https://api.wows-numbers.com/personal/rating/expected/json/",
		Retry: maxRetry,
	})
	localFile := infra.NewLocalFile()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(infra.RequestConfig{
		URL:   "https://api.github.com",
		Retry: maxRetry,
	})

	// service
	var parallels uint = 5
	watchInterval := 1 * time.Second
	configService := service.NewConfig(localFile, wargaming)
	screenshotService := service.NewScreenshot(localFile, runtime.SaveFileDialog)
	battleService := service.NewBattle(parallels, wargaming, localFile, numbers, unregistered)
	watcherService := service.NewWatcher(watchInterval, localFile, runtime.EventsEmit)
	reportService := service.NewReport(version, localFile, discord)
	updaterService := service.NewUpdater(version, github)
	logger := service.NewLogger(env, version, runtime.EventsEmit, *reportService)

	return NewApp(
		env,
		version,
		*configService,
		*screenshotService,
		*watcherService,
		*battleService,
		*reportService,
		*updaterService,
		*logger,
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
