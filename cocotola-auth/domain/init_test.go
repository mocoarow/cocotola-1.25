package domain_test

import "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"

var systemAdmin *domain.SystemAdmin

func init() {
	systemAdmin = domain.NewSystemAdmin(NewTestSystemToken())
}

type SystemToken struct {
}

func NewTestSystemToken() *SystemToken {
	return &SystemToken{}
}

func (t *SystemToken) IsSystemToken() bool {
	return true
}
