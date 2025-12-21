package service

import (
	"context"
	"errors"
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

var ErrSystemOwnerNotFound = errors.New("system owner not found")

var CreateUserAction = domain.NewRBACAction("CreateUser")
var ListUsersAction = domain.NewRBACAction("ListUsers")
var GetUserAction = domain.NewRBACAction("GetUser")
var UpdateUserAction = domain.NewRBACAction("UpdateUser")
var DeleteUserAction = domain.NewRBACAction("DeleteUser")

var CreateOwnerAction = domain.NewRBACAction("CreateOwner")

type CreateUserParameter struct {
	LoginID              string `validate:"required,max=255"`
	Username             string `validate:"required,max=255"`
	Password             string `validate:"required,min=8,max=255"`
	Provider             string
	ProviderLoginID      string
	ProviderAuthToken    string
	providerRefreshToken string
}

func NewCreateUserParameter(loginID, username, password, provider, providerLoginID, providerAuthToken, providerRefreshToken string) (*CreateUserParameter, error) {
	m := CreateUserParameter{
		LoginID:              loginID,
		Username:             username,
		Password:             password,
		Provider:             provider,
		ProviderLoginID:      providerLoginID,
		ProviderAuthToken:    providerAuthToken,
		providerRefreshToken: providerRefreshToken,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("new create user parameter: %w", err)
	}

	return &m, nil
}

type UserRepository interface {
	FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID) (*domain.SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.SystemOwner, error)

	GetUser(ctx context.Context, operator domain.UserInterface) (*domain.User, error)

	FindUserByID(ctx context.Context, operator domain.UserInterface, id *domain.UserID) (*domain.User, error)

	FindUserByLoginID(ctx context.Context, operator domain.UserInterface, loginID string) (*domain.User, error)

	FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwnerInterface, loginID string) (*domain.Owner, error)

	CreateUser(ctx context.Context, operator domain.UserInterface, param *CreateUserParameter) (*domain.UserID, error)

	CreateSystemOwner(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID) (*domain.UserID, error)

	VerifyPassword(ctx context.Context, operator domain.SystemOwnerInterface, loginID, password string) (bool, error)
}
