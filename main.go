package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"
	"time"
	"wfs/backend/data"
	"wfs/backend/infra"
	"wfs/backend/service"

	"github.com/mitchellh/go-ps"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/ratelimit"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var (
	AppName                string
	Semver                 string
	IsDev                  string
	AlertDiscordWebhookURL string
	InfoDiscordWebhookURL  string
)

func main() {
	if isAlreadyRunning() {
		os.Exit(0)
		return
	}

	isDev, _ := strconv.ParseBool(IsDev)
	env := data.Env{
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
		OnStartup: app.onStartup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func initApp(env data.Env) *App {
	timeout := 5 * time.Second
	retryCount := 2

	alertDiscord := infra.NewDiscord(*infra.NewAPIConfig(
		AlertDiscordWebhookURL,
		timeout,
		retryCount,
	))
	infoDiscord := infra.NewDiscord(*infra.NewAPIConfig(
		InfoDiscordWebhookURL,
		timeout,
		retryCount,
	))
	localFile := infra.NewLocalFile()

	ownIGN, _ := localFile.IGN()

	logger := infra.NewLogger(env, alertDiscord, infoDiscord)
	logger.SetOwnIGN(ownIGN)

	wargaming := infra.NewWargaming(*infra.NewAPIConfig(
		"https://api.worldofwarships.asia",
		timeout,
		retryCount,
	), ratelimit.New(10))
	uwargaming := infra.NewUnofficialWargaming(*infra.NewAPIConfig(
		"https://clans.worldofwarships.asia",
		timeout,
		retryCount,
	))
	numbers := infra.NewNumbers(*infra.NewAPIConfig(
		"https://api.wows-numbers.com",
		timeout,
		retryCount,
	))
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(*infra.NewAPIConfig(
		"https://api.github.com",
		timeout,
		retryCount,
	))

	// usecase
	watchInterval := 1 * time.Second
	config := service.NewConfig(localFile, wargaming, logger)
	screenshot := service.NewScreenshot(localFile, logger)
	battle := service.NewBattle(
		wargaming,
		uwargaming,
		localFile,
		numbers,
		unregistered,
		logger,
		runtime.EventsEmit,
	)
	watcher := service.NewWatcher(watchInterval, localFile, logger, runtime.EventsEmit)
	updater := service.NewUpdater(env, github, logger)

	return NewApp(
		env,
		logger,
		*config,
		*screenshot,
		*watcher,
		*battle,
		*updater,
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
