package profile

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type GetMyProfileQueryRepository interface { //nolint:iface
	service.OrganizationRepositoryGetOrganization
	service.UserRepositoryGetUser
	service.SpacemanagerGetPersonalSpaceInterface
}

type GetMyProfileQuery struct {
	repo   GetMyProfileQueryRepository
	logger *slog.Logger
}

func NewGetMyProfileQuery(repo GetMyProfileQueryRepository) *GetMyProfileQuery {
	return &GetMyProfileQuery{
		repo:   repo,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "GetMyProfileQuery")),
	}
}

func (u *GetMyProfileQuery) Execute(ctx context.Context, operator domain.UserInterface) (*domain.ProfileModel, error) {
	org, err := u.repo.GetOrganization(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("GetOrganization: %w", err)
	}
	user, err := u.repo.GetUser(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	var personalSpaceID *domain.SpaceID
	personalSpace, err := u.repo.GetPersonalSpace(ctx, operator)
	if err != nil {
		if errors.Is(err, service.ErrSpaceNotFound) {
			personalSpaceID = nil
		} else {
			return nil, fmt.Errorf("GetPersonalSpace: %w", err)
		}
	} else {
		personalSpaceID = personalSpace.SpaceID
	}

	return &domain.ProfileModel{
		LoginID:          user.LoginID,
		Username:         user.Username,
		OrganizationID:   org.OrganizationID,
		OrganizationName: org.Name,
		PersonalSpaceID:  personalSpaceID,
	}, nil
}
