package usecase

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

// func findUserByID(ctx context.Context, mbrf service.RepositoryFactory, operator domain.UserInterface, userID *domain.UserID) (*domain.User, error) {
// 	userRepo := mbrf.NewUserRepository(ctx)
// 	user, err := userRepo.FindUserByID(ctx, operator, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("find user by id(%d): %w", userID.Int(), err)
// 	}
// 	return user, nil
// }

func findUserbyLoginID(ctx context.Context, mbrf service.RepositoryFactory, operator domain.UserInterface, loginID string) (*domain.User, error) {
	userRepo := mbrf.NewUserRepository(ctx)
	user, err := userRepo.FindUserByLoginID(ctx, operator, loginID)
	if err != nil {
		return nil, fmt.Errorf("find user by login id(%s): %w", loginID, err)
	}
	return user, nil
}

func getOrganization(ctx context.Context, mbrf service.RepositoryFactory, operator domain.UserInterface) (*domain.Organization, error) {
	orgRepo := mbrf.NewOrganizationRepository(ctx)
	org, err := orgRepo.GetOrganization(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	return org, nil
}

func createTokenSet(ctx context.Context, authTokenManager service.AuthTokenManager, user *domain.User, organization *domain.Organization) (*service.AuthTokenSet, error) {
	tokenSet, err := authTokenManager.CreateTokenSet(ctx, user, organization.OrganizationID, organization.Name)
	if err != nil {
		return nil, fmt.Errorf("create token set: %w", err)
	}
	return tokenSet, nil
}
