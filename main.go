package main

import (
	"embed"
	"wfs/backend/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	appName    string
	appVersion string
	env        string
)

func main() {
	config := config.NewConfig(appName, appVersion, env)
	app := NewApp(config)
	options := &options.App{
		Title:     config.Basic.Name,
		Width:     config.Basic.Width,
		Height:    config.Basic.Height,
		MinWidth:  config.Basic.MinWidth,
		MinHeight: config.Basic.MinHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.OnStartup,
		Bind:      []any{app},
	}
	if err := wails.Run(options); err != nil {
		panic(err)
	}
}
