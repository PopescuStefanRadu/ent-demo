package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"sync/atomic"
	"time"
)

type Config struct {
	ShutdownTimeout time.Duration
	Address         string
}

type HTTPServer struct {
	Config
	Logger zerolog.Logger
	Server *http.Server
	status *atomic.Value
	Gin    *gin.Engine
}

func NewHTTPServer(config Config, logger zerolog.Logger) *HTTPServer {
	router := newRouter()
	srv := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	return &HTTPServer{
		Config: config,
		Logger: logger,
		Server: srv,
		Gin:    router,
	}
}

func (h *HTTPServer) Start(ctx context.Context) error {
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
