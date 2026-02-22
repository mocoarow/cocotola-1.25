package profile

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type UsecaseRepository interface {
	service.OrganizationRepositoryGetOrganization
	service.UserRepositoryGetUser
	service.SpacemanagerGetPersonalSpaceInterface
}

type Usecase struct {
	repo   UsecaseRepository
	logger *slog.Logger
}

func NewUsecase(repo UsecaseRepository) *Usecase {
	return &Usecase{
		repo:   repo,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "Usecase")),
	}
}

func (u *Usecase) GetMyProfile(ctx context.Context, operator domain.UserInterface) (*domain.ProfileModel, error) {
	command := NewGetMyProfileQuery(u.repo)
	profile, err := command.Execute(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("GetMyProfileQuery.Execute: %w", err)
	}
	return profile, nil
}
