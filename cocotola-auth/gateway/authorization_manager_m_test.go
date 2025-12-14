//go:build medium

package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

func TestAuthorizationManager_CheckAuthorization_shouldReflectGroupMembership_whenUserJoinsGroup(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		// given
		// organization
		//   <ORGANIZATION> with random name
		// groups
		//   TEST_GROUP
		// users
		//   TARGET_USER
		//   OTHER_USER
		orgID, _, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		authorizationManager, err := gateway.NewAuthorizationManager(ctx, tr.dialect, tr.db, tr.rf)
		require.NoError(t, err)
		membershipRepo := gateway.NewPairOfUserAndGroupRepository(ctx, tr.dialect, tr.db, tr.rf)

		group := testAddUserGroup(t, ctx, tr, owner, "TEST_GROUP", "GROUP_NAME_Auth", "GROUP_DESC_Auth")
		rbacGroup := domain.NewRBACRoleFromGroup(orgID, group.UserGroupID)
		if err := authorizationManager.AttachPolicyToGroup(ctx, owner, rbacGroup, service.CreateUserAction, service.AnyObject, service.RBACAllowEffect); err != nil {
			require.NoError(t, err)
		}

		targetUser := testAddUser(t, ctx, tr, owner, "TARGET_USER", "USERNAME_TARGET", "PASSWORD_TARGET")
		otherUser := testAddUser(t, ctx, tr, owner, "OTHER_USER", "USERNAME_OTHER", "PASSWORD_OTHER")

		ok, err := authorizationManager.CheckAuthorization(ctx, targetUser, service.CreateUserAction, service.AnyObject)
		require.NoError(t, err)
		assert.False(t, ok, "user without group membership must not have permission")

		ok, err = authorizationManager.CheckAuthorization(ctx, otherUser, service.CreateUserAction, service.AnyObject)
		require.NoError(t, err)
		assert.False(t, ok)

		groupsBefore, err := membershipRepo.FindUserGroupsByUserID(ctx, targetUser, targetUser.GetUserID())
		require.NoError(t, err)
		assert.Empty(t, groupsBefore)

		// when
		require.NoError(t, authorizationManager.AddUserToGroup(ctx, owner, targetUser.GetUserID(), group.UserGroupID))

		groupsAfter, err := membershipRepo.FindUserGroupsByUserID(ctx, targetUser, targetUser.GetUserID())
		require.NoError(t, err)
		assert.Len(t, groupsAfter, 1)

		rbacRepo, err := gateway.NewRBACRepository(ctx, tr.db)
		require.NoError(t, err)
		e := rbacRepo.GetEnforcer()
		require.NoError(t, e.LoadPolicy())

		// then
		// TARGET_USER should gain permission via group membership
		manual, err := e.Enforce(domain.NewRBACUserFromUser(targetUser.GetUserID()).Subject(), service.AnyObject.Object(), service.CreateUserAction.Action(), domain.NewRBACDomainFromOrganization(orgID).Domain())
		require.NoError(t, err)
		assert.True(t, manual)

		// TARGET_USER should gain permission via group membership
		ok, err = authorizationManager.CheckAuthorization(ctx, targetUser, service.CreateUserAction, service.AnyObject)
		require.NoError(t, err)
		assert.True(t, ok, "user should gain permission after joining group")

		// OTHER_USER should remain unauthorized
		ok, err = authorizationManager.CheckAuthorization(ctx, otherUser, service.CreateUserAction, service.AnyObject)
		require.NoError(t, err)
		assert.False(t, ok, "unrelated user must remain unauthorized")
	}

	testDB(t, fn)
}
