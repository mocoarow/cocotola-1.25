//go:build medium

package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
)

const orgNameLength = 8

func Test_organizationRepository_CreateOrganization(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testResource) {
		t.Helper()
		sysAd := domain.NewSystemAdmin()
		orgName := RandString(orgNameLength)

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)
		_, err := orgRepo.CreateOrganization(ctx, sysAd, orgName)
		require.NoError(t, err)
	}
	testDB(t, fn)
}

func Test_organizationRepository_FindOrganizationByID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testResource) {
		t.Helper()
		sysAd := domain.NewSystemAdmin()
		orgName := RandString(orgNameLength)

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)
		orgID, err := orgRepo.CreateOrganization(ctx, sysAd, orgName)
		require.NoError(t, err)
		org, err := orgRepo.FindOrganizationByID(ctx, sysAd, orgID)
		require.NoError(t, err)
		assert.Equal(t, orgName, org.Name)
	}
	testDB(t, fn)
}

func Test_organizationRepository_FindOrganizationByName(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testResource) {
		t.Helper()
		sysAd := domain.NewSystemAdmin()
		orgName := RandString(orgNameLength)

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)
		_, err := orgRepo.CreateOrganization(ctx, sysAd, orgName)
		require.NoError(t, err)
		org, err := orgRepo.FindOrganizationByName(ctx, sysAd, orgName)
		require.NoError(t, err)
		assert.Equal(t, orgName, org.Name)
	}
	testDB(t, fn)
}
