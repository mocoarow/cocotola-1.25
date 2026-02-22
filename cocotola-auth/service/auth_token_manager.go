package service

import (
	"context"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

type AuthTokenSet struct {
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	LoginID          string
	Username         string
	OrganizationID   int
	OrganizationName string
}

type CreateTokenSetFunc func(ctx context.Context, user domain.UserInterface, organizationID *domain.OrganizationID, organizationName string) (*AuthTokenSet, error)

type AuthTokenManagerCreateTokenSet interface {
	CreateTokenSet(ctx context.Context, user domain.UserInterface, organizationID *domain.OrganizationID, organizationName string) (*AuthTokenSet, error)
}
type AuthTokenManagerGetUserInfo interface {
	GetUserInfo(ctx context.Context, tokenString string) (*UserInfo, error)
}
type AuthTokenManager interface {
	// SignInWithIDToken(ctx context.Context, idToken string) (*AuthTokenSet, error)
	GetUserInfo(ctx context.Context, tokenString string) (*UserInfo, error)

	CreateTokenSet(ctx context.Context, user domain.UserInterface, organizationID *domain.OrganizationID, organizationName string) (*AuthTokenSet, error)
	// RefreshToken(ctx context.Context, accessToken string) (string, error)
}
