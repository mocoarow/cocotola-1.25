package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
	serviceMocks "github.com/mocoarow/cocotola-1.25/moonbeam/user/service/mocks"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/usecase"
)

func TestVerifyPasswordCommand_Execute_shouldReturnNil_whenPasswordMatches(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	operator := newSystemOwner(t)

	txMock := serviceMocks.NewMockTransactionManager(t)
	rfMock := serviceMocks.NewMockRepositoryFactory(t)
	userRepoMock := serviceMocks.NewMockUserRepository(t)

	txMock.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, fn func(service.RepositoryFactory) error) error {
			return fn(rfMock)
		},
	)
	rfMock.EXPECT().NewUserRepository(mock.Anything).Return(userRepoMock)
	userRepoMock.EXPECT().VerifyPassword(mock.Anything, operator, "login-id", "secret").Return(true, nil)

	cmd := usecase.NewVerifyPasswordCommand(txMock)

	err := cmd.Execute(ctx, operator, "login-id", "secret")
	require.NoError(t, err)
}

func TestVerifyPasswordCommand_Execute_shouldReturnUnauthenticated_whenPasswordDoesNotMatch(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	operator := newSystemOwner(t)

	txMock := serviceMocks.NewMockTransactionManager(t)
	rfMock := serviceMocks.NewMockRepositoryFactory(t)
	userRepoMock := serviceMocks.NewMockUserRepository(t)

	txMock.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, fn func(service.RepositoryFactory) error) error {
			return fn(rfMock)
		},
	)
	rfMock.EXPECT().NewUserRepository(mock.Anything).Return(userRepoMock)
	userRepoMock.EXPECT().VerifyPassword(mock.Anything, operator, "login-id", "secret").Return(false, nil)

	cmd := usecase.NewVerifyPasswordCommand(txMock)

	err := cmd.Execute(ctx, operator, "login-id", "secret")
	require.ErrorIs(t, err, service.ErrUnauthenticated)
}

func TestVerifyPasswordCommand_Execute_shouldReturnError_whenRepositoryFails(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	operator := newSystemOwner(t)

	txMock := serviceMocks.NewMockTransactionManager(t)
	rfMock := serviceMocks.NewMockRepositoryFactory(t)
	userRepoMock := serviceMocks.NewMockUserRepository(t)

	txMock.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, fn func(service.RepositoryFactory) error) error {
			return fn(rfMock)
		},
	)
	rfMock.EXPECT().NewUserRepository(mock.Anything).Return(userRepoMock)
	userRepoMock.EXPECT().VerifyPassword(mock.Anything, operator, "login-id", "secret").Return(false, errors.New("verify failure"))

	cmd := usecase.NewVerifyPasswordCommand(txMock)

	err := cmd.Execute(ctx, operator, "login-id", "secret")
	require.ErrorContains(t, err, "VerifyPassword")
}

func newSystemOwner(t *testing.T) domain.SystemOwnerInterface {
	t.Helper()

	orgID, err := domain.NewOrganizationID(1)
	require.NoError(t, err)
	userID, err := domain.NewUserID(10)
	require.NoError(t, err)
	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)

	user, err := domain.NewUser(baseModel, userID, orgID, "login", "username", nil)
	require.NoError(t, err)
	owner, err := domain.NewOwner(user)
	require.NoError(t, err)
	systemOwner, err := domain.NewSystemOwner(owner)
	require.NoError(t, err)

	return systemOwner
}
