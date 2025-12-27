package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

func NewRBACDomainFromOrganization(organizationID *OrganizationID) libdomain.RBACDomain {
	return libdomain.NewRBACDomain(fmt.Sprintf("domain:%d", organizationID.Int()))
}

func NewRBACUserFromUser(userID *UserID) libdomain.RBACUser {
	return libdomain.NewRBACUser(fmt.Sprintf("user:%d", userID.Int()))
}

func NewRBACRoleFromGroup(organizationID *OrganizationID, userGroupID *UserGroupID) libdomain.RBACRole {
	return libdomain.NewRBACRole(fmt.Sprintf("domain:%d,role:%d", organizationID.Int(), userGroupID.Int()))
}

func NewRBACRoleFromSpace(organizationID *OrganizationID, spaceID *SpaceID) libdomain.RBACRole {
	return libdomain.NewRBACRole(fmt.Sprintf("domain:%d,space:%d", organizationID.Int(), spaceID.Int()))
}

func NewRBACObjectFromGroup(organizationID *OrganizationID, userRoleID *UserGroupID) libdomain.RBACObject {
	return libdomain.NewRBACObject(fmt.Sprintf("domain:%d,role:%d", organizationID.Int(), userRoleID.Int()))
}

func NewRBACAllUserRolesObjectFromOrganization(organizationID *OrganizationID) libdomain.RBACObject {
	return libdomain.NewRBACObject(fmt.Sprintf("domain:%d,role:*", organizationID.Int()))
}

func NewOrganizationAndSpaceIDsFromRole(role libdomain.RBACRole) (*OrganizationID, *SpaceID, error) {
	var orgID, spaceID int
	if _, err := fmt.Sscanf(role.Role(), "domain:%d,space:%d", &orgID, &spaceID); err != nil {
		return nil, nil, fmt.Errorf("parse role(%s): %w", role.Role(), err)
	}

	org, err := NewOrganizationID(orgID)
	if err != nil {
		return nil, nil, fmt.Errorf("domain.NewOrganizationID: %w", err)
	}

	sp, err := NewSpaceID(spaceID)
	if err != nil {
		return nil, nil, fmt.Errorf("domain.NewSpaceID: %w", err)
	}

	return org, sp, nil
}

func NewOrganizationAndUserGroupIDsFromRole(role libdomain.RBACRole) (*OrganizationID, *UserGroupID, error) {
	var orgIDValue, groupIDValue int
	if _, err := fmt.Sscanf(role.Role(), "domain:%d,role:%d", &orgIDValue, &groupIDValue); err != nil {
		return nil, nil, fmt.Errorf("parse rbac role(%s): %w", role.Role(), err)
	}

	orgID, err := NewOrganizationID(orgIDValue)
	if err != nil {
		return nil, nil, fmt.Errorf("domain.NewOrganizationID: %w", err)
	}

	groupID, err := NewUserGroupID(groupIDValue)
	if err != nil {
		return nil, nil, fmt.Errorf("domain.NewUserGroupID: %w", err)
	}

	return orgID, groupID, nil
}
