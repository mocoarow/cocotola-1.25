package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type SystemOwner struct {
	*Owner `validate:"required"`
}

func NewSystemOwner(user *Owner) (*SystemOwner, error) {
	m := SystemOwner{
		Owner: user,
	}

	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate system owner: %w", err)
	}

	return &m, nil
}

func (m *SystemOwner) IsOwner() bool {
	return true
}
func (m *SystemOwner) IsSystemOwner() bool {
	return true
}
func (m *SystemOwner) GetOrganizationID() *OrganizationID {
	return m.Owner.GetOrganizationID()
}
func (m *SystemOwner) GetUserID() *UserID {
	return m.Owner.GetUserID()
}
