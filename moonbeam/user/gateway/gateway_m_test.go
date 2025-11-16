//go:build medium

package gateway_test

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"
	testlibgateway "github.com/mocoarow/cocotola-1.25/moonbeam/testlib/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
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

func testDB(t *testing.T, fn func(t *testing.T, ctx context.Context, ts testResource)) {
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
