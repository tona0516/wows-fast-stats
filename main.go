package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"wfs/backend/domain/model"
	"wfs/backend/infra"
	"wfs/backend/usecase"

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
	env := model.Env{
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

func initApp(env model.Env) *App {
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
	})
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
	var parallels uint = 5
	watchInterval := 1 * time.Second
	config := usecase.NewConfig(localFile, wargaming, storage, logger)
	screenshot := usecase.NewScreenshot(localFile, logger)
	battle := usecase.NewBattle(
		parallels,
		wargaming,
		uwargaming,
		localFile,
		numbers,
		unregistered,
		storage,
		logger,
		runtime.EventsEmit,
	)
	watcher := usecase.NewWatcher(watchInterval, localFile, storage, logger, runtime.EventsEmit)
	updater := usecase.NewUpdater(env, github, logger)
	configMigrator := usecase.NewConfigMigrator(configV0, storage, logger)

	return NewApp(
		env,
		logger,
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
