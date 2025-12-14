package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

type OrganizationID struct {
	Value int `validate:"required,gte=1"`
}

func NewOrganizationID(value int) (*OrganizationID, error) {
	m := OrganizationID{
		Value: value,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate organization ID: %w", err)
	}

	return &m, nil
}

func (v *OrganizationID) Int() int {
	return v.Value
}
func (v *OrganizationID) IsOrganizationID() bool {
	return true
}

type Organization struct {
	*libdomain.BaseModel
	OrganizationID *OrganizationID `validate:"required"`
	Name           string          `validate:"required,max=20"`
}

func NewOrganization(basemodel *libdomain.BaseModel, organizationID *OrganizationID, name string) (*Organization, error) {
	m := &Organization{
		BaseModel:      basemodel,
		OrganizationID: organizationID,
		Name:           name,
	}
	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, fmt.Errorf("validate organization: %w", err)
	}

	return m, nil
}
