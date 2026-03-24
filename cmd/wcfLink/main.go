package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"wcfLink/internal/app"
	"wcfLink/internal/config"
)

func main() {
	cfg := config.Load()

	level := new(slog.LevelVar)
	level.Set(cfg.LogLevel())
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", "err", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		logger.Error("application exited with error", "err", err)
		os.Exit(1)
	}
}
