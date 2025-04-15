package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "Tasklight", // App title
		Width:             600,         // Spotlight-style width
		Height:            200,         // Spotlight-style height
		Frameless:         true,        // Frameless window
		DisableResize:     true,        // Non-resizable window
		AlwaysOnTop:       true,        // Keep it always on top
		HideWindowOnClose: true,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 0}, // White background
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			hideAppFromDock() // <-- Hide from dock and CMD+Tab on MacOS
			focusAppWindow()  // <-- Focus Window after hiding
			app.startup(ctx)  // <-- Main Startup Entry
		},
		Bind: []interface{}{
			app, // Bind the app struct for function calls
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
			TitleBar: &mac.TitleBar{
				HideTitle:    true,
				HideTitleBar: true,
			},
			Appearance: mac.NSAppearanceNameDarkAqua,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
