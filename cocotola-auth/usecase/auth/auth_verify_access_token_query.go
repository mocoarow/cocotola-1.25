package auth

import (
	"context"
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type VerifyAccessTokenQueryGateway interface {
	service.UserRepositoryFindUserByLoginID
	service.UserRepositoryFindSystemOwnerByOrganizationName
	service.AuthTokenManagerGetUserInfo
}

type VerifyAccessTokenQuery struct {
	gw VerifyAccessTokenQueryGateway
}

func NewVerifyAccessTokenQuery(gw VerifyAccessTokenQueryGateway) *VerifyAccessTokenQuery {
	return &VerifyAccessTokenQuery{
		gw: gw,
	}
}

func (u *VerifyAccessTokenQuery) Execute(ctx context.Context, systemAdmin domain.SystemAdminInterface, bearerToken string) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "AuthVerifyAccessTokenQuery.Execute")
	defer span.End()

	// TODO: Check whether the token is registered in the Database
	userInfo, err := u.gw.GetUserInfo(ctx, bearerToken)
	if err != nil {
		return nil, fmt.Errorf("GetUserInfo: %w", err)
	}
	sysOwner, err := u.gw.FindSystemOwnerByOrganizationName(ctx, systemAdmin, userInfo.OrganizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}
	user, err := u.gw.FindUserByLoginID(ctx, sysOwner, userInfo.LoginID)
	if err != nil {
		return nil, fmt.Errorf("findUserbyLoginID: %w", err)
	}
	return user, nil
}
