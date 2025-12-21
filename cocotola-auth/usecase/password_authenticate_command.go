package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type PasswordAuthenticateCommand struct {
	txManager        service.TransactionManager
	nonTxManager     service.TransactionManager
	authTokenManager service.AuthTokenManager
}

func NewPasswordAuthenticateCommand(_ context.Context, txManager, nonTxManager service.TransactionManager, authTokenManager service.AuthTokenManager) *PasswordAuthenticateCommand {
	return &PasswordAuthenticateCommand{
		txManager:        txManager,
		nonTxManager:     nonTxManager,
		authTokenManager: authTokenManager,
	}
}

func (u *PasswordAuthenticateCommand) Execute(ctx context.Context, systemOwner domain.SystemOwnerInterface, loginID, password string) (*service.AuthTokenSet, error) {
	// 1. Check authorization
	if err := u.checkAuthorization(ctx, systemOwner, loginID); err != nil {
		return nil, fmt.Errorf("checkAuthorization: %w", err)
	}

	// 2. Execute
	tokenSet, err := u.execute(ctx, systemOwner, loginID, password)
	if err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}

	// 3. Callback
	if err := u.callback(ctx, systemOwner); err != nil {
		return nil, fmt.Errorf("callback: %w", err)
	}
	return tokenSet, nil
}

func (u *PasswordAuthenticateCommand) checkAuthorization(_ context.Context, _ domain.SystemOwnerInterface, loginID string) error {
	if strings.Contains(loginID, "guest@@") {
		return service.ErrUnauthenticated
	}

	return nil
}

func (u *PasswordAuthenticateCommand) execute(ctx context.Context, systemOwner domain.SystemOwnerInterface, loginID, password string) (*service.AuthTokenSet, error) {
	fn := func(mbrf service.RepositoryFactory) error {
		userRepo := mbrf.NewUserRepository(ctx)
		ok, err := userRepo.VerifyPassword(ctx, systemOwner, loginID, password)
		if err != nil {
			return fmt.Errorf("action.userRepo.VerifyPassword: %w", err)
		} else if !ok {
			return service.ErrUnauthenticated
		}
		return nil
	}
	if err := libservice.Do0(ctx, u.nonTxManager, fn); err != nil {
		return nil, err //nolint:wrapcheck
	}

	user, err := u.findUserbyLoginID(ctx, systemOwner, loginID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, service.ErrUnauthenticated
		}
		return nil, fmt.Errorf("findUserbyLoginID: %w", err)
	}
	org, err := u.getOrganization(ctx, systemOwner)
	if err != nil {
		return nil, fmt.Errorf("getOrganization: %w", err)
	}
	tokenSet, err := u.createTokenSet(ctx, user, org)
	if err != nil {
		return nil, fmt.Errorf("create token set: %w", err)
	}
	return tokenSet, nil
}

func (u *PasswordAuthenticateCommand) callback(_ context.Context, _ domain.SystemOwnerInterface) error {
	return nil
}

func (u *PasswordAuthenticateCommand) findUserbyLoginID(ctx context.Context, operator domain.UserInterface, loginID string) (*domain.User, error) {
	return libservice.Do1(ctx, u.nonTxManager, func(mbrf service.RepositoryFactory) (*domain.User, error) { //nolint:wrapcheck
		return findUserbyLoginID(ctx, mbrf, operator, loginID)
	})
}

func (u *PasswordAuthenticateCommand) getOrganization(ctx context.Context, operator domain.UserInterface) (*domain.Organization, error) {
	return libservice.Do1(ctx, u.nonTxManager, func(mbrf service.RepositoryFactory) (*domain.Organization, error) { //nolint:wrapcheck
		return getOrganization(ctx, mbrf, operator)
	})
}

func (u *PasswordAuthenticateCommand) createTokenSet(ctx context.Context, user *domain.User, organization *domain.Organization) (*service.AuthTokenSet, error) {
	return createTokenSet(ctx, u.authTokenManager, user, organization)
}
