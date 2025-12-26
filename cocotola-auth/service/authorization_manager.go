package service

import (
	"context"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

type AuthorizationManager interface {
	// Init(ctx context.Context) error

	AddUserToGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	AddUserToGroupBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	// RemoveUserFromGroup()

	// AddGroupToGroup(ctx context.Context, operator domain.User, src domain.UserGroupID, dst domain.UserGroupID) error
	AddObjectToObject(ctx context.Context, operator domain.SystemOwnerInterface, child, parent libdomain.RBACObject) error

	// RemoveGroupFromGroup()

	// AddObjectToObject()

	// RemoveObjectFromObject()

	AttachPolicyToUser(ctx context.Context, operator domain.UserInterface, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error

	AttachPolicyToUserBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error

	AttachPolicyToUserBySystemOwner(ctx context.Context, operator domain.SystemOwnerInterface, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error
	AttachPolicyToGroup(ctx context.Context, operator domain.UserInterface, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error

	AttachPolicyToGroupBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error

	// AddPolicyToGroup()

	// RemovePolicyToGroup()

	CheckAuthorization(ctx context.Context, operator domain.UserInterface, rbacAction libdomain.RBACAction, rbacObject libdomain.RBACObject) (bool, error)
}
