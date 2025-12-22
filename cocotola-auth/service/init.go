package service

import (
	"errors"
	"strconv"
)

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

func NewPersonalSpaceKey(userID int) string {
	return "__personal_space@@" + strconv.Itoa(userID)
}
func NewPersonalSpaceName(loginID string) string {
	return "Personal Space(" + loginID + ")"
}
