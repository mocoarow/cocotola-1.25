package service

import (
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

// var RBACSetAction = libdomain.NewRBACAction("Set")
// var RBACUnsetAction = libdomain.NewRBACAction("Unset")

var AnyObject = libdomain.NewRBACObject("*") //nolint:gochecknoglobals

var RBACAllowEffect = libdomain.NewRBACEffect("allow") //nolint:gochecknoglobals
var RBACDenyEffect = libdomain.NewRBACEffect("deny")   //nolint:gochecknoglobals

type RBACEnforcer interface {
	LoadPolicy() error
	Enforce(rvals ...any) (bool, error)
}
