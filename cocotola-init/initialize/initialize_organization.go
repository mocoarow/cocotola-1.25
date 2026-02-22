package initialize

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	"gorm.io/gorm"

	"github.com/mocoarow/cocotola-1.25/cocotola-init/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-init/usecase"
)

func initOrganization(ctx context.Context, systemToken authdomain.SystemToken, dbc *libgateway.DBConnection, organizationName, loginID, password string) error {
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, domain.AppName+"InitApp1"))

	sysAdmin := authdomain.NewSystemAdmin(systemToken)

	// 1. check whether the organization already exists
	{
		found, err := findOrganizationAndSystemOwnerAndPublicDefaultSpace(ctx, sysAdmin, dbc, organizationName)
		if err != nil {
			return fmt.Errorf("findOrganizationAndPublicDefaultSpace: %w", err)
		}
		if found {
			return nil
		}
	}

	// 2. create organization
	orgID2, err := createOrganization(ctx, sysAdmin, dbc, organizationName)
	if err != nil {
		return fmt.Errorf("create organization: %w", err)
	}
	logger.InfoContext(ctx, fmt.Sprintf("organizationID: %d", orgID2.Int()))

	// 3. find system owner
	userRepo := authgateway.NewUserRepository(dbc)
	sysOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, sysAdmin, organizationName)
	if err != nil {
		return fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	// 4. create first owner
	firstOwnerID, err := createFirstOwnerToOrganization(ctx, sysOwner, dbc, loginID, password)
	if err != nil {
		return fmt.Errorf("create first owner: %w", err)
	}
	logger.InfoContext(ctx, fmt.Sprintf("firstOwnerID: %d", firstOwnerID.Int()))

	// 5. find public default space
	spaceRepo := authgateway.NewSpaceRepository(dbc)
	if _, err := spaceRepo.FindPublicSpaceByKey(ctx, sysOwner, authservice.PublicDefaultSpaceKey); err != nil {
		return fmt.Errorf("find public default space by key(%s): %w", authservice.PublicDefaultSpaceKey, err)
	}

	return nil
}

func findOrganizationAndSystemOwnerAndPublicDefaultSpace(ctx context.Context, systemAdmin authdomain.SystemAdminInterface, dbc *libgateway.DBConnection, organizationName string) (bool, error) {
	orgRepo := authgateway.NewOrganizationRepository(dbc)
	if _, err := orgRepo.FindOrganizationByName(ctx, systemAdmin, organizationName); err != nil {
		if !errors.Is(err, authservice.ErrOrganizationNotFound) {
			return false, fmt.Errorf("find organization by name: %w", err)
		}
		return false, nil
	}

	userRepo := authgateway.NewUserRepository(dbc)
	sysOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, organizationName)
	if err != nil {
		return false, fmt.Errorf("find system owner by organization name: %w", err)
	}

	spaceRepo := authgateway.NewSpaceRepository(dbc)
	if _, err := spaceRepo.FindPublicSpaceByKey(ctx, sysOwner, authservice.PublicDefaultSpaceKey); err != nil {
		if !errors.Is(err, authservice.ErrSpaceNotFound) {
			return false, fmt.Errorf("find public default space by key: %w", err)
		}
		return false, nil
	}

	return true, nil
}

type CreateOrganizationCommandGateway struct {
	dbc *libgateway.DBConnection
}

// NewCreateOrganizationCommandGateway returns a new transaction manager.
func NewCreateOrganizationCommandGateway(dbc *libgateway.DBConnection) *CreateOrganizationCommandGateway {
	return &CreateOrganizationCommandGateway{
		dbc: dbc,
	}
}

