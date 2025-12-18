package usecase

import (
	"context"
	"fmt"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"
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

func (u *CreateGuestCommand) Execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter, aoeList []authdomain.ActionObjectEffect) (*authdomain.UserID, error) {
	// 1. Check authorization
	if err := u.checkAuthorization(ctx, operator); err != nil {
		return nil, fmt.Errorf("checkAuthorization: %w", err)
	}

	// 2. Execute
	newUserID, err := u.execute(ctx, operator, param, aoeList)
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

func (u *CreateGuestCommand) execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter, aoeList []authdomain.ActionObjectEffect) (*authdomain.UserID, error) {
	userID, err := libservice.Do1(ctx, u.txManager, func(rf authservice.RepositoryFactory) (*authdomain.UserID, error) {
		return AddUser(ctx, operator, rf, param, aoeList)
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return userID, nil
}

func (u *CreateGuestCommand) callback(_ context.Context, _ authdomain.SystemOwnerInterface, _ *authdomain.UserID) error {
	return nil
}
