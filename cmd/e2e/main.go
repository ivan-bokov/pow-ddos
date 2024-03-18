package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ivan-bokov/pow-ddos/config"
	"github.com/ivan-bokov/pow-ddos/internal/app/client"
	"github.com/ivan-bokov/pow-ddos/internal/app/server"
	"github.com/ivan-bokov/pow-ddos/internal/service/pow/hashcash"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(l)
	wg := sync.WaitGroup{}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		cfg, err := config.NewConfig[config.ServerConfig](ctx)
		if err != nil {
			slog.Info("failed to read config", "error", err)
			os.Exit(1)
		}
		srv, err := server.New(cfg)
		if err != nil {
			slog.Info("failed to create server", "error", err)
			os.Exit(1)
		}
		if err = srv.Run(ctx, cfg.Addr); err != nil {
			slog.Error("failed to run server", "error", err)
			cancel()
			return
		}
	}()
	time.Sleep(1 * time.Second)
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := client.New(":8000", hashcash.NewSolver())
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			q, err := c.GetQuote(ctx)
			if err != nil {
				slog.Error("failed to get quote", "error", err)
				return
			}
			slog.Info("got quote", "quote", q)
		}()
	}

	wg.Wait()

}
