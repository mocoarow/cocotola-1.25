package service

import (
	"context"

	mblibservice "github.com/mocoarow/cocotola-1.25/moonbeam/lib/service"
)

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) OrganizationRepository
}

type TransactionManager mblibservice.TransactionManagerT[RepositoryFactory]
