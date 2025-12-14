package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

type UserID struct {
	Value int `validate:"gte=0"`
}

func NewUserID(value int) (*UserID, error) {
	m := UserID{
		Value: value,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate user ID: %w", err)
	}

	return &m, nil
}

func (v *UserID) Int() int {
	return v.Value
}
func (v *UserID) IsUserID() bool {
	return true
}
func (v *UserID) GetRBACSubject() RBACSubject {
	return NewRBACUserFromUser(v)
}

type User struct {
	*libdomain.BaseModel `validate:"required"`
	UserID               *UserID         `validate:"required"`
	OrganizationID       *OrganizationID `validate:"required"`
	LoginID              string          `validate:"required"`
	Username             string          `validate:"required"`
	UserGroups           []*UserGroup
}

func NewUser(baseModel *libdomain.BaseModel, userID *UserID, organizationID *OrganizationID, loginID, username string, userGroups []*UserGroup) (*User, error) {
	m := User{
		BaseModel:      baseModel,
		UserID:         userID,
		OrganizationID: organizationID,
		LoginID:        loginID,
		Username:       username,
		UserGroups:     userGroups,
	}

	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate user: %w", err)
	}

	return &m, nil
}

func (m *User) GetUserID() *UserID {
	return m.UserID
}
func (m *User) GetOrganizationID() *OrganizationID {
	return m.OrganizationID
}
