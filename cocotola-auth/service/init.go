package service

import "errors"

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

const (
	SystemAdminLoginID = "__system_admin"
	SystemOwnerLoginID = "__system_owner"

	SystemOwnerGroupKey   = "__system_owner"
	OwnerGroupKey         = "__owner"
	PublicGroupKey        = "__public_group"
	PublicDefaultSpaceKey = "__public_default_space"

	SystemOwnerGroupName   = "System Owner"
	OwnerGroupName         = "Owner"
	PublicGroupName        = "Public Group"
	PublicDefaultSpaceName = "Public Default Space"
)
