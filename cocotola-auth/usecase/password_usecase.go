package usecase

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"
)

type PasswordUsecae struct {
	systemToken      domain.SystemToken
	txManager        service.TransactionManager
	nonTxManager     service.TransactionManager
	authTokenManager service.AuthTokenManager
}

func NewPassword(systemToken domain.SystemToken, txManager, nonTxManager service.TransactionManager, authTokenManager service.AuthTokenManager) *PasswordUsecae {
	return &PasswordUsecae{
		systemToken:      systemToken,
		txManager:        txManager,
		nonTxManager:     nonTxManager,
		authTokenManager: authTokenManager,
	}
}

func (u *PasswordUsecae) Authenticate(ctx context.Context, loginID, password, organizationName string) (*service.AuthTokenSet, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.findSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	command := NewPasswordAuthenticateCommand(ctx, u.txManager, u.nonTxManager, u.authTokenManager)
	tokenSet, err := command.Execute(ctx, sysOwner, loginID, password)
	if err != nil {
		return nil, fmt.Errorf("command.Execute: %w", err)
	}
	return tokenSet, nil
}

func (u *PasswordUsecae) findSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.SystemOwner, error) {
	systemOwner, err := libservice.Do1(ctx, u.nonTxManager, func(rf service.RepositoryFactory) (*domain.SystemOwner, error) {
		return service.FindSystemOwnerByOrganizationName(ctx, rf, operator, organizationName)
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return systemOwner, nil
}
