package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"wfs/backend/data"
	"wfs/backend/infra"
	"wfs/backend/infra/webapi"
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
	WGAppID                string
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
		WGAppID: WGAppID,
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

	alertDiscord := webapi.NewDiscord(webapi.RequestConfig{
		URL:     AlertDiscordWebhookURL,
		Retry:   maxRetry,
		Timeout: timeout,
	})
	infoDiscord := webapi.NewDiscord(webapi.RequestConfig{
		URL:     InfoDiscordWebhookURL,
		Retry:   maxRetry,
		Timeout: timeout,
	})

	options := badger.DefaultOptions("./persistent_data")
	if env.IsDev {
		options = options.WithBypassLockGuard(true)
	}
	db, err := badger.Open(options)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	storage := infra.NewStorage(db)

	logger := infra.NewLogger(env, alertDiscord, infoDiscord, *storage)

	wargaming := webapi.NewWargaming(webapi.RequestConfig{
		URL:     "https://api.worldofwarships.asia",
		Retry:   maxRetry,
		Timeout: timeout,
	}, ratelimit.New(10), WGAppID)
	uwargaming := webapi.NewUnofficialWargaming(webapi.RequestConfig{
		URL:     "https://clans.worldofwarships.asia",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	numbers := webapi.NewNumbers(webapi.RequestConfig{
		URL:     "https://api.wows-numbers.com",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	localFile := infra.NewLocalFile()
	unregistered := infra.NewUnregistered()
	github := webapi.NewGithub(webapi.RequestConfig{
		URL:     "https://api.github.com",
		Retry:   maxRetry,
		Timeout: timeout,
	})
	warshipStore := infra.NewWarshipStore(
		db,
		wargaming,
		*unregistered,
		numbers,
	)
	clanFercher := infra.NewClanFetcher(
		wargaming,
		uwargaming,
	)
	rawStatFetcher := infra.NewRawStatFetcher(wargaming)
	battleMetaFetcher := infra.NewBattleMetaFetcher(wargaming)
	accountFetcher := infra.NewAccountFetcher(wargaming)
	userConfig := infra.NewUserConfigStore(db)
	alertPlayer := infra.NewAlertPlayerStore(db)
	versionFetcher := infra.NewVersionFetcher(github)

	// usecase
	watchInterval := 1 * time.Second
	config := service.NewConfig(accountFetcher, userConfig, alertPlayer)
	screenshot := service.NewScreenshot(localFile)
	battle := service.NewBattle(
		localFile,
		warshipStore,
		clanFercher,
		rawStatFetcher,
		battleMetaFetcher,
		accountFetcher,
		logger,
	)
	watcher := service.NewWatcher(watchInterval, localFile, userConfig, logger, runtime.EventsEmit)
	configMigrator := service.NewConfigMigrator(storage, userConfig, alertPlayer)

	return NewApp(
		env,
		logger,
		versionFetcher,
		*config,
		*screenshot,
		*watcher,
		*battle,
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
