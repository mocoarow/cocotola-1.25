//go:build medium

package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	testlibgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/testlib/gateway"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"

	"github.com/mocoarow/cocotola-1.25/cocotola-init/initialize"
	"github.com/mocoarow/cocotola-1.25/cocotola-init/usecase"
)

func TestCreateFirstOwnerCommand_Execute_shouldCreateFirstOwner_whenCalled(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	for dialect, dbc := range testlibgateway.ListDB() {
		dialect := dialect
		dbc := dbc

		t.Run(dialect.Name(), func(t *testing.T) {
			t.Parallel()

			sqlDB, err := dbc.DB.DB()
			require.NoError(t, err)

			// Create organization first
			createOrgGw := initialize.NewCreateOrganizationCommandGateway(dbc)
			createOrgCmd := usecase.NewCreateOrganizationCommand(ctx, createOrgGw)
			orgID, err := createOrgCmd.Execute(ctx, systemAdmin, fmt.Sprintf("org-%s", randString(8)))
			require.NoError(t, err)
			require.NotNil(t, orgID)

			t.Cleanup(func() {
				cleanupOrganization(t, dbc, orgID)
				sqlDB.Close()
			})

			userRepo := authgateway.NewUserRepository(dbc)
			sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, systemAdmin, orgID)
			require.NoError(t, err)

			// when
			createFirstOwnerGw := initialize.NewCreateFirstOwnerCommandGateway(dbc)
			cmd := usecase.NewCreateFirstOwnerCommand(ctx, createFirstOwnerGw)

			firstOwnerParam, err := authservice.NewCreateUserParameter(
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

			firstOwner, err := authdomain.NewOwner(ownerUser)
			require.NoError(t, err)

			userGroupRepo := authgateway.NewUserGroupRepository(dbc)
			ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, sysOwner, authservice.OwnerGroupKey)
			require.NoError(t, err)

			t.Run("first owner is added to owner group", func(t *testing.T) { //nolint:paralleltest
				t.Helper()

				pairRepo := authgateway.NewPairOfUserAndGroupRepository(ctx, dbc)
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

				authorizationManager, err := authgateway.NewAuthorizationManager(ctx, dbc)
				require.NoError(t, err)

				ok, err := authorizationManager.CheckAuthorization(ctx, firstOwner, authservice.CreateUserAction, authservice.AnyObject)
				require.NoError(t, err)
				require.True(t, ok)
			})

			t.Run("first owner can create additional users", func(t *testing.T) { //nolint:paralleltest
				t.Helper()

				additionalUserParam, err := authservice.NewCreateUserParameter(
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
