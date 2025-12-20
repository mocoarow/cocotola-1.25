package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-lib/process"
)

type ShutdownConfig struct {
	TimeSec1 int `yaml:"timeSec1" validate:"gte=1"`
	TimeSec2 int `yaml:"timeSec2" validate:"gte=1"`
}

func WithWebServerProcess(router http.Handler, port int, readHeaderTimeout, shutdownTime time.Duration) process.RunProcessFunc {
	return func(ctx context.Context) process.RunProcess {
		return func() error {
			return WebServerProcess(ctx, router, port, readHeaderTimeout, shutdownTime)
		}
	}
}
func WebServerProcess(ctx context.Context, router http.Handler, port int, readHeaderTimeout time.Duration, shutdownTime time.Duration) error {
	logger := slog.Default().With(slog.String(domain.LoggerNameKey, "WebServerProcess"))

	httpServer := http.Server{ //nolint:exhaustruct
		Addr:              ":" + strconv.Itoa(port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.InfoContext(ctx, fmt.Sprintf("http server listening at %v", httpServer.Addr))

	errCh := make(chan error)

	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.InfoContext(ctx, fmt.Sprintf("failed to ListenAndServe: %v", err))

			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTime)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.InfoContext(ctx, fmt.Sprintf("Server forced to shutdown: %v", err))

			return fmt.Errorf("httpServer.Shutdown: %w", err)
		}

		return nil
	case err := <-errCh:
		return fmt.Errorf("httpServer.ListenAndServe: %w", err)
	}
}
