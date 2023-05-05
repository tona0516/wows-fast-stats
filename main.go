package main

import (
	"changeme/backend/vo"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	semver   string
	revision string
	env      string
)

func main() {
	// Create an instance of the app structure
	app := NewApp(
		vo.Env{Str: env},
		vo.Version{Semver: semver, Revision: revision},
	)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wows-fast-stats",
		Width:  1920,
		Height: 1080,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
