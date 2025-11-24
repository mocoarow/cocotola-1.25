//go:build medium

package gateway_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway"
)

func setupTestOrganization(ctx context.Context, t *testing.T, tr testResource) (*domain.OrganizationID, domain.SystemOwnerInterface, domain.OwnerInterface) {
	t.Helper()
	orgRepo := gateway.NewOrganizationRepository(ctx, tr.db)
	userRepo := tr.rf.NewUserRepository(ctx)
	sysAdmin := domain.NewSystemAdmin()

	orgName := fmt.Sprintf("org-%s", RandString(8))
	orgID, err := orgRepo.CreateOrganization(ctx, sysAdmin, orgName)
	require.NoError(t, err)

	_, err = userRepo.CreateSystemOwner(ctx, sysAdmin, orgID)
	require.NoError(t, err)

	sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, sysAdmin, orgID)
	require.NoError(t, err)

	ownerParam := testNewCreateUserParameter(t, fmt.Sprintf("owner_%s", RandString(6)), fmt.Sprintf("Owner %s", RandString(6)), "password-owner")
	ownerID, err := userRepo.CreateUser(ctx, sysOwner, ownerParam)
	require.NoError(t, err)

	ownerUser, err := userRepo.FindUserByID(ctx, sysOwner, ownerID)
	require.NoError(t, err)

	owner, err := domain.NewOwner(ownerUser)
	require.NoError(t, err)

	return orgID, sysOwner, owner
}
