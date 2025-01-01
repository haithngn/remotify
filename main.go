package main

import (
	"embed"
	"notification-deployer/internal/data"
	"notification-deployer/internal/domain/values"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Remotify",
		MinWidth:  1024,
		MinHeight: 748,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		HideWindowOnClose: true,
		OnStartup:         app.startup,
		Bind: []interface{}{
			app,
		},
		EnumBind: []interface{}{
			data.AppErrorCodes,
			values.FCMDeviceTypes,
			values.APNSPriorities,
			values.APNSPushTypes,
			values.PayloadTemplates,
			values.AppBundle,
		},
		EnableDefaultContextMenu: false,
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Remotify",
				Message: "Ultimate APNS/FCM debugging tool",
			},
		},
		Menu: app.getMenus(),
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:     true,
			DisableWebViewDrop: true,
		},
		OnDomReady: app.domReady,
	})

	if err != nil {
		panic(err)
	}
}
