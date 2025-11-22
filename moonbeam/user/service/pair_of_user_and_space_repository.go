package service

import (
	"context"
	"errors"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
)

var ErrPairOfUserAndSpaceAlreadyExists = errors.New("pair of user and space already exists")
var ErrPairOfUserAndSpaceNotFound = errors.New("pair of user and space not found")

type PairOfUserAndSpaceRepository interface {
	CreatePairOfUserAndSpace(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, spaceID *domain.SpaceID) error

	FindMySpaces(ctx context.Context, operator domain.UserInterface) ([]*domain.Space, error)
}
