package domain

import "errors"

type Owner struct {
	*User
}

func NewOwner(user *User) (*Owner, error) {
	if user == nil {
		return nil, errors.New("new owner: user is nil")
	}

	return &Owner{
		User: user,
	}, nil
}

func (m *Owner) IsOwner() bool {
	return true
}
func (m *Owner) GetOrganizationID() *OrganizationID {
	return m.User.GetOrganizationID()
}
func (m *Owner) GetUserID() *UserID {
	return m.User.GetUserID()
}
