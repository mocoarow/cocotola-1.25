package service

import (
	"context"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

// var RBACSetAction = domain.NewRBACAction("Set")
// var RBACUnsetAction = domain.NewRBACAction("Unset")

var AnyObject = domain.NewRBACObject("*")

var RBACAllowEffect = domain.NewRBACEffect("allow")
var RBACDenyEffect = domain.NewRBACEffect("deny")

type RBACEnforcer interface {
	LoadPolicy() error
	Enforce(rvals ...any) (bool, error)
}

type RBACRepository interface {
	GetEnforcer() RBACEnforcer
	// who can do what actions on which resources
	CreatePolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	// add user(or group) to parent group
	CreateSubjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACSubject, parent domain.RBACSubject) error

	// add child object to parent object
	CreateObjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error

	DeletePolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error

	DeleteSubjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACSubject, parent domain.RBACSubject) error
	DeleteObjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error

	NewEnforcerWithGroupsAndUsers(ctx context.Context, roles []domain.RBACRole, users []domain.RBACUser) (RBACEnforcer, error)

	// retrieve all groups (including inherited ones) a subject belongs to within a domain
	GetGroupsForSubject(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject) ([]domain.RBACRole, error)
}
