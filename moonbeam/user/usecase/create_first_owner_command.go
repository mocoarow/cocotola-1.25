package usecase

import (
	"context"
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
	libservice "github.com/mocoarow/cocotola-1.25/moonbeam/lib/service"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type CreateFirstOwnerCommand struct {
	txManager    service.TransactionManager
	nonTxManager service.TransactionManager
}

func NewCreateFirstOwnerCommand(txManager service.TransactionManager, nonTxManager service.TransactionManager) *CreateFirstOwnerCommand {
	return &CreateFirstOwnerCommand{
		txManager:    txManager,
		nonTxManager: nonTxManager,
	}
}

func (u *CreateFirstOwnerCommand) checkAuthorization(ctx context.Context, operator domain.SystemOwnerInterface, param *service.CreateUserParameter) error {
	fn1 := func(rf service.RepositoryFactory) error {
		authorizationManager, err := rf.NewAuthorizationManager(ctx)
		if err != nil {
			return fmt.Errorf("failed to NewAuthorizationManager: %w", err)
		}
		rbacAllUserRolesObject := domain.NewRBACAllUserRolesObjectFromOrganization(operator.GetOrganizationID())

		// Can "operator" "CreateOwner" "*" ?
		ok, err := authorizationManager.CheckAuthorization(ctx, operator, service.CreateOwnerAction, rbacAllUserRolesObject)
		if err != nil {
			return fmt.Errorf("CheckAuthorization: %w", err)
		} else if !ok {
			return libdomain.ErrPermissionDenied
		}
		return nil
	}
	if err := libservice.Do0(ctx, u.txManager, fn1); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

func (u *CreateFirstOwnerCommand) Execute(ctx context.Context, operator domain.SystemOwnerInterface, param *service.CreateUserParameter) (*domain.UserID, error) {
	if err := u.checkAuthorization(ctx, operator, param); err != nil {
		return nil, err //nolint:wrapcheck
	}

	fn2 := func(rf service.RepositoryFactory) (*domain.UserID, error) {
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

		ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, operator, service.OwnerGroupKey)
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
		action := service.CreateUserAction
		object := service.AnyObject
		effect := service.RBACAllowEffect

		if err := authorizationManager.AttachPolicyToUserBySystemOwner(ctx, operator, subject, action, object, effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemOwner: %w", err)
		}
		return firstOwnerID, nil
	}
	firstOwnerID, err := libservice.Do1(ctx, u.nonTxManager, fn2)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return firstOwnerID, nil
}
