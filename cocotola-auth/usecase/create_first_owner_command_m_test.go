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

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase"
)

func TestCreateFirstOwnerCommand_Execute_shouldCreateFirstOwner_whenCalled(t *testing.T) { //nolint:paralleltest
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

			createOrgCmd := usecase.NewCreateOrganizationCommand(ctx, txManager, nonTxManager)
			orgID, err := createOrgCmd.Execute(ctx, sysAdmin, fmt.Sprintf("org-%s", randString(8)))
			require.NoError(t, err)
			require.NotNil(t, orgID)

			t.Cleanup(func() {
				cleanupOrganization(t, db, orgID)
				sqlDB.Close()
			})

			userRepo := rf.NewUserRepository(ctx)
			sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, sysAdmin, orgID)
			require.NoError(t, err)

			// when
			cmd := usecase.NewCreateFirstOwnerCommand(txManager, nonTxManager)

			firstOwnerParam, err := service.NewCreateUserParameter(
				fmt.Sprintf("first-owner-%s", randString(6)),
				fmt.Sprintf("First Owner %s", randString(4)),
				"first-owner-password",
				"",
				"",
				"",
				"",
			)
			require.NoError(t, err)

			firstOwnerID, err := cmd.Execute(ctx, sysOwner, firstOwnerParam)
			require.NoError(t, err)
			require.NotNil(t, firstOwnerID)

			ownerUser, err := userRepo.FindUserByID(ctx, sysOwner, firstOwnerID)
			require.NoError(t, err)

			firstOwner, err := domain.NewOwner(ownerUser)
			require.NoError(t, err)

			userGroupRepo := rf.NewUserGroupRepository(ctx)
			ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, sysOwner, service.OwnerGroupKey)
			require.NoError(t, err)

			t.Run("first owner is added to owner group", func(t *testing.T) { //nolint:paralleltest
				t.Helper()

				pairRepo := gateway.NewPairOfUserAndGroupRepository(ctx, dialect, db, rf)
				groups, err := pairRepo.FindUserGroupsByUserID(ctx, sysOwner, firstOwnerID)
				require.NoError(t, err)

				found := false
				for _, group := range groups {
					if group.UserGroupID.Int() == ownerGroup.UserGroupID.Int() {
						found = true
						break
					}
				}
				require.True(t, found, "first owner must belong to owner group")
			})

			t.Run("first owner policies allow creating user", func(t *testing.T) { //nolint:paralleltest
				t.Helper()

				authorizationManager, err := gateway.NewAuthorizationManager(ctx, dialect, db, rf)
				require.NoError(t, err)

				ok, err := authorizationManager.CheckAuthorization(ctx, firstOwner, service.CreateUserAction, service.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)
			})

			t.Run("first owner can create additional users", func(t *testing.T) { //nolint:paralleltest
				t.Helper()

				additionalUserParam, err := service.NewCreateUserParameter(
					fmt.Sprintf("member-%s", randString(6)),
					fmt.Sprintf("Member %s", randString(4)),
					"member-password",
					"",
					"",
					"",
					"",
				)
				require.NoError(t, err)

				_, err = userRepo.CreateUser(ctx, firstOwner, additionalUserParam)
				require.NoError(t, err)
			})
		})
	}
}
