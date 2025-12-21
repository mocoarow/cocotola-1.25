package main

import (
	"context"
	"fmt"
	"log"
	"os"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-init/config"
	"github.com/mocoarow/cocotola-1.25/cocotola-init/initialize"
)

const AppName = "cocotola-init"

func main() {
	exitCode, err := run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}

func run() (int, error) {
	ctx := context.Background()

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return 0, fmt.Errorf("LoadConfig: %w", err)
	}

	systemToken := authdomain.NewSystemToken()

	// init log
	shutdownlog, err := libgateway.InitLog(ctx, cfg.Log, AppName)
	if err != nil {
		return 0, fmt.Errorf("init log: %w", err)
	}
	defer shutdownlog()

	// init tracer
	shutdownTrace, err := libgateway.InitTracerProvider(ctx, cfg.Trace, AppName)
	if err != nil {
		return 0, fmt.Errorf("init trace: %w", err)
	}
	defer shutdownTrace()

	// init db
	dbConn, shutdownDB, err := libgateway.InitDB(ctx, cfg.DB, cfg.Log, AppName)
	if err != nil {
		return 0, fmt.Errorf("init db: %w", err)
	}
	defer shutdownDB()

	if err := initialize.Initialize(ctx, systemToken, dbConn, cfg.Log, cfg.App, AppName); err != nil {
		return 0, fmt.Errorf("initialize: %w", err)
	}

	return 0, nil
}
