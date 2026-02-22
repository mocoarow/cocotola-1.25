package usecase

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type CreateOrganizationCommandGateway interface {
	// authservice.UserRepositoryCreateSystemOwner
	WithTransaction(ctx context.Context, fn func(
		createOrganization authservice.CreateOrganizationFunc,
		createSystemOwner authservice.CreateSystemOwnerFunc,
		findSystemOwnerByOrganizationName authservice.FindSystemOwnerByOrganizationNameFunc,
		attachPolicyToUserBySystemAdmin authservice.AttachPolicyToUserBySystemAdminFunc,
		createOwnerGroup authservice.CreateOwnerGroupFunc,
		attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc,
		createPublicGroup authservice.CreatePublicGroupFunc,
		createPublicDefaultSpace authservice.CreatePublicDefaultSpaceFunc,
	) (*authdomain.OrganizationID, error)) (*authdomain.OrganizationID, error)
}

type CreateOrganizationCommand struct {
	gw     CreateOrganizationCommandGateway
	logger *slog.Logger
}

func NewCreateOrganizationCommand(_ context.Context, gw CreateOrganizationCommandGateway) *CreateOrganizationCommand {
	return &CreateOrganizationCommand{
		gw:     gw,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "CreateOrganizationCommand")),
	}
}

func (u *CreateOrganizationCommand) Execute(ctx context.Context, operator authdomain.SystemAdminInterface, organizationName string) (*authdomain.OrganizationID, error) {
	organizationID, err := u.gw.WithTransaction(ctx, func(
		createOrganization authservice.CreateOrganizationFunc,
		createSystemOwner authservice.CreateSystemOwnerFunc,
		findSystemOwnerByOrganizationName authservice.FindSystemOwnerByOrganizationNameFunc,
		attachPolicyToUserBySystemAdmin authservice.AttachPolicyToUserBySystemAdminFunc,
		createOwnerGroup authservice.CreateOwnerGroupFunc,
		attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc,
		createPublicGroup authservice.CreatePublicGroupFunc,
		createPublicDefaultSpace authservice.CreatePublicDefaultSpaceFunc,
	) (*authdomain.OrganizationID, error) {
		organizationID, err := u.executeCreatingOrganizationProcessBySystemAdmin(ctx, operator, createOrganization, createSystemOwner, attachPolicyToUserBySystemAdmin, organizationName)
		if err != nil {
			return nil, fmt.Errorf("executeCreatingOrganizationProcessBySystemAdmin: %w", err)
		}

		// system-owner creates organization resources
		systemOwner, err := findSystemOwnerByOrganizationName(ctx, operator, organizationName)
		if err != nil {
			return nil, fmt.Errorf("FindSystemOwnerByOrganizationName: %w", err)
		}
		if err := u.executeCreatingOrganizationProcessBySystemOwner(ctx, systemOwner,
			createOwnerGroup, attachPolicyToUserBySystemOwner, createPublicGroup, createPublicDefaultSpace,
			organizationID); err != nil {
			return nil, fmt.Errorf("executeCreatingOrganizationProcessBySystemOwner: %w", err)
		}

		return organizationID, nil
	})
	if err != nil {
		return nil, fmt.Errorf("do in transaction: %w", err)
	}

	return organizationID, nil
}

