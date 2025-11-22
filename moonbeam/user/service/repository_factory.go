package service

import (
	"context"

	mblibservice "github.com/mocoarow/cocotola-1.25/moonbeam/lib/service"
)

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) OrganizationRepository
	// NewUserRepository(ctx context.Context) UserRepository
	// NewUserGroupRepository(ctx context.Context) UserGroupRepository
	// NewSpaceRepository(ctx context.Context) SpaceRepository
	// NewSpaceManager(ctx context.Context) (SpaceManager, error)
}

type TransactionManager mblibservice.TransactionManagerT[RepositoryFactory]
