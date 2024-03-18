package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ivan-bokov/pow-ddos/config"
	"github.com/ivan-bokov/pow-ddos/internal/app/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg, err := config.NewConfig[config.ServerConfig](ctx)
	if err != nil {
		slog.Info("failed to read config", "error", err)
		os.Exit(1)
	}
	var lvl slog.LevelVar

	if err = lvl.UnmarshalText([]byte(cfg.Logger.Level)); err != nil {
		slog.Info("failed to read config", "error", err)
		os.Exit(1)
	}
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl.Level(),
	}))

	slog.SetDefault(l)

	srv, err := server.New(cfg)
	if err != nil {
		slog.Info("failed to create server", "error", err)
		os.Exit(1)
	}
	if err = srv.Run(ctx, cfg.Addr); err != nil {
		slog.Info("failed to run server", "error", err)
		os.Exit(1)
	}
}
