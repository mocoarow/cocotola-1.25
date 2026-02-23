package service

import (
	"context"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

type CreatePersonalSpaceParameter struct {
	UserID  *domain.UserID
	KeyName string
	Name    string
}

type SpacemanagerGetPersonalSpaceInterface interface {
	GetPersonalSpace(ctx context.Context, operator domain.UserInterface) (*domain.Space, error)
}

type CreatePublicDefaultSpaceFunc func(ctx context.Context, operator domain.SystemOwnerInterface) (*domain.SpaceID, error)

type CreatePersonalSpaceFunc func(ctx context.Context, operator domain.UserInterface, param *CreatePersonalSpaceParameter) (*domain.SpaceID, error)
