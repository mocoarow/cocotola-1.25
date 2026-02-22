package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase/auth"
)

func TestVerifyPasswordCommand_Execute_shouldReturnNil_whenPasswordMatches(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newSystemOwner(t)

	gwMock := NewMockAuthVerifyPasswordCommandGateway(t)
	gwMock.EXPECT().VerifyPassword(ctx, operator, "login-id", "secret").Return(true, nil)

	cmd := auth.NewAuthVerifyPasswordCommand(gwMock)

	// when
	err := cmd.Execute(ctx, operator, "login-id", "secret")

	// then
	require.NoError(t, err)
}

func TestVerifyPasswordCommand_Execute_shouldReturnUnauthenticated_whenPasswordDoesNotMatch(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newSystemOwner(t)

	gwMock := NewMockAuthVerifyPasswordCommandGateway(t)
	gwMock.EXPECT().VerifyPassword(ctx, operator, "login-id", "secret").Return(false, nil)

	cmd := auth.NewAuthVerifyPasswordCommand(gwMock)

	// when
	err := cmd.Execute(ctx, operator, "login-id", "secret")

	// then
	require.ErrorIs(t, err, service.ErrUnauthenticated)
}

func TestVerifyPasswordCommand_Execute_shouldReturnError_whenRepositoryFails(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	operator := newSystemOwner(t)

	gwMock := NewMockAuthVerifyPasswordCommandGateway(t)
	gwMock.EXPECT().VerifyPassword(ctx, operator, "login-id", "secret").Return(false, errors.New("verify failure"))

	cmd := auth.NewAuthVerifyPasswordCommand(gwMock)

	// when
	err := cmd.Execute(ctx, operator, "login-id", "secret")

	// then
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
