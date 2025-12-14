package process

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
)

type RunProcess func() error

type RunProcessFunc func(ctx context.Context) RunProcess

func WithAppServerProcess(router http.Handler, port int, readHeaderTimeout, shutdownTime time.Duration) RunProcessFunc {
	return func(ctx context.Context) RunProcess {
		return func() error {
			return libcontroller.AppServerProcess(ctx, router, port, readHeaderTimeout, shutdownTime)
		}
	}
}
func WithMetricsServerProcess(port int, shutdownTime int) RunProcessFunc {
	return func(ctx context.Context) RunProcess {
		return func() error {
			return MetricsServerProcess(ctx, port, shutdownTime)
		}
	}
}

func WithSignalWatchProcess() RunProcessFunc {
	return func(ctx context.Context) RunProcess {
		return func() error {
			return SignalWatchProcess(ctx)
		}
	}
}

func Run(ctx context.Context, runFuncs ...RunProcessFunc) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	errMu := &sync.Mutex{}
	var nonCanceledErr error

	for _, rf := range runFuncs {
		eg.Go(func() error {
			err := rf(ctx)()
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
