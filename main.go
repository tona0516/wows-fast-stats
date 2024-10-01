package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"wfs/backend/data"
	"wfs/backend/infra"
	"wfs/backend/service"

	"github.com/dgraph-io/badger/v4"
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
	// infra
	var maxRetry uint64 = 2
	timeout := 10 * time.Second

	alertDiscord := infra.NewDiscord(infra.RequestConfig{
		URL:     AlertDiscordWebhookURL,
		Retry:   maxRetry,
		Timeout: timeout,
	})
	infoDiscord := infra.NewDiscord(infra.RequestConfig{
		URL:     InfoDiscordWebhookURL,
		Retry:   maxRetry,
		Timeout: timeout,
	})
	db, err := badger.Open(badger.DefaultOptions("./persistent_data"))
	if err != nil {
		logger := infra.NewLogger(env, alertDiscord, infoDiscord)
		logger.Fatal(err, nil)
		return nil
	}

	storage := infra.NewStorage(db)
	ownIGN, _ := storage.OwnIGN()

	logger := infra.NewLogger(env, alertDiscord, infoDiscord)
	logger.SetOwnIGN(ownIGN)

	wargaming := infra.NewWargaming(infra.RequestConfig{
		URL:     "https://api.worldofwarships.asia",
		Retry:   maxRetry,
		Timeout: timeout,
	}, ratelimit.New(10))
	uwargaming := infra.NewUnofficialWargaming(infra.RequestConfig{
		URL:     "https://clans.worldofwarships.asia",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	numbers := infra.NewNumbers(infra.RequestConfig{
		URL:     "https://api.wows-numbers.com",
		Retry:   maxRetry,
		Timeout: timeout,
		// workaround for expired SSL certificate
		Transport: &http.Transport{
			//nolint:gosec
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	localFile := infra.NewLocalFile()
	configV0 := infra.NewConfigV0()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(infra.RequestConfig{
		URL:     "https://api.github.com",
		Retry:   maxRetry,
		Timeout: timeout,
	})

	// usecase
	watchInterval := 1 * time.Second
	config := service.NewConfig(localFile, wargaming, storage, logger)
	battle := service.NewBattle(
		wargaming,
		uwargaming,
		localFile,
		numbers,
		unregistered,
		storage,
		logger,
		runtime.EventsEmit,
	)
	watcher := service.NewWatcher(watchInterval, localFile, storage, logger, runtime.EventsEmit)
	updater := service.NewUpdater(env, github, logger)
	configMigrator := service.NewConfigMigrator(configV0, storage, logger)
	battleHistory := service.NewBattleHistory(storage, logger)

	return NewApp(
		env,
		logger,
		*config,
		*watcher,
		*battle,
		*updater,
		*configMigrator,
		*battleHistory,
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
