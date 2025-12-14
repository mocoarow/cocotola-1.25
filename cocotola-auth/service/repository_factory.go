package service

import (
	"context"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"
)

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) OrganizationRepository
	NewUserRepository(ctx context.Context) UserRepository
	NewUserGroupRepository(ctx context.Context) UserGroupRepository
	NewSpaceRepository(ctx context.Context) SpaceRepository
	NewSpaceManager(ctx context.Context) (SpaceManager, error)

	NewAuthorizationManager(ctx context.Context) (AuthorizationManager, error)
}

type TransactionManager libservice.TransactionManagerT[RepositoryFactory]
