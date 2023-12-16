package main

import (
	"context"
	"github.com/PopescuStefanRadu/ent-demo/pkg"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server"
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

	srv, err := server.NewHTTPServer(server.Config{
		ShutdownTimeout: 30 * time.Second,
		Address:         ":8080",
		AppConfig: &pkg.Config{
			DBUrl:            "file:ent?mode=memory&cache=shared&_fk=1",
			DebugPersistence: true,
		},
	}, l)

	if err != nil {
		l.Err(err).Msg("Could not create http server")
		return
	}

	if err := srv.Start(ctx); err != nil {
		l.Err(err).Msg("Server closed with an unexpected error")
		return
	}

	l.Info().Msg("Application closed")
}
