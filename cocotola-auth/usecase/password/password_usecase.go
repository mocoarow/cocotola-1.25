package password

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type PasswordUsecaseGateway interface {
	service.UserRepositoryFindSystemOwnerByOrganizationName
	service.UserRepositoryVerifyPassword
	service.UserRepositoryFindUserByLoginID
	service.OrganizationRepositoryGetOrganization
	service.AuthTokenManagerCreateTokenSet
}

type PasswordUsecase struct {
	systemToken domain.SystemToken
	gw          PasswordUsecaseGateway
}

func NewPassword(systemToken domain.SystemToken, gw PasswordUsecaseGateway) *PasswordUsecase {
	return &PasswordUsecase{
		systemToken: systemToken,
		gw:          gw,
	}
}

func (u *PasswordUsecase) Authenticate(ctx context.Context, loginID, password, organizationName string) (*service.AuthTokenSet, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.gw.FindSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	command := NewPasswordAuthenticateCommand(ctx, u.gw)
	tokenSet, err := command.Execute(ctx, sysOwner, loginID, password)
	if err != nil {
		return nil, fmt.Errorf("command.Execute: %w", err)
	}
	return tokenSet, nil
}
