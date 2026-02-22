package profile

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type ProfileUsecaseRepository interface {
	service.OrganizationRepositoryGetOrganization
	service.UserRepositoryGetUser
	service.SpacemanagerGetPersonalSpaceInterface
}

type ProfileUsecase struct {
	repo   ProfileUsecaseRepository
	logger *slog.Logger
}

func NewProfileUsecase(repo ProfileUsecaseRepository) *ProfileUsecase {
	return &ProfileUsecase{
		repo:   repo,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "ProfileUsecase")),
	}
}

func (u *ProfileUsecase) GetMyProfile(ctx context.Context, operator domain.UserInterface) (*domain.ProfileModel, error) {
	command := NewGetMyProfileQuery(u.repo)
	profile, err := command.Execute(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("GetMyProfileQuery.Execute: %w", err)
	}
	return profile, nil
}
