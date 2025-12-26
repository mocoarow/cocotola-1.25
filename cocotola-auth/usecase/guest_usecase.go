package usecase

import (
	"context"
	"fmt"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type GuestUsecase struct {
	systemToken      domain.SystemToken
	mbTxManager      service.TransactionManager
	mbNonTxManager   service.TransactionManager
	authTokenManager service.AuthTokenManager
}

func NewGuest(systemToken domain.SystemToken, mbTxManager, mbNonTxManager service.TransactionManager, authTokenManager service.AuthTokenManager) *GuestUsecase {
	return &GuestUsecase{
		systemToken:      systemToken,
		mbTxManager:      mbTxManager,
		mbNonTxManager:   mbNonTxManager,
		authTokenManager: authTokenManager,
	}
}

func (u *GuestUsecase) Authenticate(ctx context.Context, organizationName string) (*service.AuthTokenSet, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.findSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	command := NewGuestAuthenticateCommand(ctx, u.mbTxManager, u.mbNonTxManager, u.authTokenManager)
	tokenSet, err := command.Execute(ctx, sysOwner, organizationName)
	if err != nil {
		return nil, fmt.Errorf("command.Execute: %w", err)
	}
	return tokenSet, nil
}

func (u *GuestUsecase) findSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.SystemOwner, error) {
	return libservice.Do1(ctx, u.mbNonTxManager, func(mbrf service.RepositoryFactory) (*domain.SystemOwner, error) { //nolint:wrapcheck
		return service.FindSystemOwnerByOrganizationName(ctx, mbrf, operator, organizationName)
	})
}
