package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ivan-bokov/pow-ddos/internal/app/client"
	"github.com/ivan-bokov/pow-ddos/internal/service/pow/hashcash"
)

var (
	addr  = flag.String("host", "0.0.0.0:8000", "server address")
	count = flag.Int("count", 10, "number of requests")
)

func main() {
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(l)
	slog.Info("run client", "host", *addr, "count", *count)
	c := client.New(*addr, hashcash.NewSolver())
	for i := 0; i < *count; i++ {
		q, err := c.GetQuote(ctx)
		if err != nil {
			slog.Error("failed to get quote", "error", err)
			return
		}
		slog.Info("got quote", "quote", q)
	}
}
