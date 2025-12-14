//go:build medium

package gateway_test

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	testlibgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/testlib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type testResource struct {
	dialect libgateway.DialectRDBMS
	db      *gorm.DB
	rf      service.RepositoryFactory
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
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

func RandInt(v int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(v)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func testDB(t *testing.T, fn func(t *testing.T, ctx context.Context, tr testResource)) {
	t.Helper()
	ctx := context.Background()

	for dialect, db := range testlibgateway.ListDB() {
		dialect := dialect
		db := db
		t.Run(dialect.Name(), func(t *testing.T) {
			// t.Parallel()
			sqlDB, err := db.DB()
			require.NoError(t, err)
			defer sqlDB.Close()

			rf, err := gateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), db, loc)
			require.NoError(t, err)
			testResource := testResource{dialect: dialect, db: db, rf: rf}

			fn(t, ctx, testResource)
		})
	}
}

// func setupOrganization(ctx context.Context, t *testing.T, tr testResource) (*domain.OrganizationID, *domain.SystemOwner, *domain.Owner) {
// 	t.Helper()
// 	orgName := RandString(orgNameLength)
// 	sysAd := domain.NewSystemAdmin()

// 	slog.Default().Info("========================== " + t.Name() + ", " + orgName + " =========================")
// 	t.Log("--------------------------" + t.Name() + ", " + orgName)

// 	firstOwnerAddParam, err := service.NewCreateUserParameter("OWNER_ID", "OWNER_NAME", "OWNER_PASSWORD", "", "", "", "")
// 	require.NoError(t, err)
// 	// orgAddParam, err := service.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
// 	// require.NoError(t, err)

// 	orgRepo := gateway.NewOrganizationRepository(ctx, tr.db)
// 	userRepo := gateway.NewUserRepository(ctx, tr.dialect, tr.db, tr.rf)
// 	userGorupRepo := gateway.NewUserGroupRepository(ctx, tr.dialect, tr.db)
// 	authorizationManager, err := gateway.NewAuthorizationManager(ctx, tr.dialect, tr.db, tr.rf)
// 	require.NoError(t, err)

// 	// 1. add organization
// 	t.Logf("add organization: %s", orgName)
// 	orgID, err := orgRepo.CreateOrganization(ctx, sysAd, orgName)
// 	if err != nil {
// 		outputOrganization(t, tr.db)
// 		require.NoError(t, err)
// 	}
// 	outputOrganization(t, tr.db)
// 	require.NoError(t, err)
// 	assert.Positive(t, orgID.Int())

// 	t.Logf("organization(%d, %s)", orgID.Int(), orgName)
// 	// 2. add "system-owner" user
// 	sysOwnerID, err := userRepo.CreateSystemOwner(ctx, sysAd, orgID)
// 	require.NoError(t, err)
// 	require.Positive(t, sysOwnerID.Int())

// 	// TODO
// 	sysOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, orgName)
// 	require.NoError(t, err)

// 	// 3. add policy to "system-owner" user
// 	t.Log(`add policy to "system-owner" user`)
// 	rbacSysOwner := domain.NewRBACUserFromUser(sysOwnerID)
// 	rbacAllUserRolesObject := domain.NewRBACAllUserRolesObjectFromOrganization(orgID)
// 	// - "system-owner" "can" "set" "all-user-roles"
// 	err = authorizationManager.AddPolicyToUserBySystemAdmin(ctx, sysAd, orgID, rbacSysOwner, service.RBACSetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
// 	require.NoError(t, err)
// 	// outputCasbinRule(t, ts.db)

// 	// - "system-owner" "can" "unset" "all-user-roles"
// 	err = authorizationManager.AddPolicyToUserBySystemAdmin(ctx, sysAd, orgID, rbacSysOwner, service.RBACUnsetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
// 	require.NoError(t, err)

// 	// 4. create owner-group
// 	t.Logf("create owner group in organization(%d)", orgID.Int())
// 	ownerGroupID, err := userGorupRepo.CreateOwnerGroup(ctx, sysOwner, orgID)
// 	require.NoError(t, err)

