package gateway

import (
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway")

	OrganizationTableName       = "mb_organization"
	UserTableName               = "mb_user"
	PairOfUserAndGroupTableName = "mb_user_n_group"
	UserGroupTableName          = "mb_user_group"
	SpaceTableName              = "mb_space"
)
