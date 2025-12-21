package service

import (
	"context"
	"fmt"

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
