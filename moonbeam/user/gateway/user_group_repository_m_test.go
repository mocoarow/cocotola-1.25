//go:build medium

package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

func TestUserGroupRepository_AddUserGroup_shouldPersistAndBeQueryable(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, sysOwner, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		repo := tr.rf.NewUserGroupRepository(ctx)
		param := testNewUserGroupAddParameter(t, "group-key", "Group Name", "desc")
		groupID, err := repo.AddUserGroup(ctx, owner, param)
		require.NoError(t, err)

		group, err := repo.FindUserGroupByID(ctx, owner, groupID)
		require.NoError(t, err)
		assert.Equal(t, param.Key, group.Key)
		assert.Equal(t, param.Name, group.Name)
		assert.Equal(t, param.Description, group.Description)

		retrievedByKey, err := repo.FindUserGroupByKey(ctx, owner, param.Key)
		require.NoError(t, err)
		assert.Equal(t, group.UserGroupID.Int(), retrievedByKey.UserGroupID.Int())

		groups, err := repo.FindAllUserGroups(ctx, owner)
		require.NoError(t, err)
		assert.NotEmpty(t, groups)

		ownerGroupID, err := repo.CreateOwnerGroup(ctx, sysOwner, orgID)
		require.NoError(t, err)
		ownerGroup, err := repo.FindUserGroupByID(ctx, owner, ownerGroupID)
		require.NoError(t, err)
		assert.Equal(t, service.OwnerGroupKey, ownerGroup.Key)
	}

	testDB(t, fn)
}
