package gateway

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	mblibgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"

	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
)

type Process func() error

type ProcessFunc func(ctx context.Context) Process

func WithAppServerProcess(router http.Handler, port int, readHeaderTimeout, shutdownTime time.Duration) ProcessFunc {
	return func(ctx context.Context) Process {
		return func() error {
			return libcontroller.AppServerProcess(ctx, router, port, readHeaderTimeout, shutdownTime)
		}
	}
}
func WithMetricsServerProcess(port int, shutdownTime int) ProcessFunc {
	return func(ctx context.Context) Process {
		return func() error {
			return mblibgateway.MetricsServerProcess(ctx, port, shutdownTime)
		}
	}
}

func WithSignalWatchProcess() ProcessFunc {
	return func(ctx context.Context) Process {
		return func() error {
			return mblibgateway.SignalWatchProcess(ctx)
		}
	}
}

func Run(ctx context.Context, processFuncs ...ProcessFunc) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	errMu := &sync.Mutex{}
	var nonCanceledErr error

	for _, pf := range processFuncs {
		eg.Go(func() error {
			err := pf(ctx)()
			if err != nil && !errors.Is(err, context.Canceled) {
				errMu.Lock()
				if nonCanceledErr == nil {
					nonCanceledErr = err
				}
				errMu.Unlock()
			}

			return err
		})
	}

	if err := eg.Wait(); err != nil {
		if nonCanceledErr == nil && errors.Is(err, context.Canceled) {
			return 0
		}

		return 1
	}

	return 0
}
