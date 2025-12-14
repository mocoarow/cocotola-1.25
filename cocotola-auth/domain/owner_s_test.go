//go:build small

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

func TestNewOwner_shouldReturnOwner_whenValidUserIsSpecified(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	userID, _ := domain.NewUserID(123)
	organizationID, _ := domain.NewOrganizationID(456)
	loginID := "owner@example.com"
	username := "owneruser"
	userGroups := []*domain.UserGroup{}

	user, _ := domain.NewUser(baseModel, userID, organizationID, loginID, username, userGroups)

	owner, err := domain.NewOwner(user)

	require.NoError(t, err)
	assert.NotNil(t, owner.User)
	assert.Equal(t, user, owner.User)
	assert.True(t, owner.IsOwner())
	assert.Equal(t, userID, owner.GetUserID())
	assert.Equal(t, organizationID, owner.GetOrganizationID())
}

func TestNewOwner_shouldReturnError_whenUserIsNil(t *testing.T) {
	t.Parallel()

	owner, err := domain.NewOwner(nil)

	require.Error(t, err)
	assert.Nil(t, owner)
	assert.Contains(t, err.Error(), "new owner")
}

func TestOwner_IsOwner_shouldReturnTrue_whenCalled(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	userID, _ := domain.NewUserID(123)
	organizationID, _ := domain.NewOrganizationID(456)
	loginID := "owner@example.com"
	username := "owneruser"
	userGroups := []*domain.UserGroup{}

	user, _ := domain.NewUser(baseModel, userID, organizationID, loginID, username, userGroups)
	owner, _ := domain.NewOwner(user)

	result := owner.IsOwner()

	assert.True(t, result)
}

func TestOwner_GetUserID_shouldReturnUserID_whenCalled(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	userID, _ := domain.NewUserID(123)
	organizationID, _ := domain.NewOrganizationID(456)
	loginID := "owner@example.com"
	username := "owneruser"
	userGroups := []*domain.UserGroup{}

	user, _ := domain.NewUser(baseModel, userID, organizationID, loginID, username, userGroups)
	owner, _ := domain.NewOwner(user)

	result := owner.GetUserID()

	require.NotNil(t, result)
	assert.Equal(t, userID, result)
	assert.Equal(t, 123, result.Value)
}

func TestOwner_GetOrganizationID_shouldReturnOrganizationID_whenCalled(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	userID, _ := domain.NewUserID(123)
	organizationID, _ := domain.NewOrganizationID(456)
	loginID := "owner@example.com"
	username := "owneruser"
	userGroups := []*domain.UserGroup{}

	user, _ := domain.NewUser(baseModel, userID, organizationID, loginID, username, userGroups)
	owner, _ := domain.NewOwner(user)

	result := owner.GetOrganizationID()

	require.NotNil(t, result)
	assert.Equal(t, organizationID, result)
	assert.Equal(t, 456, result.Value)
}
