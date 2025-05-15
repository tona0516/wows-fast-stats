package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"
	"wfs/backend/infra"
	"wfs/backend/service"

	"github.com/dgraph-io/badger/v4"
	"github.com/goccy/go-yaml"
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
var Base64ConfigYml string

func main() {
	var err error

	if isAlreadyRunning() {
		os.Exit(0)
		return
	}

	config, err := getConfig()
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	app := initApp(config)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  config.App.Name,
		Width:  config.App.Width,
		Height: config.App.Height,
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

func getConfig() (Config, error) {
	configYml, err := base64.StdEncoding.DecodeString(Base64ConfigYml)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := yaml.Unmarshal(configYml, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func initApp(config Config) *App {
	alertDiscord := infra.NewDiscord(infra.RequestConfig{
		URL:     config.Discord.AlertURL,
		Retry:   uint64(config.Discord.MaxRetry),
		Timeout: time.Duration(config.Discord.TimeoutSec) * time.Second,
	})
	infoDiscord := infra.NewDiscord(infra.RequestConfig{
		URL:     config.Discord.InfoURL,
		Retry:   uint64(config.Discord.MaxRetry),
		Timeout: time.Duration(config.Discord.TimeoutSec) * time.Second,
	})

	options := badger.DefaultOptions("./persistent_data")
	// TODO: onStartupで初期化するようにする
	if config.Logger.ZerologLogLevel == "debug" {
		options = options.WithBypassLockGuard(true)
	}
	db, err := badger.Open(options)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	storage := infra.NewStorage(db)
	ownIGN, _ := storage.OwnIGN()

	logger := infra.NewLogger(
		config.App.Name,
		config.App.Semver,
		config.Logger.ZerologLogLevel,
		alertDiscord,
		infoDiscord,
	)
	logger.SetOwnIGN(ownIGN)

	wargaming := infra.NewWargaming(
		infra.RequestConfig{
			URL:     config.Wargaming.URL,
			Retry:   uint64(config.Wargaming.MaxRetry),
			Timeout: time.Duration(config.Wargaming.TimeoutSec) * time.Second,
		},
		ratelimit.New(10), // TODO: onStartupで初期化するようにする
		config.Wargaming.AppID,
	)
	uwargaming := infra.NewUnofficialWargaming(infra.RequestConfig{
		URL:     config.UnofficialWargaming.URL,
		Retry:   uint64(config.UnofficialWargaming.MaxRetry),
		Timeout: time.Duration(config.UnofficialWargaming.TimeoutSec) * time.Second,
	})
	numbers := infra.NewNumbers(infra.RequestConfig{
		URL:     config.Numbers.URL,
		Retry:   uint64(config.Numbers.MaxRetry),
		Timeout: time.Duration(config.Numbers.TimeoutSec) * time.Second,
	})
	localFile := infra.NewLocalFile()
	configV0 := infra.NewConfigV0()
	unregistered := infra.NewUnregistered()
	github := infra.NewGithub(infra.RequestConfig{
		URL:     config.Github.URL,
		Retry:   uint64(config.Github.MaxRetry),
		Timeout: time.Duration(config.Github.TimeoutSec) * time.Second,
	})

	// usecase
	watchInterval := 1 * time.Second
	configService := service.NewConfig(localFile, wargaming, storage, logger)
	screenshot := service.NewScreenshot(localFile, logger)
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
	updater := service.NewUpdater(config.App.Semver, github, logger)
	configMigrator := service.NewConfigMigrator(configV0, storage, logger)

	return NewApp(
		config.App.Semver,
		logger,
		*configService,
		*screenshot,
		*watcher,
		*battle,
		*updater,
		*configMigrator,
	)
}
