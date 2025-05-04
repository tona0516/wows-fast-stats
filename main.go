package main

import (
	"embed"
	"encoding/base64"
	"log"
	"os"

	"github.com/mitchellh/go-ps"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var Base64ConfigYml string

func main() {
	if isAlreadyRunning() {
		os.Exit(0)
		return
	}

	config, err := getConfig()
	if err != nil {
		log.Fatalln(err)
	}

	app := NewApp(*config)

	err = wails.Run(&options.App{
		Title:  config.App.Name,
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
		log.Fatalln(err)
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

func getConfig() (*Config, error) {
	configYml, err := base64.StdEncoding.DecodeString(Base64ConfigYml)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(configYml, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
