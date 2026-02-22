package auth

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type UsecaseGateway interface {
	service.UserRepositoryFindUserByLoginID
	service.UserRepositoryFindSystemOwnerByOrganizationName
	service.AuthTokenManagerGetUserInfo
	service.UserRepositoryVerifyPassword
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

func (u *Usecase) VerifyAccessToken(ctx context.Context, accessToken string) (*domain.User, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	query := NewVerifyAccessTokenQuery(u.gw)
	user, err := query.Execute(ctx, sysAdmin, accessToken)
	if err != nil {
		return nil, fmt.Errorf("query.Execute: %w", err)
	}
	return user, nil
}

func (u *Usecase) VerifyPassword(ctx context.Context, organizationName, loginID, password string) error {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.gw.FindSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	query := NewVerifyPasswordCommand(u.gw)
	if err := query.Execute(ctx, sysOwner, loginID, password); err != nil {
		return fmt.Errorf("query.Execute: %w", err)
	}
	return nil
}
