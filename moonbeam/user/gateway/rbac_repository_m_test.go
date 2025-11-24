//go:build medium

package gateway_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type testSDOA struct {
	subject string
	domain  string
	object  string
	action  string
	want    bool
}

func (t *testSDOA) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%v", t.subject, t.domain, t.object, t.action, t.want)
}

func addPolicy(t *testing.T, ctx context.Context, rbacRepository service.RBACRepository, dom, sub, act, obj string, allowed bool) {
	t.Helper()
	effect := service.RBACAllowEffect
	if !allowed {
		effect = service.RBACDenyEffect
	}

	err := rbacRepository.CreatePolicy(ctx, domain.NewRBACDomain(dom), domain.NewRBACUser(sub), domain.NewRBACAction(act), domain.NewRBACObject(obj), effect)
	require.NoError(t, err)
}

func addObjectGroupingPolicy(t *testing.T, ctx context.Context, rbacRepository service.RBACRepository, dom, child, parent string) {
	t.Helper()
	err := rbacRepository.CreateObjectGroupingPolicy(ctx, domain.NewRBACDomain(dom), domain.NewRBACObject(child), domain.NewRBACObject(parent))
	require.NoError(t, err)
}

func addSubjectGroupingPolicy(t *testing.T, ctx context.Context, rbacRepository service.RBACRepository, dom, sub, obj string) {
	t.Helper()
	err := rbacRepository.CreateSubjectGroupingPolicy(ctx, domain.NewRBACDomain(dom), domain.NewRBACUser(sub), domain.NewRBACRole(obj))
	require.NoError(t, err)
}

func TestA(t *testing.T) { //nolint:paralleltest
	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		// unlock := lockCasbin(t)
		// defer unlock()
		// rbacRepo := gateway.RBACRepository{
		// 	DB:   ts.db,
		// 	Conf: gateway.Conf,
		// }
		rbacRepo, err := gateway.NewRBACRepository(ctx, tr.db)
		require.NoError(t, err)
		e := rbacRepo.GetEnforcer()
		domainName := RandString(orgNameLength)
		domainID, err := RandInt(1000000000)
		require.NoError(t, err)
		// rbacRepo.Init()

		// err := initRBACRepository(t, ts.db, gateway.Conf)
		// require.NoError(t, err)
		addPolicy(t, ctx, rbacRepo, domainName, "alice", "read", fmt.Sprintf("domain:%d,data:1", domainID), true)
		addPolicy(t, ctx, rbacRepo, domainName, "bob", "write", fmt.Sprintf("domain:%d,data:2", domainID), true)
		// rbacRepo.AddPolicy(domain.NewRBACDomain("domain1"), domain.NewRBACUser("alice"), domain.NewRBACAction("write"), domain.NewRBACObject("data1"), service.RBACAllowEffect)
		addObjectGroupingPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,child:1", domainID), fmt.Sprintf("domain:%d,data:1", domainID))

		tests := []testSDOA{
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "read", want: true},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "write", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "read", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "write", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "read", want: true},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "write", want: false},

			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "read", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "write", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "read", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "write", want: true},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "read", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "write", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:1", action: "write", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:2", action: "read", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:2", action: "write", want: true},

			// {subject: "charlie", domain: "domain1", object: "domain:1_data:2", action: "read", want: true},
			// {subject: "charlie", domain: "domain1", object: "domain:1_data_parent", action: "read", want: true},
		}
		// e := initEnforcer(t, ts.db, gateway.Conf)
		// e, err := rbacRepo.InitEnforcer(ctx)
		require.NoError(t, err)
		for _, tt := range tests {
			t.Run(tt.String(), func(t *testing.T) {
				t.Parallel()
				ok, err := e.Enforce(tt.subject, tt.object, tt.action, tt.domain)
				require.NoError(t, err)
				assert.Equal(t, tt.want, ok)
			})
		}
	}
	testDB(t, fn)
}

