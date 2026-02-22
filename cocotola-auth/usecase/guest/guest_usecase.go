package guest

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type GuestUsecaseGateway interface {
	service.UserRepositoryFindSystemOwnerByOrganizationName
	service.AuthTokenManagerCreateTokenSet
	service.UserRepositoryFindUserByLoginID
	service.OrganizationRepositoryGetOrganization
}

type GuestUsecase struct {
	systemToken domain.SystemToken
	gw          GuestUsecaseGateway
}

func NewGuest(systemToken domain.SystemToken, gw GuestUsecaseGateway) *GuestUsecase {
	return &GuestUsecase{
		systemToken: systemToken,
		gw:          gw,
	}
}

func (u *GuestUsecase) Authenticate(ctx context.Context, organizationName string) (*service.AuthTokenSet, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.gw.FindSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	command := NewGuestAuthenticateCommand(ctx, u.gw)
	tokenSet, err := command.Execute(ctx, sysOwner, organizationName)
	if err != nil {
		return nil, fmt.Errorf("command.Execute: %w", err)
	}
	return tokenSet, nil
}
