package service

import (
	"context"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
)

type AuthorizationManager interface {
	// Init(ctx context.Context) error

	AddUserToGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	AddUserToGroupBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	// RemoveUserFromGroup()

	// AddGroupToGroup(ctx context.Context, operator domain.User, src domain.UserGroupID, dst domain.UserGroupID) error
	AddObjectToObject(ctx context.Context, operator domain.SystemOwnerInterface, child, parent domain.RBACObject) error

	// RemoveGroupFromGroup()

	// AddObjectToObject()

	// RemoveObjectFromObject()

	AttachPolicyToUser(ctx context.Context, operator domain.UserInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AttachPolicyToUserBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AttachPolicyToUserBySystemOwner(ctx context.Context, operator domain.SystemOwnerInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error
	AttachPolicyToGroup(ctx context.Context, operator domain.UserInterface, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	AttachPolicyToGroupBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	// AddPolicyToGroup()

	// RemovePolicyToGroup()

	CheckAuthorization(ctx context.Context, operator domain.UserInterface, rbacAction domain.RBACAction, rbacObject domain.RBACObject) (bool, error)
}