func (u *CreateOrganizationCommand) executeCreatingOrganizationProcessBySystemAdmin(ctx context.Context, operator authdomain.SystemAdminInterface, createOrganization authservice.CreateOrganizationFunc, createSystemOwner authservice.CreateSystemOwnerFunc, attachPolicyToUserBySystemAdmin authservice.AttachPolicyToUserBySystemAdminFunc, organizationName string) (*authdomain.OrganizationID, error) {
	// 1. create organization
	organizationID, err := createOrganization(ctx, operator, organizationName)
	if err != nil {
		return nil, fmt.Errorf("CreateOrganization: %w", err)
	}

	// 2. create "system-owner" user
	// 3. attach policy to "system-owner" user
	systemOwnerID, err := createSystemOwner(ctx, operator, organizationID)
	if err != nil {
		return nil, fmt.Errorf("CreateSystemOwner: %w", err)
	}

	// 3. attach policy to "system-owner" user
	rbacSystemOwner := systemOwnerID.GetRBACSubject()
	// rbacAllUserRolesObject := authdomain.NewRBACAllUserRolesObjectFromOrganization(organizationID)
	for _, aoe := range []libdomain.ActionObjectEffect{
		{ // "system-owner" "can" "CreateUser" "*"
			Action: authservice.CreateUserAction,
			Object: authservice.AnyObject,
			Effect: authservice.RBACAllowEffect,
		},
		// { //"system-owner" user "can" "unset" "all-user-roles"
		// 	Action: authservice.RBACUnsetAction,
		// 	Object: rbacAllUserRolesObject,
		// 	Effect: authservice.RBACAllowEffect,
		// },
	} {
		if err := attachPolicyToUserBySystemAdmin(ctx, operator, organizationID, rbacSystemOwner, aoe.Action, aoe.Object, aoe.Effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemAdmin: %w", err)
		}
	}
	u.logger.InfoContext(ctx, fmt.Sprintf("organizationID: %d, systemOwnerID: %d", organizationID.Int(), systemOwnerID.Int()))

	return organizationID, nil
}

func (u *CreateOrganizationCommand) executeCreatingOrganizationProcessBySystemOwner(ctx context.Context, operator authdomain.SystemOwnerInterface, createOwnerGroup authservice.CreateOwnerGroupFunc, attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc, createPublicGroup authservice.CreatePublicGroupFunc, createPublicDefaultSpace authservice.CreatePublicDefaultSpaceFunc, organizationID *authdomain.OrganizationID) error {
	// 4. create owner-group
	// 5. attach policy to "owner" group
	if _, err := u.createOwnerGroupForOrganization(ctx, operator, createOwnerGroup, attachPolicyToUserBySystemOwner, organizationID); err != nil {
		return fmt.Errorf("addOwnergroupToOrganization: %w", err)
	}

	// 7. create public-group
	if _, err := createPublicGroup(ctx, operator, organizationID); err != nil {
		return fmt.Errorf("create public group: %w", err)
	}

	// 9. create public default space
	if _, err := createPublicDefaultSpace(ctx, operator); err != nil {
		return fmt.Errorf("create public space(%s): %w", authservice.PublicDefaultSpaceKey, err)
	}
	return nil
}

func (u *CreateOrganizationCommand) createOwnerGroupForOrganization(ctx context.Context, operator authdomain.SystemOwnerInterface, createOwnerGroup authservice.CreateOwnerGroupFunc, attachPolicyToUserBySystemOwner authservice.AttachPolicyToUserBySystemOwnerFunc, organizationID *authdomain.OrganizationID) (*authdomain.UserGroupID, error) {
	u.logger.InfoContext(ctx, "createOwnerGroupForOrganization", "organizationID", organizationID.Int())
	// 4. create owner-group
	ownerGroupID, err := createOwnerGroup(ctx, operator, organizationID)
	if err != nil {
		return nil, fmt.Errorf("CreateOwnerGroup: %w", err)
	}

	// 5. attach policy to "owner" group
	rbacOwnerGroup := authdomain.NewRBACRoleFromGroup(organizationID, ownerGroupID)
	// rbacAllUserRolesObject := authdomain.NewRBACAllUserRolesObjectFromOrganization(organizationID)

	for _, aoe := range []libdomain.ActionObjectEffect{
		{ // "owner" group "can" "CreateUser" "*"
			Action: authservice.CreateUserAction,
			Object: authservice.AnyObject,
			Effect: authservice.RBACAllowEffect,
		},
		// { // "owner" group "can" "unset" "all-user-roles"
		// 	Action: authservice.RBACUnsetAction,
		// 	Object: rbacAllUserRolesObject,
		// 	Effect: authservice.RBACAllowEffect,
		// },
	} {
		if err := attachPolicyToUserBySystemOwner(ctx, operator, rbacOwnerGroup, aoe.Action, aoe.Object, aoe.Effect); err != nil {
			return nil, fmt.Errorf("AttachPolicyToUserBySystemOwner: %w", err)
		}
	}
	return ownerGroupID, nil
}
