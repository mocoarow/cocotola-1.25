package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"
)

type GuestAuthenticateCommand struct {
	mbTxManager      service.TransactionManager
	mbNonTxManager   service.TransactionManager
	authTokenManager service.AuthTokenManager
}

func NewGuestAuthenticateCommand(_ context.Context, mbTxManager, mbNonTxManager service.TransactionManager, authTokenManager service.AuthTokenManager) *GuestAuthenticateCommand {
	return &GuestAuthenticateCommand{
		mbTxManager:      mbTxManager,
		mbNonTxManager:   mbNonTxManager,
		authTokenManager: authTokenManager,
	}
}

func (u *GuestAuthenticateCommand) Execute(ctx context.Context, systemOwner domain.SystemOwnerInterface, organizationName string) (*service.AuthTokenSet, error) {
	// 1. Check authorization
	if err := u.checkAuthorization(ctx, systemOwner); err != nil {
		return nil, fmt.Errorf("checkAuthorization: %w", err)
	}

	// 2. Execute
	tokenSet, err := u.execute(ctx, systemOwner, organizationName)
	if err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}

	// 3. Callback
	if err := u.callback(ctx, systemOwner); err != nil {
		return nil, fmt.Errorf("callback: %w", err)
	}
	return tokenSet, nil
}

func (u *GuestAuthenticateCommand) checkAuthorization(_ context.Context, _ domain.SystemOwnerInterface) error {
	return nil
}

func (u *GuestAuthenticateCommand) execute(ctx context.Context, systemOwner domain.SystemOwnerInterface, organizationName string) (*service.AuthTokenSet, error) {
	guestLoginID := domain.NewGuestLoginID(organizationName)
	user, err := u.findUserbyLoginID(ctx, systemOwner, guestLoginID)
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

func (u *GuestAuthenticateCommand) callback(_ context.Context, _ domain.SystemOwnerInterface) error {
	return nil
}

func (u *GuestAuthenticateCommand) findUserbyLoginID(ctx context.Context, operator domain.UserInterface, loginID string) (*domain.User, error) {
	return libservice.Do1(ctx, u.mbNonTxManager, func(mbrf service.RepositoryFactory) (*domain.User, error) { //nolint:wrapcheck
		return findUserbyLoginID(ctx, mbrf, operator, loginID)
	})
}

func (u *GuestAuthenticateCommand) getOrganization(ctx context.Context, operator domain.UserInterface) (*domain.Organization, error) {
	return libservice.Do1(ctx, u.mbNonTxManager, func(mbrf service.RepositoryFactory) (*domain.Organization, error) { //nolint:wrapcheck
		return getOrganization(ctx, mbrf, operator)
	})
}

func (u *GuestAuthenticateCommand) createTokenSet(ctx context.Context, user *domain.User, organization *domain.Organization) (*service.AuthTokenSet, error) {
	return createTokenSet(ctx, u.authTokenManager, user, organization)
}
