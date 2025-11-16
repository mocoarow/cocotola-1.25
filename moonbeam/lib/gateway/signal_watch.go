package gateway

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

func SignalWatchProcess(ctx context.Context) error {
	logger := slog.Default().With(slog.String(domain.LoggerNameKey, "SignalWatch"))
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-sigs:
		logger.InfoContext(ctx, "shutdown signal received", slog.String("signal", sig.String()))
		return context.Canceled
	}
}
