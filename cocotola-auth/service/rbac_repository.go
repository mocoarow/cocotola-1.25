package service

import (
	"context"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

// var RBACSetAction = libdomain.NewRBACAction("Set")
// var RBACUnsetAction = libdomain.NewRBACAction("Unset")

var AnyObject = libdomain.NewRBACObject("*")

var RBACAllowEffect = libdomain.NewRBACEffect("allow")
var RBACDenyEffect = libdomain.NewRBACEffect("deny")

type RBACEnforcer interface {
	LoadPolicy() error
	Enforce(rvals ...any) (bool, error)
}

type RBACRepository interface {
	GetEnforcer() RBACEnforcer
	// who can do what actions on which resources
	CreatePolicy(ctx context.Context, domain libdomain.RBACDomain, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error

	// add user(or group) to parent group
	CreateSubjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomain, child libdomain.RBACSubject, parent libdomain.RBACSubject) error

	// add child object to parent object
	CreateObjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomain, child libdomain.RBACObject, parent libdomain.RBACObject) error

	DeletePolicy(ctx context.Context, domain libdomain.RBACDomain, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error

	DeleteSubjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomain, child libdomain.RBACSubject, parent libdomain.RBACSubject) error
	DeleteObjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomain, child libdomain.RBACObject, parent libdomain.RBACObject) error

	NewEnforcerWithGroupsAndUsers(ctx context.Context, roles []libdomain.RBACRole, users []libdomain.RBACUser) (RBACEnforcer, error)

	// retrieve all groups (including inherited ones) a subject belongs to within a domain
	GetGroupsForSubject(ctx context.Context, domain libdomain.RBACDomain, subject libdomain.RBACSubject) ([]libdomain.RBACRole, error)
}
