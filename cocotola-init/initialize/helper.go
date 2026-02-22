package initialize

import (
	"context"
	"fmt"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

func findSystemOwnerByOrganizationName(ctx context.Context, systemAdmin authdomain.SystemAdminInterface, dbc *libgateway.DBConnection, organizationName string) (*authdomain.SystemOwner, error) {
	userRepo := authgateway.NewUserRepository(dbc)
	sysOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, organizationName)
	if err != nil {
		return nil, fmt.Errorf("find system owner by organization name(%s): %w", organizationName, err)
	}

	return sysOwner, nil
}

func findUserByLoginID(ctx context.Context, systemOwner authdomain.SystemOwnerInterface, dbc *libgateway.DBConnection, loginID string) (*authdomain.User, error) {
	userRepo := authgateway.NewUserRepository(dbc)
	user, err := userRepo.FindUserByLoginID(ctx, systemOwner, loginID)
	if err != nil {
		return nil, fmt.Errorf("find user by login id(%s): %w", loginID, err)
	}

	return user, nil
}

func findPublicSpaceByKey(ctx context.Context, operator authdomain.SystemOwnerInterface, dbc *libgateway.DBConnection, key string) (*authdomain.Space, error) {
	spaceRepo := authgateway.NewSpaceRepository(dbc)
	space, err := spaceRepo.FindPublicSpaceByKey(ctx, operator, key)
	if err != nil {
		return nil, fmt.Errorf("find public space by key(%s): %w", key, err)
	}

	return space, nil
}
