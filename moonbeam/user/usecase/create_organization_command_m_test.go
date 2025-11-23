//go:build medium

package usecase_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"
	testlibgateway "github.com/mocoarow/cocotola-1.25/moonbeam/testlib/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/usecase"
)

func TestCreateOrganizationCommand_Execute_shouldProvisionOrganizationResources_whenCalled(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()
	sysAdmin := domain.NewSystemAdmin()

	for dialect, db := range testlibgateway.ListDB() {
		dialect := dialect
		db := db

		t.Run(dialect.Name(), func(t *testing.T) {
			t.Parallel()

			unlock := acquireCasbinLock(t)
			defer unlock()

			sqlDB, err := db.DB()
			require.NoError(t, err)
			defer sqlDB.Close()

			rf, err := gateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), db, time.UTC)
			require.NoError(t, err)

			txManager, err := libgateway.NewTransactionManagerT(db, func(ctx context.Context, txDB *gorm.DB) (service.RepositoryFactory, error) {
				return gateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), txDB, time.UTC)
			})
			require.NoError(t, err)
			nonTxManager, err := libgateway.NewNonTransactionManagerT[service.RepositoryFactory](rf)
			require.NoError(t, err)

			cmd := usecase.NewCreateOrganizationCommand(ctx, txManager, nonTxManager)

			orgName := fmt.Sprintf("org-%s", randString(10))
			orgID, err := cmd.Execute(ctx, sysAdmin, orgName)
			require.NoError(t, err)
			require.NotNil(t, orgID)

			t.Cleanup(func() {
				cleanupOrganization(t, db, orgID)
			})

			userRepo := rf.NewUserRepository(ctx)
			sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, sysAdmin, orgID)
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

			t.Run("system owner policies allow managing user roles", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				allUserRoles := domain.NewRBACAllUserRolesObjectFromOrganization(orgID)
				rbacDomain := domain.NewRBACDomainFromOrganization(orgID)
				rbacSystemOwner := sysOwner.GetUserID().GetRBACSubject()

				requireCasbinPolicy(t, ctx, db, []string{
					"p",
					rbacSystemOwner.Subject(),
					allUserRoles.Object(),
					service.RBACSetAction.Action(),
					service.RBACAllowEffect.Effect(),
					rbacDomain.Domain(),
				})

				requireCasbinPolicy(t, ctx, db, []string{
					"p",
					rbacSystemOwner.Subject(),
					allUserRoles.Object(),
					service.RBACUnsetAction.Action(),
					service.RBACAllowEffect.Effect(),
					rbacDomain.Domain(),
				})
			})

			t.Run("owner group policies allow managing user roles", func(t *testing.T) { //nolint:paralleltest
				t.Helper()
				require.NotNil(t, ownerGroup, "owner group must be created before checking policies")

				allUserRoles := domain.NewRBACAllUserRolesObjectFromOrganization(orgID)
				rbacDomain := domain.NewRBACDomainFromOrganization(orgID)
				rbacOwnerGroup := domain.NewRBACRoleFromGroup(orgID, ownerGroup.UserGroupID)

				requireCasbinPolicy(t, ctx, db, []string{
					"p",
					rbacOwnerGroup.Subject(),
					allUserRoles.Object(),
					service.RBACSetAction.Action(),
					service.RBACAllowEffect.Effect(),
					rbacDomain.Domain(),
				})

				requireCasbinPolicy(t, ctx, db, []string{
					"p",
					rbacOwnerGroup.Subject(),
					allUserRoles.Object(),
					service.RBACUnsetAction.Action(),
					service.RBACAllowEffect.Effect(),
					rbacDomain.Domain(),
				})
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

var casbinMu sync.Mutex

func acquireCasbinLock(t *testing.T) func() {
	t.Helper()
	casbinMu.Lock()
	return func() {
		casbinMu.Unlock()
	}
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

func requireCasbinPolicy(t *testing.T, ctx context.Context, db *gorm.DB, expected []string) {
	t.Helper()

	require.Len(t, expected, 6)

	type casbinRule struct {
		ID int
	}

	var rule casbinRule
	err := db.WithContext(ctx).Table("casbin_rule").Where(
		"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ? AND v4 = ?",
		expected[0], expected[1], expected[2], expected[3], expected[4], expected[5],
	).First(&rule).Error
	require.NoError(t, err, "casbin policy not found: %v", expected)
}
