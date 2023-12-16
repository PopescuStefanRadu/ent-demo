package server

import (
	"context"
	"errors"
	"github.com/PopescuStefanRadu/ent-demo/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Config struct {
	ShutdownTimeout time.Duration
	Address         string
	AppConfig       *app.Config
}

type HTTPServer struct {
	Config
	App    *app.App
	Logger zerolog.Logger
	Server *http.Server
	Gin    *gin.Engine
}

func NewHTTPServer(config Config, logger zerolog.Logger) (*HTTPServer, error) {
	app, err := app.NewAppFromConfig(logger, config.AppConfig)
	if err != nil {
		return nil, err
	}

	router := NewRouter(app)
	srv := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	return &HTTPServer{
		App:    app,
		Config: config,
		Logger: logger,
		Server: srv,
		Gin:    router,
	}, nil
}

func (h *HTTPServer) Start(ctx context.Context) error {
	err := h.App.Init(ctx)
	if err != nil {
		return err
	}

	serverErr := make(chan error)
	defer close(serverErr)
	go func() {
		if err := h.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), h.ShutdownTimeout)
		defer cancel()
		return h.Shutdown(timeout)
	}
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	h.Logger.Info().Msg("Server shutting down")
	return h.Server.Shutdown(ctx)
}
