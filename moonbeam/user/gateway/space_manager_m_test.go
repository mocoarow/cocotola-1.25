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

func TestSpaceManager_CreatePersonalSpace_and_GetPersonalSpace(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, sysOwner, owner := setupTestOrganization(ctx, t, tr)
		_ = sysOwner
		defer teardownOrganization(t, tr, orgID)

		mgr, err := gateway.NewSpaceManager(ctx, tr.dialect, tr.db, tr.rf)
		require.NoError(t, err)

		target := testAddUser(t, ctx, tr, owner, "space_user", "SPACE USER", "password")

		spaceID, err := mgr.CreatePersonalSpace(ctx, owner, &service.CreatePersonalSpaceParameter{
			UserID:  target.GetUserID(),
			KeyName: "personal-space",
			Name:    "Personal Space",
		})
		require.NoError(t, err)
		assert.Positive(t, spaceID.Int())

		space, err := mgr.GetPersonalSpace(ctx, target)
		require.NoError(t, err)
		assert.Equal(t, spaceID.Int(), space.SpaceID.Int())
		assert.Equal(t, "personal", space.SpaceType)
	}

	testDB(t, fn)
}

func TestSpaceManager_AddUserToSpace_shouldAttachExistingSpace(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		unlock := acquireCasbinLock(t)
		defer unlock()
		defer teardownCasbin(t, tr)

		orgID, sysOwner, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		mgr, err := gateway.NewSpaceManager(ctx, tr.dialect, tr.db, tr.rf)
		require.NoError(t, err)

		spaceRepo := tr.rf.NewSpaceRepository(ctx)
		spaceID, err := spaceRepo.CreateSpace(ctx, owner, &service.CreateSpaceParameter{
			Key:       "shared-space",
			Name:      "Shared Space",
			SpaceType: "private",
		})
		require.NoError(t, err)

		other := testAddUser(t, ctx, tr, owner, "new_member", "NEW MEMBER", "password")
		err = mgr.AddUserToSpace(ctx, sysOwner, *other.GetUserID(), spaceID)
		require.NoError(t, err)

		rbacRepo, err := gateway.NewRBACRepository(ctx, tr.db)
		require.NoError(t, err)
		rbacDomain := domain.NewRBACDomainFromOrganization(orgID)
		rbacUser := domain.NewRBACUserFromUser(other.GetUserID())
		roles, err := rbacRepo.GetGroupsForSubject(ctx, rbacDomain, rbacUser)
		require.NoError(t, err)
		rbacSpace := domain.NewRBACRoleFromSpace(orgID, spaceID)
		found := false
		for _, role := range roles {
			if role.Role() == rbacSpace.Role() {
				found = true
				break
			}
		}
		assert.True(t, found)

		_, err = mgr.GetPersonalSpace(ctx, other)
		assert.ErrorIs(t, err, service.ErrSpaceNotFound)
	}

	testDB(t, fn)
}
