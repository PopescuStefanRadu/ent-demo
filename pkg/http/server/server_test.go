package server_test

import (
	"context"
	"fmt"
	"github.com/PopescuStefanRadu/ent-demo/pkg"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestStartAndGracefulShutdown(t *testing.T) {
	l := zerolog.New(zerolog.NewTestWriter(t))

	getFreePort := func() (port int, err error) {
		var a *net.TCPAddr
		if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
			var l *net.TCPListener
			if l, err = net.ListenTCP("tcp", a); err == nil {
				defer l.Close()
				return l.Addr().(*net.TCPAddr).Port, nil
			}
		}
		return
	}

	port, err := getFreePort()
	require.NoError(t, err)

	subject, err := server.NewHTTPServer(server.Config{
		ShutdownTimeout: 5 * time.Second,
		Address:         fmt.Sprintf(":%d", port),
		AppConfig: &pkg.Config{
			DBUrl:            "file:ent?mode=memory&cache=shared&_fk=1",
			DebugPersistence: true,
		},
	}, l)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		defer func() {
			close(errCh)
			cancel()
		}()

		maxDuration := 3 * time.Second

		reqCtx, cancelReq := context.WithTimeout(ctx, maxDuration)
		defer cancelReq()

		req, err := http.NewRequestWithContext(reqCtx, "GET", fmt.Sprintf("http://localhost:%d/health", port), nil)
		require.NoError(t, err)

		bo := backoff.NewExponentialBackOff()
		bo.MaxElapsedTime = maxDuration

		errCh <- backoff.Retry(func() error {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("could not execute http request: %w", err)
			}

			if status := res.StatusCode; status != 200 {
				all, _ := io.ReadAll(res.Body)
				return fmt.Errorf("expected status: 200, got: %d, body: %s", status, string(all))
			}

			return nil
		}, bo)
	}()

	require.NoError(t, subject.Start(ctx))
	require.NoError(t, <-errCh)
}
