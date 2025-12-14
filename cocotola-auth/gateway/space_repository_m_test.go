//go:build medium

package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

func TestSpaceRepository_CreateAndQuerySpaces(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		orgID, _, owner := setupTestOrganization(ctx, t, tr)
		defer teardownOrganization(t, tr, orgID)

		repo := gateway.NewSpaceRepository(ctx, tr.dialect, tr.db)
		param := &service.CreateSpaceParameter{
			Key:       "space-key",
			Name:      "Test Space",
			SpaceType: "public",
		}
		spaceID, err := repo.CreateSpace(ctx, owner, param)
		require.NoError(t, err)

		space, err := repo.GetSpaceByID(ctx, owner, spaceID)
		require.NoError(t, err)
		assert.Equal(t, param.Key, space.KeyName)
		assert.Equal(t, param.Name, space.Name)

		publicSpaces, err := repo.FindPublicSpaces(ctx, owner)
		require.NoError(t, err)
		assert.NotEmpty(t, publicSpaces)

		spaceByKey, err := repo.FindPublicSpaceByKey(ctx, owner, param.Key)
		require.NoError(t, err)
		assert.Equal(t, spaceID.Int(), spaceByKey.SpaceID.Int())
	}

	testDB(t, fn)
}
