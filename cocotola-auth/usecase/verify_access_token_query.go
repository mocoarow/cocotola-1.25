package usecase

import (
	"context"
	"fmt"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type VerifyAccessTokenQuery struct {
	mbNonTxManager   service.TransactionManager
	authTokenManager service.AuthTokenManager
}

func NewVerifyAccessTokenQuery(mbNonTxManager service.TransactionManager, authTokenManager service.AuthTokenManager) *VerifyAccessTokenQuery {
	return &VerifyAccessTokenQuery{
		mbNonTxManager:   mbNonTxManager,
		authTokenManager: authTokenManager,
	}
}

func (u *VerifyAccessTokenQuery) Execute(ctx context.Context, systemAdmin domain.SystemAdminInterface, bearerToken string) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "VerifyAccessTokenQuery.Execute")
	defer span.End()

	// TODO: Check whether the token is registered in the Database
	userInfo, err := u.authTokenManager.GetUserInfo(ctx, bearerToken)
	if err != nil {
		return nil, fmt.Errorf("GetUserInfo: %w", err)
	}
	sysOwner, err := u.findSystemOwnerByOrganizationName(ctx, systemAdmin, userInfo.OrganizationName)
	if err != nil {
		return nil, fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}
	user, err := u.findUserbyLoginID(ctx, sysOwner, userInfo.LoginID)
	if err != nil {
		return nil, fmt.Errorf("findUserbyLoginID: %w", err)
	}
	return user, nil
}

func (u *VerifyAccessTokenQuery) findSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.SystemOwner, error) {
	return libservice.Do1(ctx, u.mbNonTxManager, func(mbrf service.RepositoryFactory) (*domain.SystemOwner, error) { //nolint:wrapcheck
		return service.FindSystemOwnerByOrganizationName(ctx, mbrf, operator, organizationName)
	})
}

func (u *VerifyAccessTokenQuery) findUserbyLoginID(ctx context.Context, operator domain.UserInterface, loginID string) (*domain.User, error) {
	return libservice.Do1(ctx, u.mbNonTxManager, func(mbrf service.RepositoryFactory) (*domain.User, error) { //nolint:wrapcheck
		return findUserbyLoginID(ctx, mbrf, operator, loginID)
	})
}
