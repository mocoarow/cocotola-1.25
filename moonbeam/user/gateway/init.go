package gateway

import (
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("github.com/mocoarow/cocotola-1.25/moonbeam/user/gateway")

	OrganizationTableName = "mb_organization"
)
