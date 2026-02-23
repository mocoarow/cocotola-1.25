package profile_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase/profile"
)

func Test_GetMyProfileQuery_Execute_shouldReturnProfile_whenAllDataFound(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newUser(t)
	org := newOrganization(t)
	user := newUserModel(t)
	space := newSpace(t)

	repoMock := NewMockGetMyProfileQueryRepository(t)
	repoMock.EXPECT().GetOrganization(mock.Anything, operator).Return(org, nil)
	repoMock.EXPECT().GetUser(mock.Anything, operator).Return(user, nil)
	repoMock.EXPECT().GetPersonalSpace(mock.Anything, operator).Return(space, nil)

	query := profile.NewGetMyProfileQuery(repoMock)

	// when
	result, err := query.Execute(ctx, operator)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, user.LoginID, result.LoginID)
	require.Equal(t, user.Username, result.Username)
	require.Equal(t, org.OrganizationID, result.OrganizationID)
	require.Equal(t, org.Name, result.OrganizationName)
	require.Equal(t, space.SpaceID, result.PersonalSpaceID)
}

func Test_GetMyProfileQuery_Execute_shouldReturnProfileWithNilSpaceID_whenPersonalSpaceNotFound(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newUser(t)
	org := newOrganization(t)
	user := newUserModel(t)

	repoMock := NewMockGetMyProfileQueryRepository(t)
	repoMock.EXPECT().GetOrganization(mock.Anything, operator).Return(org, nil)
	repoMock.EXPECT().GetUser(mock.Anything, operator).Return(user, nil)
	repoMock.EXPECT().GetPersonalSpace(mock.Anything, operator).Return(nil, service.ErrSpaceNotFound)

	query := profile.NewGetMyProfileQuery(repoMock)

	// when
	result, err := query.Execute(ctx, operator)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Nil(t, result.PersonalSpaceID)
}

func Test_GetMyProfileQuery_Execute_shouldReturnError_whenGetOrganizationFails(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newUser(t)

	repoMock := NewMockGetMyProfileQueryRepository(t)
	repoMock.EXPECT().GetOrganization(mock.Anything, operator).Return(nil, errors.New("org error"))

	query := profile.NewGetMyProfileQuery(repoMock)

	// when
	result, err := query.Execute(ctx, operator)

	// then
	require.ErrorContains(t, err, "GetOrganization")
	require.Nil(t, result)
}

func Test_GetMyProfileQuery_Execute_shouldReturnError_whenGetUserFails(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newUser(t)
	org := newOrganization(t)

	repoMock := NewMockGetMyProfileQueryRepository(t)
	repoMock.EXPECT().GetOrganization(mock.Anything, operator).Return(org, nil)
	repoMock.EXPECT().GetUser(mock.Anything, operator).Return(nil, errors.New("user error"))

	query := profile.NewGetMyProfileQuery(repoMock)

	// when
	result, err := query.Execute(ctx, operator)

	// then
	require.ErrorContains(t, err, "GetUser")
	require.Nil(t, result)
}

func Test_GetMyProfileQuery_Execute_shouldReturnError_whenGetPersonalSpaceFailsWithUnexpectedError(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newUser(t)
	org := newOrganization(t)
	user := newUserModel(t)

	repoMock := NewMockGetMyProfileQueryRepository(t)
	repoMock.EXPECT().GetOrganization(mock.Anything, operator).Return(org, nil)
	repoMock.EXPECT().GetUser(mock.Anything, operator).Return(user, nil)
	repoMock.EXPECT().GetPersonalSpace(mock.Anything, operator).Return(nil, errors.New("unexpected error"))

	query := profile.NewGetMyProfileQuery(repoMock)

	// when
	result, err := query.Execute(ctx, operator)

	// then
	require.ErrorContains(t, err, "GetPersonalSpace")
	require.Nil(t, result)
}

func newUser(t *testing.T) *domain.User {
	t.Helper()

	orgID, err := domain.NewOrganizationID(1)
	require.NoError(t, err)
	userID, err := domain.NewUserID(10)
	require.NoError(t, err)
	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)

	user, err := domain.NewUser(baseModel, userID, orgID, "login-id", "username", nil)
	require.NoError(t, err)

	return user
}

func newUserModel(t *testing.T) *domain.User {
	t.Helper()

	orgID, err := domain.NewOrganizationID(1)
	require.NoError(t, err)
	userID, err := domain.NewUserID(10)
	require.NoError(t, err)
	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)

	user, err := domain.NewUser(baseModel, userID, orgID, "login-id", "username", nil)
	require.NoError(t, err)

	return user
}

func newOrganization(t *testing.T) *domain.Organization {
	t.Helper()

	orgID, err := domain.NewOrganizationID(1)
	require.NoError(t, err)
	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)

	org, err := domain.NewOrganization(baseModel, orgID, "test-org")
	require.NoError(t, err)

	return org
}

func newSpace(t *testing.T) *domain.Space {
	t.Helper()

	orgID, err := domain.NewOrganizationID(1)
	require.NoError(t, err)
	userID, err := domain.NewUserID(10)
	require.NoError(t, err)
	spaceID, err := domain.NewSpaceID(100)
	require.NoError(t, err)
	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)

	space, err := domain.NewSpace(baseModel, spaceID, orgID, userID, "personal-key", "Personal Space", "personal")
	require.NoError(t, err)

	return space
}
