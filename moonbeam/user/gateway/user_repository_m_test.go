//go:build medium

package gateway_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

func TestUserRepository_CreateAndFindUser_shouldReturnUser_whenOwnerCreates(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, _, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		userRepo := tr.rf.NewUserRepository(ctx)
		param := testNewCreateUserParameter(t, fmt.Sprintf("login_%s", RandString(4)), "USERNAME_U", "PASSWORD_U")
		userID, err := userRepo.CreateUser(ctx, owner, param)
		require.NoError(t, err)

		user, err := userRepo.FindUserByID(ctx, owner, userID)
		require.NoError(t, err)
		assert.Equal(t, param.LoginID, user.LoginID)
		assert.Equal(t, param.Username, user.Username)

		userByLogin, err := userRepo.FindUserByLoginID(ctx, owner, param.LoginID)
		require.NoError(t, err)
		assert.Equal(t, user.GetUserID().Int(), userByLogin.GetUserID().Int())

		sysAdmin := domain.NewSystemAdmin()
		sysOwnerByID, err := userRepo.FindSystemOwnerByOrganizationID(ctx, sysAdmin, orgID)
		require.NoError(t, err)
		assert.Equal(t, service.SystemOwnerLoginID, sysOwnerByID.LoginID)

		orgRepo := gateway.NewOrganizationRepository(ctx, tr.db)
		org, err := orgRepo.GetOrganization(ctx, owner)
		require.NoError(t, err)
		sysOwnerByName, err := userRepo.FindSystemOwnerByOrganizationName(ctx, sysAdmin, org.Name)
		require.NoError(t, err)
		assert.Equal(t, sysOwnerByID.GetUserID().Int(), sysOwnerByName.GetUserID().Int())
	}

	testDB(t, fn)
}
