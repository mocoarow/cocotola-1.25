package usecase

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type GetMyProfileQuery struct {
	nonTxManager service.TransactionManager
	logger       *slog.Logger
}

func NewGetMyProfileQuery(nonTxManager service.TransactionManager) *GetMyProfileQuery {
	return &GetMyProfileQuery{
		nonTxManager: nonTxManager,
		logger:       slog.Default().With(slog.String(libdomain.LoggerNameKey, "GetMyProfileQuery")),
	}
}

func (u *GetMyProfileQuery) Execute(ctx context.Context, operator domain.UserInterface) (*domain.ProfileModel, error) {
	fn := func(rf service.RepositoryFactory) (*domain.ProfileModel, error) {
		orgRepo := rf.NewOrganizationRepository(ctx)
		org, err := orgRepo.GetOrganization(ctx, operator)
		if err != nil {
			return nil, fmt.Errorf("GetOrganization: %w", err)
		}
		userRepo := rf.NewUserRepository(ctx)
		user, err := userRepo.GetUser(ctx, operator)
		if err != nil {
			return nil, fmt.Errorf("GetUser: %w", err)
		}
		spaceManager, err := rf.NewSpaceManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("NewSpaceManager: %w", err)
		}
		privateSpace, err := spaceManager.GetPersonalSpace(ctx, operator)
		if err != nil {
			return nil, fmt.Errorf("GetPersonalSpace: %w", err)
		}

		return &domain.ProfileModel{
			LoginID:          user.LoginID,
			Username:         user.Username,
			OrganizationID:   org.OrganizationID,
			OrganizationName: org.Name,
			PrivateSpaceID:   privateSpace.SpaceID,
		}, nil
	}
	profileModel, err := libservice.Do1(ctx, u.nonTxManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return profileModel, nil
}
