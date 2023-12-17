package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"
	"time"
	"wfs/backend/application/usecase"
	"wfs/backend/application/vo"
	"wfs/backend/infra"

	"github.com/dgraph-io/badger/v4"
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
		OnStartup: app.onStartup,
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
	timeout := 10 * time.Second
	wargaming := infra.NewWargaming(infra.RequestConfig{
		URL:     "https://api.worldofwarships.asia",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	numbers := infra.NewNumbers(infra.RequestConfig{
		URL:     "https://api.wows-numbers.com",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	localFile := infra.NewLocalFile()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(infra.RequestConfig{
		URL:     "https://api.github.com",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	discord := infra.NewDiscord(infra.RequestConfig{
		URL:     DiscordWebhookURL,
		Retry:   maxRetry,
		Timeout: timeout,
	})
	db, err := badger.Open(badger.DefaultOptions("./persistent_data"))
	storage := infra.NewStorage(db)
	report := infra.NewReport(env, *localFile, *discord, *storage)

	if err != nil {
		report.Send("fatal has occurred!", err)
		panic(err)
	}

	// usecase
	var parallels uint = 5
	watchInterval := 1 * time.Second
	config := usecase.NewConfig(localFile, wargaming, storage)
	screenshot := usecase.NewScreenshot(localFile)
	battle := usecase.NewBattle(parallels, wargaming, localFile, numbers, unregistered, storage)
	watcher := usecase.NewWatcher(watchInterval, localFile, storage, runtime.EventsEmit)
	updater := usecase.NewUpdater(env, github)
	configMigrator := usecase.NewConfigMigrator(localFile, storage)

	return NewApp(
		env,
		report,
		*config,
		*screenshot,
		*watcher,
		*battle,
		*updater,
		*configMigrator,
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
