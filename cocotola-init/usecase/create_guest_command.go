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

type CreateGuestCommandGateway interface {
	WithTransaction(ctx context.Context, fn func(
		findPublicSpaceByKey authservice.FindPublicSpaceByKeyFunc,
		createUser authservice.CreateUserFunc,
		findUserGroupByKey authservice.FindUserGroupByKeyFunc,
		addUserToGroup authservice.AddUserToGroupFunc,
		attachPolicyToUser authservice.AttachPolicyToUserFunc,
	) (*authdomain.UserID, error)) (*authdomain.UserID, error)
}

type CreateGuestCommand struct {
	gw     CreateGuestCommandGateway
	logger *slog.Logger
}

func NewCreateGuestCommand(_ context.Context, gw CreateGuestCommandGateway) *CreateGuestCommand {
	return &CreateGuestCommand{
		gw:     gw,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "CreateGuestCommand")),
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
	userID, err := u.gw.WithTransaction(ctx, func(
		findPublicSpaceByKey authservice.FindPublicSpaceByKeyFunc,
		createUser authservice.CreateUserFunc,
		findUserGroupByKey authservice.FindUserGroupByKeyFunc,
		addUserToGroup authservice.AddUserToGroupFunc,
		attachPolicyToUser authservice.AttachPolicyToUserFunc,
	) (*authdomain.UserID, error) {
		// 1. find public default space
		publicDefaultSpace, err := findPublicSpaceByKey(ctx, operator, authservice.PublicDefaultSpaceKey)
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

		// 2. create new user
		userID, err := createUser(ctx, operator, param)
		if err != nil {
			return nil, fmt.Errorf("CreateUser: %w", err)
		}

		// 3. add user to public-group
		publicGroup, err := findUserGroupByKey(ctx, operator, authservice.PublicGroupKey)
		if err != nil {
			return nil, fmt.Errorf("find public group(%s): %w", authservice.PublicGroupKey, err)
		}
		if err := addUserToGroup(ctx, operator, userID, publicGroup.UserGroupID); err != nil {
			return nil, fmt.Errorf("AddUserToGroup: %w", err)
		}

		// 4. attach policy to user
		subject := userID.GetRBACSubject()
		for _, aoe := range aoeList {
			if err := attachPolicyToUser(ctx, operator, subject, aoe.Action, aoe.Object, aoe.Effect); err != nil {
				return nil, fmt.Errorf("AttachPolicyToUser: %w", err)
			}
		}

		return userID, nil
	})
	if err != nil {
		return nil, fmt.Errorf("WithTransaction: %w", err)
	}

	return userID, nil
}

func (u *CreateGuestCommand) callback(_ context.Context, _ authdomain.SystemOwnerInterface, _ *authdomain.UserID) error {
	return nil
}
