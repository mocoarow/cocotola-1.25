package auth

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type AuthUsecaseGateway interface {
	service.UserRepositoryFindUserByLoginID
	service.UserRepositoryFindSystemOwnerByOrganizationName
	service.AuthTokenManagerGetUserInfo
	service.UserRepositoryVerifyPassword
}

type AuthUsecase struct {
	systemToken domain.SystemToken
	gw          AuthUsecaseGateway
}

func NewAuthUsecase(systemToken domain.SystemToken, gw AuthUsecaseGateway) *AuthUsecase {
	return &AuthUsecase{
		systemToken: systemToken,
		gw:          gw,
	}
}

func (u *AuthUsecase) VerifyAccessToken(ctx context.Context, accessToken string) (*domain.User, error) {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	query := NewAuthVerifyAccessTokenQuery(u.gw)
	user, err := query.Execute(ctx, sysAdmin, accessToken)
	if err != nil {
		return nil, fmt.Errorf("query.Execute: %w", err)
	}
	return user, nil
}

func (u *AuthUsecase) VerifyPassword(ctx context.Context, organizationName, loginID, password string) error {
	sysAdmin := domain.NewSystemAdmin(u.systemToken)
	sysOwner, err := u.gw.FindSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	query := NewAuthVerifyPasswordCommand(u.gw)
	if err := query.Execute(ctx, sysOwner, loginID, password); err != nil {
		return fmt.Errorf("query.Execute: %w", err)
	}
	return nil
}
