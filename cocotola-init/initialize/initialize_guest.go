package initialize

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	usecase "github.com/mocoarow/cocotola-1.25/cocotola-init/usecase"
)

func initGuest(ctx context.Context, systemToken authdomain.SystemToken, dbc *libgateway.DBConnection, organizationName string, appName string) error {
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, appName+"-InitGuest"))

	sysAdmin := authdomain.NewSystemAdmin(systemToken)

	sysOwner, err := findSystemOwnerByOrganizationName(ctx, sysAdmin, dbc, organizationName)
	if err != nil {
		return fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	guestLoginID := authdomain.NewGuestLoginID(organizationName)
	guestUserName := authdomain.NewGuestUserName(organizationName)
	// 1. check whether the guest user already exists
	{
		guest, err := findUserByLoginID(ctx, sysOwner, dbc, guestLoginID)
		if err == nil {
			logger.InfoContext(ctx, fmt.Sprintf("guest already exists. id: %d", guest.GetUserID().Int()))
			return nil
		} else if !errors.Is(err, authservice.ErrUserNotFound) {
			return fmt.Errorf("find user by login id(%s): %w", guestLoginID, err)
		}
	}

	// 2. find public default space
	publicDefaultSpace, err := findPublicSpaceByKey(ctx, sysOwner, dbc, authservice.PublicDefaultSpaceKey)
	if err != nil {
		return fmt.Errorf("find public default space by key(%s): %w", authservice.PublicDefaultSpaceKey, err)
	}

	// 3. add guest user
	if err := createGuestUser(ctx, dbc, sysOwner, guestLoginID, guestUserName, publicDefaultSpace.SpaceID); err != nil {
		return fmt.Errorf("createGuestUser: %w", err)
	}

	return nil
}

func createGuestUser(ctx context.Context, dbc *libgateway.DBConnection, systemOwner authdomain.SystemOwnerInterface, guestLoginID, guestUserName string, _ *authdomain.SpaceID) error {
	createGuestCommandGateway := NewCreateGuestCommandGateway(dbc)
	addGuestCommand := usecase.NewCreateGuestCommand(ctx, createGuestCommandGateway)
	addUserParam, err := authservice.NewCreateUserParameter(guestLoginID, guestUserName, "DUMMY_PASSWORD", "", "", "", "")
	if err != nil {
		return fmt.Errorf("NewCreateUserParameter: %w", err)
	}

	if _, err := addGuestCommand.Execute(ctx, systemOwner, addUserParam); err != nil {
		return fmt.Errorf("execute addGuestCommand: %w", err)
	}

	return nil
}

type CreateGuestCommandGateway struct {
	dbc *libgateway.DBConnection
}

func NewCreateGuestCommandGateway(dbc *libgateway.DBConnection) *CreateGuestCommandGateway {
	return &CreateGuestCommandGateway{
		dbc: dbc,
	}
}

func (gw *CreateGuestCommandGateway) WithTransaction(ctx context.Context, fn func(
	findPublicSpaceByKey authservice.FindPublicSpaceByKeyFunc,
	createUser authservice.CreateUserFunc,
	findUserGroupByKey authservice.FindUserGroupByKeyFunc,
	addUserToGroup authservice.AddUserToGroupFunc,
	attachPolicyToUser authservice.AttachPolicyToUserFunc,
) (*authdomain.UserID, error)) (*authdomain.UserID, error) {
	var userID *authdomain.UserID
	err := gw.dbc.DB.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
		spaceRepo := authgateway.NewSpaceRepository(gw.dbc)
		userRepo := authgateway.NewUserRepository(gw.dbc)
		groupRepo := authgateway.NewUserGroupRepository(gw.dbc)
		authManager, err := authgateway.NewAuthorizationManager(ctx, gw.dbc)
		if err != nil {
			return fmt.Errorf("new AuthorizationManager: %w", err)
		}

		userID, err = fn(
			spaceRepo.FindPublicSpaceByKey,
			userRepo.CreateUser,
			groupRepo.FindUserGroupByKey,
			authManager.AddUserToGroup,
			authManager.AttachPolicyToUser)
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
