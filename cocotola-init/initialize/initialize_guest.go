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

func initGuest(ctx context.Context, systemToken authdomain.SystemToken, mbTxManager, mbNonTxManager authservice.TransactionManager, organizationName string, appName string) error {
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, appName+"-InitGuest"))

	sysAdmin := authdomain.NewSystemAdmin(systemToken)

	sysOwner, err := findSystemOwnerByOrganizationName(ctx, sysAdmin, mbNonTxManager, organizationName)
	if err != nil {
		return fmt.Errorf("findSystemOwnerByOrganizationName: %w", err)
	}

	guestLoginID := authdomain.NewGuestLoginID(organizationName)
	guestUserName := authdomain.NewGuestUserName(organizationName)
	// 1. check whether the guest user already exists
	{
		guest, err := findUserByLoginID(ctx, sysOwner, mbNonTxManager, guestLoginID)
		if err == nil {
			logger.InfoContext(ctx, fmt.Sprintf("guest already exists. id: %d", guest.GetUserID().Int()))
			return nil
		} else if !errors.Is(err, authservice.ErrUserNotFound) {
			return fmt.Errorf("find user by login id(%s): %w", guestLoginID, err)
		}
	}

	// 2. find public default space
	publicDefaultSpace, err := findPublicSpaceByKey(ctx, sysOwner, mbNonTxManager, authservice.PublicDefaultSpaceKey)
	if err != nil {
		return fmt.Errorf("find public default space by key(%s): %w", authservice.PublicDefaultSpaceKey, err)
	}

	// 3. add guest user
	if err := createGuestUser(ctx, mbTxManager, mbNonTxManager, sysOwner, guestLoginID, guestUserName, publicDefaultSpace.SpaceID); err != nil {
		return fmt.Errorf("createGuestUser: %w", err)
	}

	return nil
}

func createGuestUser(ctx context.Context, mbTxManager, mbNonTxManager authservice.TransactionManager, systemOwner authdomain.SystemOwnerInterface, guestLoginID, guestUserName string, _ *authdomain.SpaceID) error {
	// allowEffect := authservice.RBACAllowEffect
	// spaceObject := spaceID.GetRBACObject()
	aoeList := []authdomain.ActionObjectEffect{
		// guest can list decks in the "public" space
		// {Action: coreservice.ListDecksAction, Object: spaceObject, Effect: allowEffect},
		// // guest cat read all decks in the "public" space
		// {Action: coreservice.ReadDeckAction, Object: spaceObject, Effect: allowEffect},
	}

	addGuestCommand := usecase.NewCreateGuestCommand(mbTxManager, mbNonTxManager)
	addUserParam, err := authservice.NewCreateUserParameter(guestLoginID, guestUserName, "DUMMY_PASSWORD", "", "", "", "")
	if err != nil {
		return fmt.Errorf("NewCreateUserParameter: %w", err)
	}

	if _, err := addGuestCommand.Execute(ctx, systemOwner, addUserParam, aoeList); err != nil {
		return fmt.Errorf("execute addGuestCommand: %w", err)
	}

	return nil
}
