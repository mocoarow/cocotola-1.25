package usecase

import (
	"context"
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type CreateGuestCommand struct {
	txManager    authservice.TransactionManager
	nonTxManager authservice.TransactionManager
}

func NewCreateGuestCommand(txManager authservice.TransactionManager, nonTxManager authservice.TransactionManager) *CreateGuestCommand {
	return &CreateGuestCommand{
		txManager:    txManager,
		nonTxManager: nonTxManager,
	}
}

func (u *CreateGuestCommand) Execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter) (*authdomain.UserID, error) {
	// 1. Check authorization
	if err := u.checkAuthorization(ctx, operator); err != nil {
		return nil, fmt.Errorf("checkAuthorization: %w", err)
	}

	// 2. Execute
	newUserID, err := u.execute(ctx, operator, param)
	if err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}

	// 3. Callback
	if err := u.callback(ctx, operator, newUserID); err != nil {
		return nil, fmt.Errorf("callback: %w", err)
	}

	return newUserID, nil
}

func (u *CreateGuestCommand) checkAuthorization(_ context.Context, _ authdomain.SystemOwnerInterface) error {
	return nil
}

func (u *CreateGuestCommand) execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter) (*authdomain.UserID, error) {
	publicDefaultSpace, err := authservice.FindPublicSpaceByKey(ctx, operator, u.nonTxManager, authservice.PublicDefaultSpaceKey)
	if err != nil {
		return nil, fmt.Errorf("find public default space by key(%s): %w", authservice.PublicDefaultSpaceKey, err)
	}

	spaceObject := publicDefaultSpace.SpaceID.GetRBACObject()

	aoeList := []libdomain.ActionObjectEffect{
		// guest can list decks in the "public" space
		{Action: libservice.ListDecksAction, Object: spaceObject, Effect: authservice.RBACAllowEffect},
		// guest can read all decks in the "public" space
		{Action: libservice.ReadDeckAction, Object: spaceObject, Effect: authservice.RBACAllowEffect},
	}

	fn2 := func(rf authservice.RepositoryFactory) (*authdomain.UserID, error) {
		userID, err := AddUser(ctx, operator, rf, param, aoeList)
		if err != nil {
			return nil, fmt.Errorf("AddUser: %w", err)
		}

		// u.logger.InfoContext(ctx, fmt.Sprintf("personalSpaceID: %d", spaceID.Int()))
		return userID, nil
	}
	userID, err := libservice.Do1(ctx, u.txManager, fn2)
	if err != nil {
		return nil, fmt.Errorf("Do1: %w", err)
	}
	return userID, nil
}

func (u *CreateGuestCommand) callback(_ context.Context, _ authdomain.SystemOwnerInterface, _ *authdomain.UserID) error {
	return nil
}
