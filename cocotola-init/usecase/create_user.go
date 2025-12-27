package usecase

import (
	"context"
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

func AddUser(ctx context.Context, operator authdomain.UserInterface, rf authservice.RepositoryFactory, param *authservice.CreateUserParameter, aoeList []libdomain.ActionObjectEffect) (*authdomain.UserID, error) {
	userRepo := rf.NewUserRepository(ctx)
	userGroupRepo := rf.NewUserGroupRepository(ctx)
	authorizationManager, err := rf.NewAuthorizationManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewAuthorizationManager: %w", err)
	}

	// 1. create new user
	userID, err := userRepo.CreateUser(ctx, operator, param)
	if err != nil {
		return nil, fmt.Errorf("m.userRepo.CreateUser. err: %w", err)
	}

	// 2. add user to public-group
	publicGroup, err := userGroupRepo.FindUserGroupByKey(ctx, operator, authservice.PublicGroupKey)
	if err != nil {
		return nil, fmt.Errorf("find public group(%s): %w", authservice.PublicGroupKey, err)
	}
	if err := authorizationManager.AddUserToGroup(ctx, operator, userID, publicGroup.UserGroupID); err != nil {
		return nil, fmt.Errorf("AddUserToGroup: %w", err)
	}

	// 3. attach policy to user
	subject := userID.GetRBACSubject()
	for _, aoe := range aoeList {
		if err := authorizationManager.AttachPolicyToUser(ctx, operator, subject, aoe.Action, aoe.Object, aoe.Effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUser: %w", err)
		}
	}

	return userID, nil
}
