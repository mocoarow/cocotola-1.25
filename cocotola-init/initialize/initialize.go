package initialize

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authgateway "github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libconfig "github.com/mocoarow/cocotola-1.25/cocotola-lib/config"

	"github.com/mocoarow/cocotola-1.25/cocotola-init/config"
)

func Initialize(ctx context.Context, systemToken authdomain.SystemToken, dbConn *libconfig.DBConnection, logConfig *libconfig.LogConfig, initConfig *config.InitConfig, appName string) error {
	ctx, span := tracer.Start(ctx, "Initialize")
	defer span.End()

	txManager, nonTxManager, err := initApp(ctx, systemToken, dbConn)
	if err != nil {
		return fmt.Errorf("initApp: %w", err)
	}

	if err := initOrganization(ctx, systemToken, txManager, nonTxManager, "cocotola", initConfig.OwnerLoginID, initConfig.OwnerPassword, appName); err != nil {
		return fmt.Errorf("initOrganization: %w", err)
	}

	if err := initGuest(ctx, systemToken, txManager, nonTxManager, "cocotola", appName); err != nil {
		return fmt.Errorf("initGuest: %w", err)
	}

	return nil
}

func initApp(ctx context.Context, systemToken authdomain.SystemToken, dbConn *libconfig.DBConnection) (authservice.TransactionManager, authservice.TransactionManager, error) {
	rff := func(ctx context.Context, db *gorm.DB) (authservice.RepositoryFactory, error) {
		return authgateway.NewRepositoryFactory(ctx, dbConn.Dialect, dbConn.DriverName, db, time.UTC)
	}
	rf, err := rff(ctx, dbConn.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("rff: %w", err)
	}

	// init transaction manager
	txManager, err := initTransactionManager(dbConn.DB, rff)
	if err != nil {
		return nil, nil, fmt.Errorf("initTransactionManager: %w", err)
	}
	nonTxManager, err := initNonTransactionManager(rf)
	if err != nil {
		return nil, nil, fmt.Errorf("initNonTransactionManager: %w", err)
	}

	return txManager, nonTxManager, nil
}
