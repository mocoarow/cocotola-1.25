package usecase

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type CreateFirstOwnerCommand struct {
	txManager    authservice.TransactionManager
	nonTxManager authservice.TransactionManager
	logger       *slog.Logger
}

func NewCreateFirstOwnerCommand(txManager authservice.TransactionManager, nonTxManager authservice.TransactionManager) *CreateFirstOwnerCommand {
	return &CreateFirstOwnerCommand{
		txManager:    txManager,
		nonTxManager: nonTxManager,
		logger:       slog.Default().With(slog.String(libdomain.LoggerNameKey, "CreateFirstOwnerCommand")),
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
		firstOwner, err := u.createFirstOwner(ctx, operator, param, userRepo)
		if err != nil {
			return nil, fmt.Errorf("CreateUser: %w", err)
		}

		// 2. add owner to owner-group
		if err := u.addFirstOwnerToOwnerGroup(ctx, operator, firstOwner.GetUserID(), userGroupRepo, authorizationManager); err != nil {
			return nil, fmt.Errorf("addUserToGroup: %w", err)
		}

		// 3. attach policy to "first-owner" user
		// first owner can create users
		subject := firstOwner.GetUserID().GetRBACSubject()
		action := authservice.CreateUserAction
		object := authservice.AnyObject
		effect := authservice.RBACAllowEffect

		if err := authorizationManager.AttachPolicyToUserBySystemOwner(ctx, operator, subject, action, object, effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemOwner: %w", err)
		}

		// 4. create personal space for first owner
		spaceManager, err := rf.NewSpaceManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("NewSpaceManager: %w", err)
		}
		createPersonalSpaceParameter := authservice.CreatePersonalSpaceParameter{
			UserID:  firstOwner.UserID,
			KeyName: authservice.NewPersonalSpaceKey(firstOwner.UserID.Int()),
			Name:    authservice.NewPersonalSpaceName(firstOwner.GetLoginID()),
		}
		spaceID, err := spaceManager.CreatePersonalSpace(ctx, operator, &createPersonalSpaceParameter)
		if err != nil {
			return nil, fmt.Errorf("CreatePersonalSpace: %w", err)
		}
		u.logger.InfoContext(ctx, fmt.Sprintf("personalSpaceID: %d", spaceID.Int()))

		return firstOwner.UserID, nil
	}
	firstOwnerID, err := libservice.Do1(ctx, u.txManager, fn2)
	if err != nil {
		return nil, fmt.Errorf("Do1: %w", err)
	}

	return firstOwnerID, nil
}

func (u *CreateFirstOwnerCommand) createFirstOwner(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter, userRepo authservice.UserRepository) (*authdomain.User, error) {
	// 1. create owner
	firstOwnerID, err := userRepo.CreateUser(ctx, operator, param)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	// 2. attach policy to "first-owner" user
	firstOwner, err := userRepo.FindUserByID(ctx, operator, firstOwnerID)
	if err != nil {
		return nil, fmt.Errorf("FindUserByLoginID: %w", err)
	}

	return firstOwner, nil
}

func (u *CreateFirstOwnerCommand) addFirstOwnerToOwnerGroup(ctx context.Context, operator authdomain.SystemOwnerInterface, userID *authdomain.UserID, userGroupRepo authservice.UserGroupRepository, authorizationManager authservice.AuthorizationManager) error {
	ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, operator, authservice.OwnerGroupKey)
	if err != nil {
		return fmt.Errorf("FindUserGroupByKey: %w", err)
	}

	// 2. add owner to owner-group
	if err := authorizationManager.AddUserToGroup(ctx, operator, userID, ownerGroup.UserGroupID); err != nil {
		return fmt.Errorf("AddUserToGroup: %w", err)
	}
	return nil
}
