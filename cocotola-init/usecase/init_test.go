package usecase_test

import "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"

var systemAdmin *domain.SystemAdmin

func init() {
	systemAdmin = domain.NewSystemAdmin(domain.NewSystemToken())
}
