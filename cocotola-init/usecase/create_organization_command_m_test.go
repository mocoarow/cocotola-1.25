//go:build medium

package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	testlibgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/testlib/gateway"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-init/initialize"
	"github.com/mocoarow/cocotola-1.25/cocotola-init/usecase"
)

func TestCreateOrganizationCommand_Execute_shouldProvisionOrganizationResources_whenCalled(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	for dialect, dbc := range testlibgateway.ListDB() {
		t.Run(dialect.Name(), func(t *testing.T) {
			t.Parallel()

			sqlDB, err := dbc.DB.DB()
			require.NoError(t, err)

			// when
			gw := initialize.NewCreateOrganizationCommandGateway(dbc)
			cmd := usecase.NewCreateOrganizationCommand(ctx, gw)

			orgName := "org-" + randString(10)
			orgID, err := cmd.Execute(ctx, systemAdmin, orgName)
			require.NoError(t, err)
			require.NotNil(t, orgID)

			t.Cleanup(func() {
				cleanupOrganization(t, dbc, orgID)
				sqlDB.Close()
			})

			userRepo := authgateway.NewUserRepository(dbc)
			sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, systemAdmin, orgID)
			require.NoError(t, err)

			authorizationManager, err := authgateway.NewAuthorizationManager(ctx, dbc)
			require.NoError(t, err)

			ownerParam, err := authservice.NewCreateUserParameter("owner_"+randString(6), "Owner "+randString(6), "owner-password", "", "", "", "")
			require.NoError(t, err)
			ownerID, err := userRepo.CreateUser(ctx, sysOwner, ownerParam)
			require.NoError(t, err)
			ownerUser, err := userRepo.FindUserByID(ctx, sysOwner, ownerID)
			require.NoError(t, err)
			owner, err := authdomain.NewOwner(ownerUser)
			require.NoError(t, err)

			t.Run("system owner is provisioned", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				require.True(t, sysOwner.IsSystemOwner())
			})

			var ownerGroup *authdomain.UserGroup
			t.Run("owner group is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				userGroupRepo := authgateway.NewUserGroupRepository(dbc)
				var err error
				ownerGroup, err = userGroupRepo.FindUserGroupByKey(ctx, sysOwner, authservice.OwnerGroupKey)
				require.NoError(t, err)
				require.Equal(t, authservice.OwnerGroupKey, ownerGroup.Key)
			})

			t.Run("public group is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				userGroupRepo := authgateway.NewUserGroupRepository(dbc)
				publicGroup, err := userGroupRepo.FindUserGroupByKey(ctx, sysOwner, authservice.PublicGroupKey)
				require.NoError(t, err)
				require.Equal(t, authservice.PublicGroupKey, publicGroup.Key)
			})

			t.Run("system owner policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				ok, err := authorizationManager.CheckAuthorization(ctx, sysOwner, authservice.CreateUserAction, authservice.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)
			})

			t.Run("owner group policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				require.NotNil(t, ownerGroup, "owner group must be created before checking policies")

				require.NoError(t, authorizationManager.AddUserToGroup(ctx, sysOwner, owner.GetUserID(), ownerGroup.UserGroupID))

				ok, err := authorizationManager.CheckAuthorization(ctx, owner, authservice.CreateUserAction, authservice.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)
			})

			t.Run("public default space is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				spaceRepo := authgateway.NewSpaceRepository(dbc)
				publicSpace, err := spaceRepo.FindPublicSpaceByKey(ctx, sysOwner, authservice.PublicDefaultSpaceKey)
				require.NoError(t, err)
				require.Equal(t, authservice.PublicDefaultSpaceName, publicSpace.Name)
			})
		})
	}
}
