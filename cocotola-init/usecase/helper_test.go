package usecase_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

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

func cleanupOrganization(t *testing.T, dbc *libgateway.DBConnection, orgID *domain.OrganizationID) {
	t.Helper()

	dbc.DB.Exec("delete from mb_user_n_space where organization_id = ?", orgID.Int())
	dbc.DB.Exec("delete from mb_space where organization_id = ?", orgID.Int())
	dbc.DB.Exec("delete from mb_group_n_group where organization_id = ?", orgID.Int())
	dbc.DB.Exec("delete from mb_user_n_group where organization_id = ?", orgID.Int())
	dbc.DB.Exec("delete from mb_user_group where organization_id = ?", orgID.Int())
	dbc.DB.Exec("delete from mb_user where organization_id = ?", orgID.Int())
	dbc.DB.Exec("delete from mb_organization where id = ?", orgID.Int())
}
