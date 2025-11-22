//go:build medium

package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

func TestAuthorizationManager_CheckAuthorization_shouldReflectGroupMembership_whenUserJoinsGroup(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, sysOwner, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		authorizationManager, err := gateway.NewAuthorizationManager(ctx, tr.dialect, tr.db, tr.rf)
		require.NoError(t, err)
		membershipRepo := gateway.NewPairOfUserAndGroupRepository(ctx, tr.dialect, tr.db, tr.rf)

		group := testAddUserGroup(t, ctx, tr, owner, "GROUP_KEY_Auth", "GROUP_NAME_Auth", "GROUP_DESC_Auth")
		rbacGroup := domain.NewRBACRoleFromGroup(orgID, group.UserGroupID)
		rbacObject := domain.NewRBACAllUserRolesObjectFromOrganization(orgID)

		require.NoError(t, authorizationManager.AddPolicyToGroup(ctx, owner, rbacGroup, service.RBACSetAction, rbacObject, service.RBACAllowEffect))

		targetUser := testAddUser(t, ctx, tr, owner, "LOGIN_ID_TARGET", "USERNAME_TARGET", "PASSWORD_TARGET")
		otherUser := testAddUser(t, ctx, tr, owner, "LOGIN_ID_OTHER", "USERNAME_OTHER", "PASSWORD_OTHER")

		ok, err := authorizationManager.CheckAuthorization(ctx, targetUser, service.RBACSetAction, rbacObject)
		require.NoError(t, err)
		assert.False(t, ok, "user without group membership must not have permission")

		ok, err = authorizationManager.CheckAuthorization(ctx, otherUser, service.RBACSetAction, rbacObject)
		require.NoError(t, err)
		assert.False(t, ok)

		groupsBefore, err := membershipRepo.FindUserGroupsByUserID(ctx, targetUser, targetUser.GetUserID())
		require.NoError(t, err)
		assert.Empty(t, groupsBefore)

		require.NoError(t, authorizationManager.AddUserToGroup(ctx, owner, targetUser.GetUserID(), group.UserGroupID))

		groupsAfter, err := membershipRepo.FindUserGroupsByUserID(ctx, targetUser, targetUser.GetUserID())
		require.NoError(t, err)
		assert.Len(t, groupsAfter, 1)

		rbacRepo, err := gateway.NewRBACRepository(ctx, tr.db)
		require.NoError(t, err)
		e := rbacRepo.GetEnforcer()
		require.NoError(t, e.LoadPolicy())
		manual, err := e.Enforce(domain.NewRBACUserFromUser(targetUser.GetUserID()).Subject(), rbacObject.Object(), service.RBACSetAction.Action(), domain.NewRBACDomainFromOrganization(orgID).Domain())
		require.NoError(t, err)
		assert.True(t, manual)

		ok, err = authorizationManager.CheckAuthorization(ctx, targetUser, service.RBACSetAction, rbacObject)
		require.NoError(t, err)
		assert.True(t, ok, "user should gain permission after joining group")

		ok, err = authorizationManager.CheckAuthorization(ctx, otherUser, service.RBACSetAction, rbacObject)
		require.NoError(t, err)
		assert.False(t, ok, "unrelated user must remain unauthorized")

		_ = sysOwner
	}

	testDB(t, fn)
}
