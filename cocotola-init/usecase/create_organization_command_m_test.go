//go:build medium

package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	testlibgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/testlib/gateway"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-init/usecase"
)

func TestCreateOrganizationCommand_Execute_shouldProvisionOrganizationResources_whenCalled(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	for dialect, db := range testlibgateway.ListDB() {
		dialect := dialect
		db := db

		t.Run(dialect.Name(), func(t *testing.T) {
			t.Parallel()

			sqlDB, err := db.DB()
			require.NoError(t, err)

			rf, err := authgateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), db, time.UTC)
			require.NoError(t, err)

			txManager, err := libgateway.NewTransactionManagerT(db, func(ctx context.Context, txDB *gorm.DB) (authservice.RepositoryFactory, error) {
				return authgateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), txDB, time.UTC)
			})
			require.NoError(t, err)
			nonTxManager, err := libgateway.NewNonTransactionManagerT(rf)
			require.NoError(t, err)

			// when
			cmd := usecase.NewCreateOrganizationCommand(ctx, txManager, nonTxManager)

			orgName := fmt.Sprintf("org-%s", randString(10))
			orgID, err := cmd.Execute(ctx, systemAdmin, orgName)
			require.NoError(t, err)
			require.NotNil(t, orgID)

			t.Cleanup(func() {
				cleanupOrganization(t, db, orgID)
				sqlDB.Close()
			})

			userRepo := rf.NewUserRepository(ctx)
			sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, systemAdmin, orgID)
			require.NoError(t, err)

			authorizationManager, err := authgateway.NewAuthorizationManager(ctx, dialect, db, rf)
			require.NoError(t, err)

			ownerParam, err := authservice.NewCreateUserParameter(fmt.Sprintf("owner_%s", randString(6)), fmt.Sprintf("Owner %s", randString(6)), "owner-password", "", "", "", "")
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
				userGroupRepo := rf.NewUserGroupRepository(ctx)
				var err error
				ownerGroup, err = userGroupRepo.FindUserGroupByKey(ctx, sysOwner, authservice.OwnerGroupKey)
				require.NoError(t, err)
				require.Equal(t, authservice.OwnerGroupKey, ownerGroup.Key)
			})

			t.Run("public group is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				userGroupRepo := rf.NewUserGroupRepository(ctx)
				publicGroup, err := userGroupRepo.FindUserGroupByKey(ctx, sysOwner, authservice.PublicGroupKey)
				require.NoError(t, err)
				require.Equal(t, authservice.PublicGroupKey, publicGroup.Key)
			})

			t.Run("system owner policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				ok, err := authorizationManager.CheckAuthorization(ctx, sysOwner, authservice.CreateUserAction, authservice.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)

				// ok, err = authorizationManager.CheckAuthorization(ctx, sysOwner, authservice.RBACUnsetAction, allUserRoles)
				// require.NoError(t, err)
				// require.True(t, ok)
			})

			t.Run("owner group policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				require.NotNil(t, ownerGroup, "owner group must be created before checking policies")

				// allUserRoles := authdomain.NewRBACAllUserRolesObjectFromOrganization(orgID)
				require.NoError(t, authorizationManager.AddUserToGroup(ctx, sysOwner, owner.GetUserID(), ownerGroup.UserGroupID))

				ok, err := authorizationManager.CheckAuthorization(ctx, owner, authservice.CreateUserAction, authservice.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)

				// ok, err = authorizationManager.CheckAuthorization(ctx, owner, authservice.RBACUnsetAction, allUserRoles)
				// require.NoError(t, err)
				// require.True(t, ok)
			})

			t.Run("public default space is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				spaceRepo := rf.NewSpaceRepository(ctx)
				publicSpace, err := spaceRepo.FindPublicSpaceByKey(ctx, sysOwner, authservice.PublicDefaultSpaceKey)
				require.NoError(t, err)
				require.Equal(t, authservice.PublicDefaultSpaceName, publicSpace.Name)
			})
		})
	}
}
