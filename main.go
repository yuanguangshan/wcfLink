package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	coreapp "wcfLink/internal/app"
	"wcfLink/internal/config"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := config.Load()

	level := new(slog.LevelVar)
	level.Set(cfg.LogLevel())
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	core, err := coreapp.New(context.Background(), cfg, logger)
	if err != nil {
		logger.Error("failed to initialize core app", "err", err)
		os.Exit(1)
	}

	bridge := NewAppBridge(core, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		_ = core.Shutdown()
	}()

	err = wails.Run(&options.App{
		Title:     "wcfLink",
		Width:     1200,
		Height:    840,
		DisableResize: true,
		MinWidth:  1200,
		MinHeight: 840,
		MaxWidth:  1200,
		MaxHeight: 840,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:         &options.RGBA{R: 245, G: 239, B: 229, A: 1},
		EnableDefaultContextMenu: true,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		OnStartup:  bridge.Startup,
		OnShutdown: bridge.Shutdown,
		Bind: []interface{}{
			bridge,
		},
	})
	if err != nil {
		logger.Error("desktop app exited with error", "err", err)
		os.Exit(1)
	}
}
