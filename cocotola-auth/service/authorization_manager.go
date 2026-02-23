package service

import (
	"context"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

type AttachPolicyToUserBySystemAdminFunc func(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject libdomain.RBACSubjectInterface, action libdomain.RBACActionInterface, object libdomain.RBACObjectInterface, effect libdomain.RBACEffectInterface) error

type AttachPolicyToUserBySystemOwnerFunc func(ctx context.Context, operator domain.SystemOwnerInterface, subject libdomain.RBACSubjectInterface, action libdomain.RBACActionInterface, object libdomain.RBACObjectInterface, effect libdomain.RBACEffectInterface) error

type AddUserToGroupFunc func(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error

type AttachPolicyToUserFunc func(ctx context.Context, operator domain.UserInterface, subject libdomain.RBACSubjectInterface, action libdomain.RBACActionInterface, object libdomain.RBACObjectInterface, effect libdomain.RBACEffectInterface) error
