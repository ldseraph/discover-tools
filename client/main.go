package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

type Module interface {
	OnStartup(ctx context.Context)
}

var binds []interface{}

func RegistryModule(m Module) {
	binds = append(binds, m)
}

func main() {

	err := wails.Run(&options.App{
		Title:  "设备发现工具",
		Width:  1024,
		Height: 768,
		Assets: assets,
		OnStartup: func(ctx context.Context) {
			for _, b := range binds {
				m, isModule := b.(Module)
				if isModule {
					m.OnStartup(ctx)
				}
			}
		},
		Frameless: true,
		Bind:      binds,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
