//go:build medium

package usecase_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	testlibgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/testlib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase"
)

func TestCreateOrganizationCommand_Execute_shouldProvisionOrganizationResources_whenCalled(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()
	sysAdmin := domain.NewSystemAdmin()

	for dialect, db := range testlibgateway.ListDB() {
		dialect := dialect
		db := db

		t.Run(dialect.Name(), func(t *testing.T) {
			t.Parallel()

			sqlDB, err := db.DB()
			require.NoError(t, err)

			rf, err := gateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), db, time.UTC)
			require.NoError(t, err)

			txManager, err := libgateway.NewTransactionManagerT(db, func(ctx context.Context, txDB *gorm.DB) (service.RepositoryFactory, error) {
				return gateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), txDB, time.UTC)
			})
			require.NoError(t, err)
			nonTxManager, err := libgateway.NewNonTransactionManagerT(rf)
			require.NoError(t, err)

			// when
			cmd := usecase.NewCreateOrganizationCommand(ctx, txManager, nonTxManager)

			orgName := fmt.Sprintf("org-%s", randString(10))
			orgID, err := cmd.Execute(ctx, sysAdmin, orgName)
			require.NoError(t, err)
			require.NotNil(t, orgID)

			t.Cleanup(func() {
				cleanupOrganization(t, db, orgID)
				sqlDB.Close()
			})

			userRepo := rf.NewUserRepository(ctx)
			sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, sysAdmin, orgID)
			require.NoError(t, err)

			authorizationManager, err := gateway.NewAuthorizationManager(ctx, dialect, db, rf)
			require.NoError(t, err)

			ownerParam, err := service.NewCreateUserParameter(fmt.Sprintf("owner_%s", randString(6)), fmt.Sprintf("Owner %s", randString(6)), "owner-password", "", "", "", "")
			require.NoError(t, err)
			ownerID, err := userRepo.CreateUser(ctx, sysOwner, ownerParam)
			require.NoError(t, err)
			ownerUser, err := userRepo.FindUserByID(ctx, sysOwner, ownerID)
			require.NoError(t, err)
			owner, err := domain.NewOwner(ownerUser)
			require.NoError(t, err)

			t.Run("system owner is provisioned", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				require.True(t, sysOwner.IsSystemOwner())
			})

			var ownerGroup *domain.UserGroup
			t.Run("owner group is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				userGroupRepo := rf.NewUserGroupRepository(ctx)
				var err error
				ownerGroup, err = userGroupRepo.FindUserGroupByKey(ctx, sysOwner, service.OwnerGroupKey)
				require.NoError(t, err)
				require.Equal(t, service.OwnerGroupKey, ownerGroup.Key)
			})

			t.Run("public group is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				userGroupRepo := rf.NewUserGroupRepository(ctx)
				publicGroup, err := userGroupRepo.FindUserGroupByKey(ctx, sysOwner, service.PublicGroupKey)
				require.NoError(t, err)
				require.Equal(t, service.PublicGroupKey, publicGroup.Key)
			})

			t.Run("system owner policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				ok, err := authorizationManager.CheckAuthorization(ctx, sysOwner, service.CreateUserAction, service.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)

				// ok, err = authorizationManager.CheckAuthorization(ctx, sysOwner, service.RBACUnsetAction, allUserRoles)
				// require.NoError(t, err)
				// require.True(t, ok)
			})

			t.Run("owner group policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				require.NotNil(t, ownerGroup, "owner group must be created before checking policies")

				// allUserRoles := domain.NewRBACAllUserRolesObjectFromOrganization(orgID)
				require.NoError(t, authorizationManager.AddUserToGroup(ctx, sysOwner, owner.GetUserID(), ownerGroup.UserGroupID))

				ok, err := authorizationManager.CheckAuthorization(ctx, owner, service.CreateUserAction, service.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)

				// ok, err = authorizationManager.CheckAuthorization(ctx, owner, service.RBACUnsetAction, allUserRoles)
				// require.NoError(t, err)
				// require.True(t, ok)
			})

			t.Run("public default space is created", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				spaceRepo := rf.NewSpaceRepository(ctx)
				publicSpace, err := spaceRepo.FindPublicSpaceByKey(ctx, sysOwner, service.PublicDefaultSpaceKey)
				require.NoError(t, err)
				require.Equal(t, service.PublicDefaultSpaceName, publicSpace.Name)
			})
		})
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			panic(err)
		}
		b[i] = letterRunes[val.Int64()]
	}
	return string(b)
}

func cleanupOrganization(t *testing.T, db *gorm.DB, orgID *domain.OrganizationID) {
	t.Helper()

	db.Exec("delete from mb_user_n_space where organization_id = ?", orgID.Int())
	db.Exec("delete from mb_space where organization_id = ?", orgID.Int())
	db.Exec("delete from mb_group_n_group where organization_id = ?", orgID.Int())
	db.Exec("delete from mb_user_n_group where organization_id = ?", orgID.Int())
	db.Exec("delete from mb_user_group where organization_id = ?", orgID.Int())
	db.Exec("delete from mb_user where organization_id = ?", orgID.Int())
	db.Exec("delete from mb_organization where id = ?", orgID.Int())
}
