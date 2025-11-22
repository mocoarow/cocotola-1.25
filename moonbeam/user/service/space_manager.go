package service

import (
	"context"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
)

type CreatePersonalSpaceParameter struct {
	UserID  *domain.UserID
	KeyName string
	Name    string
}
type SpaceManager interface {
	CreatePersonalSpace(ctx context.Context, operator domain.UserInterface, param *CreatePersonalSpaceParameter) (*domain.SpaceID, error)
	AddUserToSpace(ctx context.Context, operator domain.SystemOwnerInterface, userID domain.UserID, spaceID *domain.SpaceID) error
	GetPersonalSpace(ctx context.Context, operator domain.UserInterface) (*domain.Space, error)
}
