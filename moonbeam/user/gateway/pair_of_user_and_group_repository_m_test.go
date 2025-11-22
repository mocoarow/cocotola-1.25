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

func TestPairOfUserAndGroupRepository_FindUserGroupsByUserID_shouldReturnGroups_whenUserBelongsToMultipleGroups(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, _, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		repo := gateway.NewPairOfUserAndGroupRepository(ctx, tr.dialect, tr.db, tr.rf)

		user1 := testAddUser(t, ctx, tr, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		user2 := testAddUser(t, ctx, tr, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")

		group1 := testAddUserGroup(t, ctx, tr, owner, "GROUP_KEY_1", "GROUP_NAME_1", "GROUP_DESC_1")
		group2 := testAddUserGroup(t, ctx, tr, owner, "GROUP_KEY_2", "GROUP_NAME_2", "GROUP_DESC_2")
		group3 := testAddUserGroup(t, ctx, tr, owner, "GROUP_KEY_3", "GROUP_NAME_3", "GROUP_DESC_3")

		for _, g := range []*domain.UserGroup{group1, group2, group3} {
			require.NoError(t, repo.CreatePairOfUserAndGroup(ctx, owner, user1.GetUserID(), g.UserGroupID))
		}
		require.NoError(t, repo.CreatePairOfUserAndGroup(ctx, owner, user2.GetUserID(), group1.UserGroupID))

		result := tr.db.WithContext(ctx).
			Table("casbin_rule").
			Where("ptype = ? AND v0 = ?", "g", fmt.Sprintf("user:%d", user1.GetUserID().Int())).
			Find(&[]struct{}{})
		require.NoError(t, result.Error)
		assert.Equal(t, int64(3), result.RowsAffected)

		groups1, err := repo.FindUserGroupsByUserID(ctx, owner, user1.GetUserID())
		require.NoError(t, err)
		groups2, err := repo.FindUserGroupsByUserID(ctx, owner, user2.GetUserID())
		require.NoError(t, err)

		assert.Len(t, groups1, 3)
		assert.ElementsMatch(t, []string{"GROUP_KEY_1", "GROUP_KEY_2", "GROUP_KEY_3"}, extractGroupKeys(groups1))
		assert.Len(t, groups2, 1)
		assert.Equal(t, "GROUP_KEY_1", groups2[0].Key)
	}

	testDB(t, fn)
}

func TestPairOfUserAndGroupRepository_DeletePairOfUserAndGroup_shouldRemoveRelation_whenPairExists(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, _, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		repo := gateway.NewPairOfUserAndGroupRepository(ctx, tr.dialect, tr.db, tr.rf)

		user := testAddUser(t, ctx, tr, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		group := testAddUserGroup(t, ctx, tr, owner, "GROUP_KEY_1", "GROUP_NAME_1", "GROUP_DESC_1")

		require.NoError(t, repo.CreatePairOfUserAndGroup(ctx, owner, user.GetUserID(), group.UserGroupID))

		before, err := repo.FindUserGroupsByUserID(ctx, owner, user.GetUserID())
		require.NoError(t, err)
		assert.Len(t, before, 1)

		require.NoError(t, repo.DeletePairOfUserAndGroup(ctx, owner, user.GetUserID(), group.UserGroupID))

		after, err := repo.FindUserGroupsByUserID(ctx, owner, user.GetUserID())
		require.NoError(t, err)
		assert.Len(t, after, 0)

		err = repo.DeletePairOfUserAndGroup(ctx, owner, user.GetUserID(), group.UserGroupID)
		assert.ErrorIs(t, err, service.ErrPairOfUserAndGroupNotFound)
	}

	testDB(t, fn)
}

func extractGroupKeys(groups []*domain.UserGroup) []string {
	keys := make([]string, len(groups))
	for i, g := range groups {
		keys[i] = g.Key
	}
	return keys
}
