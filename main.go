//package main
//
//import (
//	"embed"
//
//	"github.com/wailsapp/wails/v2"
//	"github.com/wailsapp/wails/v2/pkg/options"
//	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
//)
//
////go:embed all:frontend/dist
//var assets embed.FS
//
//func main() {
//	// Create an instance of the app structure
//	app := NewApp()
//
//	// Create application with options
//	err := wails.Run(&options.App{
//		Title:  "tasklight",
//		Width:  1024,
//		Height: 768,
//		AssetServer: &assetserver.Options{
//			Assets: assets,
//		},
//		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
//		OnStartup:        app.startup,
//		Bind: []interface{}{
//			app,
//		},
//	})
//
//	if err != nil {
//		println("Error:", err.Error())
//	}
//}

package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//func main() {
//	app := NewApp()
//
//	err := wails.Run(&options.App{
//		Title:  "Tasklight",
//		Width:  1024,
//		Height: 768,
//		AssetServer: &assetserver.Options{
//			Assets: assets,
//		},
//		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
//		OnStartup:        app.startup,
//		Bind: []interface{}{
//			app,
//		},
//	})
//
//	if err != nil {
//		println("Error:", err.Error())
//	}
//}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "Tasklight",                                 // App title
		Width:            600,                                         // Spotlight-style width
		Height:           200,                                         // Spotlight-style height
		Frameless:        true,                                        // Frameless window
		DisableResize:    true,                                        // Non-resizable window
		AlwaysOnTop:      true,                                        // Keep it always on top
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0}, // White background
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup, // Call app's startup function
		Bind: []interface{}{
			app, // Bind the app struct for function calls
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
