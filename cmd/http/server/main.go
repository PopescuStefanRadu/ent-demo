package main

import (
	"context"
	"github.com/PopescuStefanRadu/ent-demo/pkg/app"
	"github.com/PopescuStefanRadu/ent-demo/pkg/external/dog"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server"
	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
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
		AppConfig: &app.Config{
			DBUrl:            "file:ent?mode=memory&cache=shared&_fk=1",
			DebugPersistence: true,
			DogClientConfig: dog.ClientConfig{
				Enabled: true,
				BaseUrl: "https://random.dog",
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
