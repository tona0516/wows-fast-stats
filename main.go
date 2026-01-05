package main

import (
	"embed"
	"log"
	"wfs/internal/di"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//nolint:gochecknoglobals
var (
	appName string
	version string
	env     string
)

func main() {
	config := di.NewConfig(appName, version, env)
	container := di.NewContainer(config)
	app := NewApp(&config, container)

	err := wails.Run(&options.App{
		Title:     config.Basic.Name,
		Width:     config.Basic.Width,
		Height:    config.Basic.Height,
		MinWidth:  config.Basic.MinWidth,
		MinHeight: config.Basic.MinHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.onStartup,
		Bind: []any{
			app,
		},
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}
