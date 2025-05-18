package main

import (
	"embed"
	"encoding/base64"
	"log"

	"github.com/goccy/go-yaml"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var Base64ConfigYml string

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	app := NewApp(config)

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
		log.Fatalln(err.Error())
	}
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
