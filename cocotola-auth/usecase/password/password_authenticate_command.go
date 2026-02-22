package password

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type AuthenticateCommandRepository interface {
	service.UserRepositoryVerifyPassword
	service.UserRepositoryFindUserByLoginID
	service.OrganizationRepositoryGetOrganization
	service.AuthTokenManagerCreateTokenSet
}

type AuthenticateCommand struct {
	repo AuthenticateCommandRepository
}

func NewAuthenticateCommand(_ context.Context, repo AuthenticateCommandRepository) *AuthenticateCommand {
	return &AuthenticateCommand{
		repo: repo,
	}
}

func (u *AuthenticateCommand) Execute(ctx context.Context, systemOwner domain.SystemOwnerInterface, loginID, password string) (*service.AuthTokenSet, error) {
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

func (u *AuthenticateCommand) checkAuthorization(_ context.Context, _ domain.SystemOwnerInterface, loginID string) error {
	if strings.Contains(loginID, "guest@@") {
		return service.ErrUnauthenticated
	}

	return nil
}

func (u *AuthenticateCommand) execute(ctx context.Context, systemOwner domain.SystemOwnerInterface, loginID, password string) (*service.AuthTokenSet, error) {
	if ok, err := u.repo.VerifyPassword(ctx, systemOwner, loginID, password); err != nil {
		return nil, fmt.Errorf("verifyPassword: %w", err)
	} else if !ok {
		return nil, service.ErrUnauthenticated
	}

	user, err := u.repo.FindUserByLoginID(ctx, systemOwner, loginID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, service.ErrUnauthenticated
		}
		return nil, fmt.Errorf("findUserbyLoginID: %w", err)
	}
	org, err := u.repo.GetOrganization(ctx, systemOwner)
	if err != nil {
		return nil, fmt.Errorf("getOrganization: %w", err)
	}
	tokenSet, err := u.repo.CreateTokenSet(ctx, user, org.OrganizationID, org.Name)
	if err != nil {
		return nil, fmt.Errorf("create token set: %w", err)
	}
	return tokenSet, nil
}

func (u *AuthenticateCommand) callback(_ context.Context, _ domain.SystemOwnerInterface) error {
	return nil
}
