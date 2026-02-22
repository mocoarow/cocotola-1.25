package initialize

import (
	"context"
	"fmt"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-init/config"
)

func Initialize(ctx context.Context, systemToken authdomain.SystemToken, dbConn *libgateway.DBConnection, _ *libgateway.LogConfig, initConfig *config.InitConfig, appName string) error {
	ctx, span := tracer.Start(ctx, "Initialize")
	defer span.End()

	if err := initOrganization(ctx, systemToken, dbConn, "cocotola", initConfig.OwnerLoginID, initConfig.OwnerPassword); err != nil {
		return fmt.Errorf("initOrganization: %w", err)
	}

	if err := initGuest(ctx, systemToken, dbConn, "cocotola", appName); err != nil {
		return fmt.Errorf("initGuest: %w", err)
	}

	return nil
}
