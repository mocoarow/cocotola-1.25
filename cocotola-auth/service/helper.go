package service

import (
	"context"
	"fmt"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

func FindSystemOwnerByOrganizationName(ctx context.Context, rf RepositoryFactory, systemAdmin domain.SystemAdminInterface, organizationName string) (*domain.SystemOwner, error) {
	userRepo := rf.NewUserRepository(ctx)
	sysOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("FindSystemOwnerByOrganizationName: %w", err)
	}
	return sysOwner, nil
}

func FindPublicSpaceByKey(ctx context.Context, systemOwner domain.SystemOwnerInterface, nonTxManager TransactionManager, key string) (*domain.Space, error) {
	fn := func(rf RepositoryFactory) (*domain.Space, error) {
		spaceRepo := rf.NewSpaceRepository(ctx)
		publicDefaultSpace, err := spaceRepo.FindPublicSpaceByKey(ctx, systemOwner, key)
		if err != nil {
			return nil, fmt.Errorf("find public default space by key(%s): %w", key, err)
		}

		return publicDefaultSpace, nil
	}
	publicDefaultSpace, err := libservice.Do1(ctx, nonTxManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return publicDefaultSpace, nil
}