// 	// 5. add policty to "owner" group
// 	rbacOwnerGroup := domain.NewRBACRoleFromGroup(orgID, ownerGroupID)
// 	// - "owner" group "can" "set" "all-user-roles"
// 	err = authorizationManager.AddPolicyToGroupBySystemAdmin(ctx, sysAd, orgID, rbacOwnerGroup, service.RBACSetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
// 	require.NoError(t, err)
// 	// - "owner" group "can" "unset" "all-user-roles"
// 	err = authorizationManager.AddPolicyToGroupBySystemAdmin(ctx, sysAd, orgID, rbacOwnerGroup, service.RBACUnsetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
// 	require.NoError(t, err)

// 	// 6. create first owner
// 	ownerID, err := userRepo.CreateUser(ctx, sysOwner, firstOwnerAddParam)
// 	require.NoError(t, err)
// 	require.Positive(t, ownerID.Int())

// 	// - owner belongs to owner-group
// 	err = authorizationManager.AddUserToGroup(ctx, sysOwner, ownerID, ownerGroupID)
// 	require.NoError(t, err)

// 	owner, err := userRepo.FindOwnerByLoginID(ctx, sysOwner, firstOwnerAddParam.LoginID)
// 	require.NoError(t, err)

// 	// logger := slog.Default()
// 	// logger.Warn(fmt.Sprintf("orgID: %d", orgID.Int()))

// 	return orgID, sysOwner, owner
// }

func teardownOrganization(t *testing.T, tr testResource, orgID *domain.OrganizationID) {
	t.Helper()
	// delete all organizations
	tr.db.Exec("delete from mb_user_n_space where organization_id = ?", orgID.Int())
	tr.db.Exec("delete from mb_space where organization_id = ?", orgID.Int())
	tr.db.Exec("delete from mb_group_n_group where organization_id = ?", orgID.Int())
	tr.db.Exec("delete from mb_user_n_group where organization_id = ?", orgID.Int())
	tr.db.Exec("delete from mb_user_group where organization_id = ?", orgID.Int())
	tr.db.Exec("delete from mb_user where organization_id = ?", orgID.Int())
	tr.db.Exec("delete from mb_organization where id = ?", orgID.Int())
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&userEntity{})
	// db.Where("true").Delete(&organizationEntity{})
	slog.Default().Info("teardown organization", "orgID", orgID.Int())
}

func testNewCreateUserParameter(t *testing.T, loginID, username, password string) *service.CreateUserParameter {
	t.Helper()
	p, err := service.NewCreateUserParameter(loginID, username, password, "", "", "", "")
	require.NoError(t, err)

	return p
}

func testAddUser(t *testing.T, ctx context.Context, tr testResource, owner domain.OwnerInterface, loginID, username, password string) *domain.User {
	t.Helper()
	userRepo := tr.rf.NewUserRepository(ctx)
	userID1, err := userRepo.CreateUser(ctx, owner, testNewCreateUserParameter(t, loginID, username, password))
	require.NoError(t, err)
	user1, err := userRepo.FindUserByID(ctx, owner, userID1)
	require.NoError(t, err)
	require.Equal(t, loginID, user1.LoginID)

	return user1
}

func testNewUserGroupAddParameter(t *testing.T, key, name, description string) *service.AddUserGroupParameter {
	t.Helper()
	p, err := service.NewAddUserGroupParameter(key, name, description)
	require.NoError(t, err)

	return p
}

func testAddUserGroup(t *testing.T, ctx context.Context, tr testResource, owner domain.OwnerInterface, key, name, description string) *domain.UserGroup {
	t.Helper()
	userGorupRepo := tr.rf.NewUserGroupRepository(ctx)
	groupID1, err := userGorupRepo.AddUserGroup(ctx, owner, testNewUserGroupAddParameter(t, key, name, description))
	require.NoError(t, err)
	group1, err := userGorupRepo.FindUserGroupByID(ctx, owner, groupID1)
	require.NoError(t, err)
	require.Equal(t, key, group1.Key)
	require.Equal(t, name, group1.Name)
	require.Equal(t, description, group1.Description)

	return group1
}
