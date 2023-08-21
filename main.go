package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"
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
	AppName           string
	Semver            string
	IsDev             string
	DiscordWebhookURL string
)

func main() {
	if isAlreadyRunning() {
		os.Exit(0)
		return
	}

	isDev, _ := strconv.ParseBool(IsDev)
	env := vo.Env{
		AppName: AppName,
		Semver:  Semver,
		IsDev:   isDev,
	}

	app := initApp(env)

	title := AppName
	if isDev {
		title += " [dev]"
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  title,
		Width:  1280,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.onStartup,
		OnShutdown: app.onShutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func initApp(env vo.Env) *App {
	// infra
	var maxRetry uint64 = 2
	wargaming := infra.NewWargaming(infra.RequestConfig{
		URL:   "https://api.worldofwarships.asia",
		Retry: maxRetry,
	})
	numbers := infra.NewNumbers(infra.RequestConfig{
		URL:   "https://api.wows-numbers.com",
		Retry: maxRetry,
	})
	localFile := infra.NewLocalFile()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(infra.RequestConfig{
		URL:   "https://api.github.com",
		Retry: maxRetry,
	})
	discord := infra.NewDiscord(infra.RequestConfig{
		URL:   DiscordWebhookURL,
		Retry: maxRetry,
	})
	report := infra.NewReport(env, *localFile, *discord)

	// service
	var parallels uint = 5
	watchInterval := 1 * time.Second
	configService := service.NewConfig(localFile, wargaming)
	screenshotService := service.NewScreenshot(localFile)
	battleService := service.NewBattle(parallels, wargaming, localFile, numbers, unregistered)
	watcherService := service.NewWatcher(watchInterval, localFile, runtime.EventsEmit)
	updaterService := service.NewUpdater(env, github)

	return NewApp(
		env,
		report,
		*configService,
		*screenshotService,
		*watcherService,
		*battleService,
		*updaterService,
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