func (gw *CreateOrganizationCommandGateway) WithTransaction(ctx context.Context, fn func(
	createOrganization authservice.CreateOrganizationFunc,
	createSystemOwner authservice.CreateSystemOwnerFunc,
	findSystemOwnerByOrganizationName authservice.FindSystemOwnerByOrganizationNameFunc,
	attachPolicyToUserBySystemAdmin authservice.AttachPolicyToUserBySystemAdminFunc,
	createOwnerGroup authservice.CreateOwnerGroupFunc,
	attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc,
	createPublicGroup authservice.CreatePublicGroupFunc,
	createPublicDefaultSpace authservice.CreatePublicDefaultSpaceFunc,
) (*authdomain.OrganizationID, error)) (*authdomain.OrganizationID, error) {
	var organizationID *authdomain.OrganizationID
	err := gw.dbc.DB.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
		orgRepo := authgateway.NewOrganizationRepository(gw.dbc)
		userRepo := authgateway.NewUserRepository(gw.dbc)
		groupRepo := authgateway.NewUserGroupRepository(gw.dbc)
		authManager, err := authgateway.NewAuthorizationManager(ctx, gw.dbc)
		if err != nil {
			return fmt.Errorf("new AuthorizationManager: %w", err)
		}
		spaceManager, err := authgateway.NewSpaceManager(ctx, gw.dbc)
		if err != nil {
			return fmt.Errorf("new SpaceManager: %w", err)
		}

		organizationID, err = fn(
			orgRepo.CreateOrganization,
			userRepo.CreateSystemOwner,
			userRepo.FindSystemOwnerByOrganizationName,
			authManager.AttachPolicyToUserBySystemAdmin,
			groupRepo.CreateOwnerGroup,
			authManager.AttachPolicyToUserBySystemOwner,
			groupRepo.CreatePublicGroup,
			spaceManager.CreatePublicDefaultSpace)
		if err != nil {
			return fmt.Errorf("execute function in transaction: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return organizationID, nil
}

func createOrganization(ctx context.Context, operator authdomain.SystemAdminInterface, dbc *libgateway.DBConnection, organizationName string) (*authdomain.OrganizationID, error) {
	createOrganizationCommandGateway := NewCreateOrganizationCommandGateway(dbc)
	command := usecase.NewCreateOrganizationCommand(ctx, createOrganizationCommandGateway)
	organizationID, err := command.Execute(ctx, operator, organizationName)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}
	return organizationID, nil
}

type CreateFirstOwnerCommandGateway struct {
	dbc *libgateway.DBConnection
}

func NewCreateFirstOwnerCommandGateway(dbc *libgateway.DBConnection) *CreateFirstOwnerCommandGateway {
	return &CreateFirstOwnerCommandGateway{
		dbc: dbc,
	}
}

func (gw *CreateFirstOwnerCommandGateway) WithTransaction(ctx context.Context, fn func(
	createUser authservice.CreateUserFunc,
	findUserByID authservice.FindUserByIDFunc,
	findUserGroupByKey authservice.FindUserGroupByKeyFunc,
	addUserToGroup authservice.AddUserToGroupFunc,
	attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc,
	createPersonalSpace authservice.CreatePersonalSpaceFunc,
) (*authdomain.UserID, error)) (*authdomain.UserID, error) {
	var userID *authdomain.UserID
	err := gw.dbc.DB.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
		userRepo := authgateway.NewUserRepository(gw.dbc)
		groupRepo := authgateway.NewUserGroupRepository(gw.dbc)
		authManager, err := authgateway.NewAuthorizationManager(ctx, gw.dbc)
		if err != nil {
			return fmt.Errorf("new AuthorizationManager: %w", err)
		}
		spaceManager, err := authgateway.NewSpaceManager(ctx, gw.dbc)
		if err != nil {
			return fmt.Errorf("new SpaceManager: %w", err)
		}

		userID, err = fn(
			userRepo.CreateUser,
			userRepo.FindUserByID,
			groupRepo.FindUserGroupByKey,
			authManager.AddUserToGroup,
			authManager.AttachPolicyToUserBySystemOwner,
			spaceManager.CreatePersonalSpace)
		if err != nil {
			return fmt.Errorf("execute function in transaction: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return userID, nil
}

func createFirstOwnerToOrganization(ctx context.Context, operator authdomain.SystemOwnerInterface, dbc *libgateway.DBConnection, loginID, password string) (*authdomain.UserID, error) {
	firstOwnerAddParam, err := authservice.NewCreateUserParameter(loginID, "Owner(cocotola)", password, "", "", "", "")
	if err != nil {
		return nil, fmt.Errorf("new UserAddParameter: %w", err)
	}
	createFirstOwnerCommandGateway := NewCreateFirstOwnerCommandGateway(dbc)
	createFirstOwnerCommand := usecase.NewCreateFirstOwnerCommand(ctx, createFirstOwnerCommandGateway)
	firstOwnerID, err := createFirstOwnerCommand.Execute(ctx, operator, firstOwnerAddParam)
	if err != nil {
		return nil, fmt.Errorf("create first owner: %w", err)
	}

	return firstOwnerID, nil
}
