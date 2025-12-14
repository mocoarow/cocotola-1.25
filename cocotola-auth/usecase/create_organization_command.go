package usecase

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type CreateOrganizationCommand struct {
	txManager    service.TransactionManager
	nonTxManager service.TransactionManager
	logger       *slog.Logger
}

func NewCreateOrganizationCommand(_ context.Context, txManager, nonTxManager service.TransactionManager) *CreateOrganizationCommand {
	return &CreateOrganizationCommand{
		txManager:    txManager,
		nonTxManager: nonTxManager,
		logger:       slog.Default().With(slog.String(libdomain.LoggerNameKey, "CreateOrganizationCommand")),
	}
}

func (u *CreateOrganizationCommand) Execute(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.OrganizationID, error) {
	fn := func(rf service.RepositoryFactory) (*domain.OrganizationID, error) {
		userRepo := rf.NewUserRepository(ctx)

		// system-admin creates organization and system-owner
		organizationID, err := u.executeCreatingOrganizationProcessBySystemAdmin(ctx, operator, rf, organizationName)
		if err != nil {
			return nil, fmt.Errorf("executeCreatingOrganizationProcessBySystemAdmin: %w", err)
		}

		// system-owner creates organization resources
		systemOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, operator, organizationName)
		if err != nil {
			return nil, fmt.Errorf("FindSystemOwnerByOrganizationName: %w", err)
		}
		if err := u.executeCreatingOrganizationProcessBySystemOwner(ctx, systemOwner, rf, organizationID); err != nil {
			return nil, fmt.Errorf("executeCreatingOrganizationProcessBySystemOwner: %w", err)
		}

		return organizationID, nil
	}
	organizationID, err := libservice.Do1(ctx, u.txManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return organizationID, nil
}

func (u *CreateOrganizationCommand) executeCreatingOrganizationProcessBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, rf service.RepositoryFactory, organizationName string) (*domain.OrganizationID, error) {
	orgRepo := rf.NewOrganizationRepository(ctx)
	userRepo := rf.NewUserRepository(ctx)
	authorizationManager, err := rf.NewAuthorizationManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to NewAuthorizationManager: %w", err)
	}

	// 1. create organization
	organizationID, err := orgRepo.CreateOrganization(ctx, operator, organizationName)
	if err != nil {
		return nil, fmt.Errorf("CreateOrganization: %w", err)
	}

	// 2. create "system-owner" user
	// 3. attach policy to "system-owner" user
	systemOwnerID, err := u.createSystemOwnerForOrganization(ctx, operator, userRepo, authorizationManager, organizationID)
	if err != nil {
		return nil, fmt.Errorf("createSystemOwnerForOrganization: %w", err)
	}
	u.logger.InfoContext(ctx, fmt.Sprintf("organizationID: %d, systemOwnerID: %d", organizationID.Int(), systemOwnerID.Int()))

	return organizationID, nil
}

func (u *CreateOrganizationCommand) executeCreatingOrganizationProcessBySystemOwner(ctx context.Context, operator domain.SystemOwnerInterface, rf service.RepositoryFactory, organizationID *domain.OrganizationID) error {
	userGroupRepo := rf.NewUserGroupRepository(ctx)
	authorizationManager, err := rf.NewAuthorizationManager(ctx)
	if err != nil {
		return fmt.Errorf("failed to NewAuthorizationManager: %w", err)
	}
	spaceManager, err := rf.NewSpaceManager(ctx)
	if err != nil {
		return fmt.Errorf("NewSpaceManager: %w", err)
	}

	// 4. create owner-group
	// 5. attach policy to "owner" group
	if _, err := u.createOwnerGroupForOrganization(ctx, operator, userGroupRepo, authorizationManager, organizationID); err != nil {
		return fmt.Errorf("addOwnergroupToOrganization: %w", err)
	}

	// 7. create public-group
	if _, err := userGroupRepo.CreatePublicGroup(ctx, operator, organizationID); err != nil {
		return fmt.Errorf("create public group: %w", err)
	}

	// 9. create public default space
	if _, err := spaceManager.CreatePublicDefaultSpace(ctx, operator); err != nil {
		return fmt.Errorf("create public space(%s): %w", service.PublicDefaultSpaceKey, err)
	}
	return nil
}

type ActionObjectEffect struct {
	Action domain.RBACAction
	Object domain.RBACObject
	Effect domain.RBACEffect
}

func (u *CreateOrganizationCommand) createSystemOwnerForOrganization(ctx context.Context, operator domain.SystemAdminInterface, userRepo service.UserRepository, authorizationManager service.AuthorizationManager, organizationID *domain.OrganizationID) (*domain.UserID, error) {
	systemOwnerID, err := userRepo.CreateSystemOwner(ctx, operator, organizationID)
	if err != nil {
		return nil, fmt.Errorf("CreateSystemOwner: %w", err)
	}

	// 3. attach policy to "system-owner" user
	rbacSystemOwner := systemOwnerID.GetRBACSubject()
	// rbacAllUserRolesObject := domain.NewRBACAllUserRolesObjectFromOrganization(organizationID)
	for _, aoe := range []ActionObjectEffect{
		{ // "system-owner" "can" "CreateUser" "*"
			Action: service.CreateUserAction,
			Object: service.AnyObject,
			Effect: service.RBACAllowEffect,
		},
		// { //"system-owner" user "can" "unset" "all-user-roles"
		// 	Action: service.RBACUnsetAction,
		// 	Object: rbacAllUserRolesObject,
		// 	Effect: service.RBACAllowEffect,
		// },
	} {
		if err := authorizationManager.AttachPolicyToUserBySystemAdmin(ctx, operator, organizationID, rbacSystemOwner, aoe.Action, aoe.Object, aoe.Effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemAdmin: %w", err)
		}
	}

	return systemOwnerID, nil
}

func (u *CreateOrganizationCommand) createOwnerGroupForOrganization(ctx context.Context, operator domain.SystemOwnerInterface, userGroupRepo service.UserGroupRepository, authorizationManager service.AuthorizationManager, organizationID *domain.OrganizationID) (*domain.UserGroupID, error) {
	u.logger.InfoContext(ctx, "createOwnerGroupForOrganization", "organizationID", organizationID.Int())
	// 4. create owner-group
	ownerGroupID, err := userGroupRepo.CreateOwnerGroup(ctx, operator, organizationID)
	if err != nil {
		return nil, fmt.Errorf("CreateOwnerGroup: %w", err)
	}

	// 5. attach policy to "owner" group
	rbacOwnerGroup := domain.NewRBACRoleFromGroup(organizationID, ownerGroupID)
	// rbacAllUserRolesObject := domain.NewRBACAllUserRolesObjectFromOrganization(organizationID)

	for _, aoe := range []ActionObjectEffect{
		{ // "owner" group "can" "CreateUser" "*"
			Action: service.CreateUserAction,
			Object: service.AnyObject,
			Effect: service.RBACAllowEffect,
		},
		// { // "owner" group "can" "unset" "all-user-roles"
		// 	Action: service.RBACUnsetAction,
		// 	Object: rbacAllUserRolesObject,
		// 	Effect: service.RBACAllowEffect,
		// },
	} {
		if err := authorizationManager.AttachPolicyToUserBySystemOwner(ctx, operator, rbacOwnerGroup, aoe.Action, aoe.Object, aoe.Effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemOwner: %w", err)
		}
	}
	return ownerGroupID, nil
}
