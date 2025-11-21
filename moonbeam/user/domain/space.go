package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type SpaceID struct {
	Value int `validate:"required,gte=1"`
}

func NewSpaceID(value int) (*SpaceID, error) {
	m := SpaceID{
		Value: value,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate space id(%d): %w", value, err)
	}
	return &m, nil
}

func (v *SpaceID) Int() int {
	return v.Value
}
func (v *SpaceID) IsSpaceID() bool {
	return true
}

type SpaceIDs []*SpaceID

func (v *SpaceIDs) IDs() []int {
	if v == nil {
		return nil
	}

	ids := make([]int, len(*v))
	for i, id := range *v {
		ids[i] = id.Int()
	}

	return ids
}

type Space struct {
	*libdomain.BaseModel
	SpaceID        *SpaceID        `validate:"required"`
	OrganizationID *OrganizationID `validate:"required"`
	OwnerID        *UserID         `validate:"required"`
	KeyName        string          `validate:"required"`
	Name           string          `validate:"required"`
	SpaceType      string          `validate:"required,oneof=personal private public"`
}

func NewSpace(baseModel *libdomain.BaseModel, spaceID *SpaceID, organizationID *OrganizationID, owernID *UserID, keyName, name string, spaceType string) (*Space, error) {
	m := Space{
		BaseModel:      baseModel,
		SpaceID:        spaceID,
		OrganizationID: organizationID,
		OwnerID:        owernID,
		KeyName:        keyName,
		Name:           name,
		SpaceType:      spaceType,
	}

	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate space model: %w", err)
	}

	return &m, nil
}
func (m *Space) IsPrivate() bool {
	return m.SpaceType == "private"
}
