package usecase

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type CreateFirstOwnerCommandGateway interface {
	WithTransaction(ctx context.Context, fn func(
		createUser authservice.CreateUserFunc,
		findUserByID authservice.FindUserByIDFunc,
		findUserGroupByKey authservice.FindUserGroupByKeyFunc,
		addUserToGroup authservice.AddUserToGroupFunc,
		attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc,
		createPersonalSpace authservice.CreatePersonalSpaceFunc,
	) (*authdomain.UserID, error)) (*authdomain.UserID, error)
}

type CreateFirstOwnerCommand struct {
	gw     CreateFirstOwnerCommandGateway
	logger *slog.Logger
}

func NewCreateFirstOwnerCommand(_ context.Context, gw CreateFirstOwnerCommandGateway) *CreateFirstOwnerCommand {
	return &CreateFirstOwnerCommand{
		gw:     gw,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "CreateFirstOwnerCommand")),
	}
}

func (u *CreateFirstOwnerCommand) Execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter) (*authdomain.UserID, error) {
	// 1. Check authorization
	if err := u.checkAuthorization(ctx, operator, param); err != nil {
		return nil, fmt.Errorf("checkAuthorization: %w", err)
	}

	// 2. Execute
	firstOwnerID, err := u.execute(ctx, operator, param)
	if err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}

	// 3. Callback
	if err := u.callback(ctx, operator, firstOwnerID); err != nil {
		return nil, fmt.Errorf("callback: %w", err)
	}

	return firstOwnerID, nil
}

func (u *CreateFirstOwnerCommand) checkAuthorization(_ context.Context, _ authdomain.SystemOwnerInterface, _ *authservice.CreateUserParameter) error {
	// system-owner can create owner
	return nil
}

func (u *CreateFirstOwnerCommand) execute(ctx context.Context, operator authdomain.SystemOwnerInterface, param *authservice.CreateUserParameter) (*authdomain.UserID, error) {
	firstOwnerID, err := u.gw.WithTransaction(ctx, func(
		createUser authservice.CreateUserFunc,
		findUserByID authservice.FindUserByIDFunc,
		findUserGroupByKey authservice.FindUserGroupByKeyFunc,
		addUserToGroup authservice.AddUserToGroupFunc,
		attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc,
		createPersonalSpace authservice.CreatePersonalSpaceFunc,
	) (*authdomain.UserID, error) {
		// 1. create owner
		firstOwnerID, err := createUser(ctx, operator, param)
		if err != nil {
			return nil, fmt.Errorf("CreateUser: %w", err)
		}

		// 2. find first owner
		firstOwner, err := findUserByID(ctx, operator, firstOwnerID)
		if err != nil {
			return nil, fmt.Errorf("FindUserByID: %w", err)
		}

		// 3. find owner group
		ownerGroup, err := findUserGroupByKey(ctx, operator, authservice.OwnerGroupKey)
		if err != nil {
			return nil, fmt.Errorf("FindUserGroupByKey: %w", err)
		}

		// 4. add owner to owner-group
		if err := addUserToGroup(ctx, operator, firstOwnerID, ownerGroup.UserGroupID); err != nil {
			return nil, fmt.Errorf("AddUserToGroup: %w", err)
		}

		// 5. attach policy to "first-owner" user
		// first owner can create users
		subject := firstOwnerID.GetRBACSubject()
		action := authservice.CreateUserAction
		object := authservice.AnyObject
		effect := authservice.RBACAllowEffect

		if err := attachPolicyToUserBySystemOwner(ctx, operator, subject, action, object, effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemOwner: %w", err)
		}

		// 6. create personal space for first owner
		createPersonalSpaceParameter := authservice.CreatePersonalSpaceParameter{
			UserID:  firstOwner.UserID,
			KeyName: authservice.NewPersonalSpaceKey(firstOwner.UserID.Int()),
			Name:    authservice.NewPersonalSpaceName(firstOwner.GetLoginID()),
		}
		spaceID, err := createPersonalSpace(ctx, operator, &createPersonalSpaceParameter)
		if err != nil {
			return nil, fmt.Errorf("CreatePersonalSpace: %w", err)
		}
		u.logger.InfoContext(ctx, fmt.Sprintf("personalSpaceID: %d", spaceID.Int()))

		return firstOwnerID, nil
	})
	if err != nil {
		return nil, fmt.Errorf("WithTransaction: %w", err)
	}

	return firstOwnerID, nil
}

func (u *CreateFirstOwnerCommand) callback(_ context.Context, _ authdomain.SystemOwnerInterface, _ *authdomain.UserID) error {
	return nil
}
