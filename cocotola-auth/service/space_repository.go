package service

import (
	"context"
	"errors"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

var ErrSpaceAlreadyExists = errors.New("space already exists")
var ErrSpaceNotFound = errors.New("space not found")

type CreateSpaceParameter struct {
	Key       string
	Name      string
	SpaceType string
}

type FindPublicSpaceByKeyFunc func(ctx context.Context, operator domain.UserInterface, key string) (*domain.Space, error)
