// +build desktop

package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Info Collector",
		Width:  500,
		Height: 350,
		Frameless: true,	
		BackgroundColour:  &options.RGBA{R:27,G:38,B:54,A:255}, // Dark blue background
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			DisableFramelessWindowDecorations: true,
			WindowIsTranslucent: true,               // needed for rounded shadow
        	BackdropType: windows.Acrylic,       // needed for rounded shadow
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}