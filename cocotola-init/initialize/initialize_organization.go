package initialize

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	usecase "github.com/mocoarow/cocotola-1.25/cocotola-init/usecase"
)

func initOrganization(ctx context.Context, systemToken authdomain.SystemToken, mbTxManager, mbNonTxManager authservice.TransactionManager, organizationName, loginID, password, appName string) error {
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, appName+"InitApp1"))

	sysAdmin := authdomain.NewSystemAdmin(systemToken)

	// 1. check whether the organization already exists
	{
		found, err := findOrganizationAndSystemOwnerAndPublicDefaultSpace(ctx, sysAdmin, mbNonTxManager, organizationName)
		if err != nil {
			return fmt.Errorf("findOrganizationAndPublicDefaultSpace: %w", err)
		}
		if found {
			return nil
		}
	}

	// 2. create organization
	orgID2, err := createOrganization(ctx, sysAdmin, mbTxManager, mbNonTxManager, organizationName)
	if err != nil {
		return fmt.Errorf("create organization: %w", err)
	}
	logger.InfoContext(ctx, fmt.Sprintf("organizationID: %d", orgID2.Int()))

	// 3. find system owner
	sysOwner, err := findSystemOwnerByOrganizationName(ctx, sysAdmin, mbNonTxManager, organizationName)
	if err != nil {
		return fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	// 4. create first owner
	firstOwnerID, err := createFirstOwnerToOrganization(ctx, sysOwner, mbTxManager, mbNonTxManager, loginID, password)
	if err != nil {
		return fmt.Errorf("create first owner: %w", err)
	}
	logger.InfoContext(ctx, fmt.Sprintf("firstOwnerID: %d", firstOwnerID.Int()))

	// 5. find public default space
	if _, err := findPublicSpaceByKey(ctx, sysOwner, mbNonTxManager, authservice.PublicDefaultSpaceKey); err != nil {
		return fmt.Errorf("find public default space by key(%s): %w", authservice.PublicDefaultSpaceKey, err)
	}

	return nil
}

func findOrganizationAndSystemOwnerAndPublicDefaultSpace(ctx context.Context, systemAdmin authdomain.SystemAdminInterface, mbNonTxManager authservice.TransactionManager, organizationName string) (bool, error) {
	if _, err := findOrganizationByName(ctx, systemAdmin, mbNonTxManager, organizationName); err != nil {
		if !errors.Is(err, authservice.ErrOrganizationNotFound) {
			return false, fmt.Errorf("find organization by name: %w", err)
		}
		return false, nil
	}

	sysOwner, err := findSystemOwnerByOrganizationName(ctx, systemAdmin, mbNonTxManager, organizationName)
	if err != nil {
		return false, fmt.Errorf("find system owner by organization name: %w", err)
	}

	if _, err := findPublicSpaceByKey(ctx, sysOwner, mbNonTxManager, authservice.PublicDefaultSpaceKey); err != nil {
		if !errors.Is(err, authservice.ErrSpaceNotFound) {
			return false, fmt.Errorf("find public default space by key: %w", err)
		}
		return false, nil
	}

	return true, nil
}

func createOrganization(ctx context.Context, operator authdomain.SystemAdminInterface, mbTxManager, mbNonTxManager authservice.TransactionManager, organizationName string) (*authdomain.OrganizationID, error) {
	command := usecase.NewCreateOrganizationCommand(ctx, mbTxManager, mbNonTxManager)
	organizationID, err := command.Execute(ctx, operator, organizationName)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}
	return organizationID, nil
}

func createFirstOwnerToOrganization(ctx context.Context, operator authdomain.SystemOwnerInterface, mbTxManager, mbNonTxManager authservice.TransactionManager, loginID, password string) (*authdomain.UserID, error) {
	firstOwnerAddParam, err := authservice.NewCreateUserParameter(loginID, "Owner(cocotola)", password, "", "", "", "")
	if err != nil {
		return nil, fmt.Errorf("new UserAddParameter: %w", err)
	}
	addFirstOwnerCommand := usecase.NewCreateFirstOwnerCommand(mbTxManager, mbNonTxManager)
	firstOwnerID, err := addFirstOwnerCommand.Execute(ctx, operator, firstOwnerAddParam)
	if err != nil {
		return nil, fmt.Errorf("add first owner: %w", err)
	}
	return firstOwnerID, nil
}
