package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller"
	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	libprocess "github.com/mocoarow/cocotola-1.25/cocotola-lib/process"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/config"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/initialize"
)

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

	systemToken := domain.NewSystemToken()

	// init log
	shutdownlog, err := libgateway.InitLog(ctx, cfg.Log, domain.AppName)
	if err != nil {
		return 0, fmt.Errorf("init log: %w", err)
	}
	defer shutdownlog()
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, domain.AppName+"-main"))

	// init tracer
	shutdownTrace, err := libgateway.InitTracerProvider(ctx, cfg.Trace, domain.AppName)
	if err != nil {
		return 0, fmt.Errorf("init trace: %w", err)
	}
	defer shutdownTrace()

	// init db
	dbConn, shutdownDB, err := libgateway.InitDB(ctx, cfg.DB, cfg.Log, domain.AppName)
	if err != nil {
		return 0, fmt.Errorf("init db: %w", err)
	}
	defer shutdownDB()

	// init gin
	router := libgin.InitRootRouterGroup(ctx, cfg.Server.Gin, domain.AppName)

	if err := initialize.Initialize(ctx, systemToken, router, dbConn, cfg.Server.Gin.Log, cfg.App); err != nil {
		return 0, fmt.Errorf("initialize: %w", err)
	}

	// run
	readHeaderTimeout := time.Duration(cfg.Server.ReadHeaderTimeoutSec) * time.Second
	shutdownTime := time.Duration(cfg.Server.Shutdown.TimeSec1) * time.Second
	result := libprocess.Run(ctx,
		libcontroller.WithWebServerProcess(router, cfg.Server.HTTPPort, readHeaderTimeout, shutdownTime),
		libcontroller.WithMetricsServerProcess(cfg.Server.MetricsPort, cfg.Server.Shutdown.TimeSec1),
		libgateway.WithSignalWatchProcess(),
	)

	gracefulShutdownTime2 := time.Duration(cfg.Server.Shutdown.TimeSec2) * time.Second
	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")

	return result, nil
}
