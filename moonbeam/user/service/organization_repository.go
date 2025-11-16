package service

import (
	"context"
	"errors"
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
)

var ErrOrganizationNotFound = errors.New("organization not found")
var ErrOrganizationAlreadyExists = errors.New("organization already exists")

type CreateOrganizationParameter struct {
	Name       string            `validate:"required"`
	FirstOwner *AddUserParameter `validate:"required"`
}

func NewCreateOrganizationParameter(name string, firstOwner *AddUserParameter) (*CreateOrganizationParameter, error) {
	m := CreateOrganizationParameter{
		Name:       name,
		FirstOwner: firstOwner,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate CreateOrganizationParameter: %w", err)
	}

	return &m, nil
}

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, operator domain.UserInterface) (*domain.Organization, error)

	FindOrganizationByName(ctx context.Context, operator domain.SystemAdminInterface, name string) (*domain.Organization, error)

	FindOrganizationByID(ctx context.Context, operator domain.SystemAdminInterface, id *domain.OrganizationID) (*domain.Organization, error)

	CreateOrganization(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.OrganizationID, error)
}
