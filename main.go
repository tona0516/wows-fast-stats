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
var env string

func main() {
	config := di.NewConfig(env)
	app := NewApp(config)

	err := wails.Run(&options.App{
		Title:     config.App.Name,
		Width:     config.App.Width,
		Height:    config.App.Height,
		MinWidth:  config.App.MinWidth,
		MinHeight: config.App.MinHeight,
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
