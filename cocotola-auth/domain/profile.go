package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

type ProfileModel struct {
	LoginID          string          `validate:"required"`
	Username         string          `validate:"required"`
	OrganizationID   *OrganizationID `validate:"required"`
	OrganizationName string          `validate:"required"`
	PersonalSpaceID  *SpaceID        // Guest user does not have personal space
}

func NewProfileModel(loginID string, username string, organizationID *OrganizationID, organizationName string, personalSpaceID *SpaceID) (*ProfileModel, error) {
	m := &ProfileModel{
		LoginID:          loginID,
		Username:         username,
		OrganizationID:   organizationID,
		OrganizationName: organizationName,
		PersonalSpaceID:  personalSpaceID,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, fmt.Errorf("validate profile model: %w", err)
	}

	return m, nil
}
