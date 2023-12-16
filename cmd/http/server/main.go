package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PopescuStefanRadu/ent-demo/pkg/app"
	"github.com/PopescuStefanRadu/ent-demo/pkg/external/dog"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server"
	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
)

const shutdownTimeout = 30 * time.Second

func main() {
	l := zerolog.New(os.Stdout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx = l.WithContext(ctx)

	srv, err := server.NewHTTPServer(server.Config{
		ShutdownTimeout: shutdownTimeout,
		Address:         ":8080",
		AppConfig: &app.Config{
			DBUrl:            "file:ent?mode=memory&cache=shared&_fk=1",
			DebugPersistence: true,
			DogClientConfig: dog.ClientConfig{
				Enabled: true,
				BaseURL: "https://random.dog",
				CircuitBreakerSettings: gobreaker.Settings{
					Name: "dog",
				},
			},
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