func TestB(t *testing.T) { //nolint:paralleltest
	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		// unlock := lockCasbin(t)
		// defer unlock()
		rbacRepo, err := gateway.NewRBACRepository(ctx, tr.db)
		require.NoError(t, err)
		e := rbacRepo.GetEnforcer()

		domainName := RandString(orgNameLength)
		domainID, err := RandInt(1000000000)
		require.NoError(t, err)

		addSubjectGroupingPolicy(t, ctx, rbacRepo, domainName, "alice", fmt.Sprintf("domain:%d,reader", domainID))
		addSubjectGroupingPolicy(t, ctx, rbacRepo, domainName, "bob", fmt.Sprintf("domain:%d,writer", domainID))

		addPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,reader", domainID), "read", fmt.Sprintf("domain:%d,data:1", domainID), true)
		addPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,writer", domainID), "write", fmt.Sprintf("domain:%d,data:2", domainID), true)
		addObjectGroupingPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,child:1", domainID), fmt.Sprintf("domain:%d,data:1", domainID))

		tests := []testSDOA{
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "read", want: true},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "write", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "read", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "write", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "read", want: true},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "write", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "read", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "write", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "read", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "write", want: true},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "read", want: false},
			{subject: "bob", domain: domainName, object: fmt.Sprintf("domain:%d,child:1", domainID), action: "write", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:1", action: "write", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:2", action: "read", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:2", action: "write", want: true},

			// {subject: "charlie", domain: "domain1", object: "domain:1_data:2", action: "read", want: true},
			// {subject: "charlie", domain: "domain1", object: "domain:1_data_parent", action: "read", want: true},
		}
		for _, tt := range tests {
			t.Run(tt.String(), func(t *testing.T) {
				t.Parallel()
				ok, err := e.Enforce(tt.subject, tt.object, tt.action, tt.domain)
				require.NoError(t, err)
				assert.Equal(t, tt.want, ok)
			})
		}
	}
	testDB(t, fn)
}

func TestC(t *testing.T) { //nolint:paralleltest
	fn := func(t *testing.T, ctx context.Context, tr testResource) {
		t.Helper()
		// unlock := lockCasbin(t)
		// defer unlock(
		rbacRepo, err := gateway.NewRBACRepository(ctx, tr.db)
		require.NoError(t, err)
		e := rbacRepo.GetEnforcer()

		domainName := RandString(orgNameLength)
		domainID, err := RandInt(1000000000)
		require.NoError(t, err)

		addSubjectGroupingPolicy(t, ctx, rbacRepo, domainName, "alice", fmt.Sprintf("domain:%d,reader", domainID))

		addPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,reader", domainID), "read", fmt.Sprintf("domain:%d,data:2", domainID), true)
		addPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,reader", domainID), "read", fmt.Sprintf("domain:%d,data:4", domainID), false)

		addObjectGroupingPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,data:2", domainID), fmt.Sprintf("domain:%d,data:1", domainID))
		addObjectGroupingPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,data:3", domainID), fmt.Sprintf("domain:%d,data:2", domainID))
		addObjectGroupingPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,data:4", domainID), fmt.Sprintf("domain:%d,data:3", domainID))
		addObjectGroupingPolicy(t, ctx, rbacRepo, domainName, fmt.Sprintf("domain:%d,data:5", domainID), fmt.Sprintf("domain:%d,data:4", domainID))
		// 1/
		// - 2/ <= alice can read
		//   - 3/ <= alice also can read
		//     - 4/ <= alice can't read
		//	     - 5/ <= alice also can't read

		tests := []testSDOA{
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:1", domainID), action: "read", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:2", domainID), action: "read", want: true},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:3", domainID), action: "read", want: true},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:4", domainID), action: "read", want: false},
			{subject: "alice", domain: domainName, object: fmt.Sprintf("domain:%d,data:5", domainID), action: "read", want: false},
		}
		for _, tt := range tests {
			t.Run(tt.String(), func(t *testing.T) {
				t.Parallel()
				ok, err := e.Enforce(tt.subject, tt.object, tt.action, tt.domain)
				require.NoError(t, err)
				assert.Equal(t, tt.want, ok)
			})
		}
	}
	testDB(t, fn)
}

// func TestD(t *testing.T) {
// 	t.Parallel()

// 	fn := func(t *testing.T, ctx context.Context, ts testService) {
// 		t.Helper()
// 		defer teardownCasbin(t, ts)

// 		rbacRepo, err := gateway.NewRBACRepository(ctx, ts.db)
// 		require.NoError(t, err)
// 		e := rbacRepo.GetEnforcer()

// 		addPolicy(t, ctx, rbacRepo, "domain1", "alice", "read", "domain:1,data:1", true)
// 		addPolicy(t, ctx, rbacRepo, "domain1", "bob", "write", "domain:1,data:2", true)

//			s1 := e.GetPermissionsForUserInDomain("alice", "domain1")
//			require.NoError(t, err)
//			for _, s2 := range s1 {
//				t.Log(s2)
//			}
//			t.Fail()
//		}
//		testDB(t, fn)
//	}
