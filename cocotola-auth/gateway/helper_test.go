package gateway_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

func organizationID(t *testing.T, organizationID int) *domain.OrganizationID {
	t.Helper()
	id, err := domain.NewOrganizationID(organizationID)
	require.NoError(t, err)
	return id
}

func userID(t *testing.T, userID int) *domain.UserID {
	t.Helper()
	id, err := domain.NewUserID(userID)
	require.NoError(t, err)
	return id
}

type user struct {
	userID         *domain.UserID
	organizationID *domain.OrganizationID
	loginID        string
	username       string
}

var _ domain.UserInterface = (*user)(nil)

func (m *user) GetUserID() *domain.UserID {
	return m.userID
}
func (m *user) GetOrganizationID() *domain.OrganizationID {
	return m.organizationID
}
func (m *user) GetUsername() string {
	return m.username
}
func (m *user) GetLoginID() string {
	return m.loginID
}
