package guest

import (
	"context"
	"errors"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type GuestAuthenticateCommandGateway interface {
	service.AuthTokenManagerCreateTokenSet
	service.UserRepositoryFindUserByLoginID
	service.OrganizationRepositoryGetOrganization
}

type GuestAuthenticateCommand struct {
	gw GuestAuthenticateCommandGateway
}

func NewGuestAuthenticateCommand(_ context.Context, gw GuestAuthenticateCommandGateway) *GuestAuthenticateCommand {
	return &GuestAuthenticateCommand{
		gw: gw,
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
	user, err := u.gw.FindUserByLoginID(ctx, systemOwner, guestLoginID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, service.ErrUnauthenticated
		}
		return nil, fmt.Errorf("findUserbyLoginID: %w", err)
	}
	org, err := u.gw.GetOrganization(ctx, systemOwner)
	if err != nil {
		return nil, fmt.Errorf("getOrganization: %w", err)
	}
	tokenSet, err := u.gw.CreateTokenSet(ctx, user, org.OrganizationID, org.Name)
	if err != nil {
		return nil, fmt.Errorf("create token set: %w", err)
	}
	return tokenSet, nil
}

func (u *GuestAuthenticateCommand) callback(_ context.Context, _ domain.SystemOwnerInterface) error {
	return nil
}
