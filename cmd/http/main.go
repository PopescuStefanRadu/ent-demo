package main

import (
	"context"
	"github.com/PopescuStefanRadu/ent-demo/cmd/http/server"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := zerolog.New(os.Stdout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx = l.WithContext(ctx)

	srv := server.NewHTTPServer(server.Config{
		ShutdownTimeout: 30 * time.Second,
		Address:         ":8080",
	}, l)

	if err := srv.Start(ctx); err != nil {
		l.Err(err).Msg("Server closed with an unexpected error")
	}

	l.Info().Msg("Application closed")
}
