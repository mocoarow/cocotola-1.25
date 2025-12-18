package usecase_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"gorm.io/gorm"
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
