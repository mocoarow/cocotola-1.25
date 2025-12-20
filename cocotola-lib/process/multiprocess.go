package process

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

type RunProcess func() error

type RunProcessFunc func(ctx context.Context) RunProcess

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
