package usecase

import (
	"context"
	"fmt"

	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type CreateFirstOwnerCommand struct {
	txManager    authservice.TransactionManager
	nonTxManager authservice.TransactionManager
}

func NewCreateFirstOwnerCommand(txManager authservice.TransactionManager, nonTxManager authservice.TransactionManager) *CreateFirstOwnerCommand {
	return &CreateFirstOwnerCommand{
		txManager:    txManager,
		nonTxManager: nonTxManager,
	}
}

func (u *CreateFirstOwnerCommand) checkAuthorization(_ context.Context, _ authdomain.SystemOwnerInterface, _ *authservice.CreateUserParameter) error {
	// system-owner can create owner
	return nil
}

func (u *CreateFirstOwnerCommand) Execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter) (*authdomain.UserID, error) {
	if err := u.checkAuthorization(ctx, operator, param); err != nil {
		return nil, fmt.Errorf("checkAuthorization: %w", err)
	}

	fn2 := func(rf authservice.RepositoryFactory) (*authdomain.UserID, error) {
		userRepo := rf.NewUserRepository(ctx)
		userGroupRepo := rf.NewUserGroupRepository(ctx)
		authorizationManager, err := rf.NewAuthorizationManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to NewAuthorizationManager: %w", err)
		}
		// 1. create owner
		firstOwnerID, err := userRepo.CreateUser(ctx, operator, param)
		if err != nil {
			return nil, fmt.Errorf("CreateUser: %w", err)
		}

		ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, operator, authservice.OwnerGroupKey)
		if err != nil {
			return nil, fmt.Errorf("FindUserGroupByKey: %w", err)
		}

		// 2. add owner to owner-group
		if err := authorizationManager.AddUserToGroup(ctx, operator, firstOwnerID, ownerGroup.UserGroupID); err != nil {
			return nil, fmt.Errorf("AddUserToGroup: %w", err)
		}

		// 3. attach policy to "first-owner" user
		firstOwner, err := userRepo.FindUserByID(ctx, operator, firstOwnerID)
		if err != nil {
			return nil, fmt.Errorf("FindUserByLoginID: %w", err)
		}

		// first owner can create users
		subject := firstOwner.GetUserID().GetRBACSubject()
		action := authservice.CreateUserAction
		object := authservice.AnyObject
		effect := authservice.RBACAllowEffect

		if err := authorizationManager.AttachPolicyToUserBySystemOwner(ctx, operator, subject, action, object, effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemOwner: %w", err)
		}
		return firstOwnerID, nil
	}
	firstOwnerID, err := libservice.Do1(ctx, u.txManager, fn2)
	if err != nil {
		return nil, fmt.Errorf("Do1: %w", err)
	}

	return firstOwnerID, nil
}
