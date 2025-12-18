package domain_test

import "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"

var systemAdmin *domain.SystemAdmin

func init() {
	systemAdmin = domain.NewSystemAdmin(NewTestSystemToken())
}

type systemToken struct {
}

func NewTestSystemToken() domain.SystemToken {
	return &systemToken{}
}

func (t *systemToken) IsSystemToken() bool {
	return true
}
