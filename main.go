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

	app := NewApp(*config)

	var title string
	if config.App.Name != nil {
		title = *config.App.Name
	}
	var width int
	if config.App.Width != nil {
		width = *config.App.Width
	}
	var height int
	if config.App.Height != nil {
		height = *config.App.Height
	}

	err = wails.Run(&options.App{
		Title:  title,
		Width:  width,
		Height: height,
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
