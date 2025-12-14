//go:build small

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

func TestNewOrganizationID_shouldReturnOrganizationID_whenValidValueIsSpecified(t *testing.T) {
	t.Parallel()

	value := 123
	organizationID, err := domain.NewOrganizationID(value)

	require.NoError(t, err)
	assert.Equal(t, value, organizationID.Value)
	assert.Equal(t, value, organizationID.Int())
	assert.True(t, organizationID.IsOrganizationID())
}

func TestNewOrganizationID_shouldReturnOrganizationID_whenMinimumValueIsSpecified(t *testing.T) {
	t.Parallel()

	value := 1
	organizationID, err := domain.NewOrganizationID(value)

	require.NoError(t, err)
	assert.Equal(t, value, organizationID.Value)
	assert.Equal(t, value, organizationID.Int())
	assert.True(t, organizationID.IsOrganizationID())
}

func TestNewOrganization_shouldReturnOrganization_whenValidParametersAreSpecified(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	organizationID, _ := domain.NewOrganizationID(123)
	name := "Test Organization"

	organization, err := domain.NewOrganization(baseModel, organizationID, name)

	require.NoError(t, err)
	assert.Equal(t, baseModel, organization.BaseModel)
	assert.Equal(t, organizationID, organization.OrganizationID)
	assert.Equal(t, name, organization.Name)
}

func TestNewOrganization_shouldReturnError_whenOrganizationIDIsNil(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	name := "Test Organization"

	organization, err := domain.NewOrganization(baseModel, nil, name)

	require.Error(t, err)
	assert.Nil(t, organization)
	assert.Contains(t, err.Error(), "validate organization")
}

func TestNewOrganization_shouldReturnError_whenNameIsEmpty(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	organizationID, _ := domain.NewOrganizationID(123)
	name := ""

	organization, err := domain.NewOrganization(baseModel, organizationID, name)

	require.Error(t, err)
	assert.Nil(t, organization)
	assert.Contains(t, err.Error(), "validate organization")
}

func TestNewOrganization_shouldReturnError_whenNameExceedsMaxLength(t *testing.T) {
	t.Parallel()

	baseModel := &libdomain.BaseModel{Version: 1}
	organizationID, _ := domain.NewOrganizationID(123)
	name := "123456789012345678901"

	organization, err := domain.NewOrganization(baseModel, organizationID, name)

	require.Error(t, err)
	assert.Nil(t, organization)
	assert.Contains(t, err.Error(), "validate organization")
}
