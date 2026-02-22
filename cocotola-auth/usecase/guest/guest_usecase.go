package guest

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type UsecaseGateway interface {
	service.UserRepositoryFindSystemOwnerByOrganizationName
	service.AuthTokenManagerCreateTokenSet
	service.UserRepositoryFindUserByLoginID
	service.OrganizationRepositoryGetOrganization
}

type Usecase struct {
	systemToken domain.SystemToken
	gw          UsecaseGateway
}

func NewUsecase(systemToken domain.SystemToken, gw UsecaseGateway) *Usecase {
	return &Usecase{
		systemToken: systemToken,
		gw:          gw,
	}
}

func (u *Usecase) Authenticate(ctx context.Context, organizationName string) (*service.AuthTokenSet, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.gw.FindSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	command := NewAuthenticateCommand(ctx, u.gw)
	tokenSet, err := command.Execute(ctx, sysOwner, organizationName)
	if err != nil {
		return nil, fmt.Errorf("command.Execute: %w", err)
	}
	return tokenSet, nil
}
