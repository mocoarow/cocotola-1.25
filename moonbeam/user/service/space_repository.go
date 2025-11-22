package service

import (
	"context"
	"errors"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
)

var ErrSpaceAlreadyExists = errors.New("space already exists")
var ErrSpaceNotFound = errors.New("space not found")

type CreateSpaceParameter struct {
	Key       string
	Name      string
	SpaceType string
}

type SpaceRepository interface {
	CreateSpace(ctx context.Context, operator domain.UserInterface, param *CreateSpaceParameter) (*domain.SpaceID, error)

	FindPublicSpaces(ctx context.Context, operator domain.UserInterface) ([]*domain.Space, error)

	FindPublicSpaceByKey(ctx context.Context, operator domain.UserInterface, key string) (*domain.Space, error)

	GetSpaceByID(ctx context.Context, operator domain.UserInterface, deckID *domain.SpaceID) (*domain.Space, error)
}
