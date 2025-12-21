package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authinit "github.com/mocoarow/cocotola-1.25/cocotola-auth/initialize"

	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller"
	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	libprocess "github.com/mocoarow/cocotola-1.25/cocotola-lib/process"

	"github.com/mocoarow/cocotola-1.25/cocotola-app/config"
)

const AppName = "cocotola-app"

var tracer = otel.Tracer("github.com/mocoarow/cocotola-1.25/cocotola-app")

func main() {
	exitCode, err := run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}

func run() (int, error) {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		return 0, fmt.Errorf("LoadConfig: %w", err)
	}

	// init log
	shutdownlog, err := libgateway.InitLog(ctx, cfg.Log, AppName)
	if err != nil {
		return 0, fmt.Errorf("init log: %w", err)
	}
	defer shutdownlog()
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, AppName+"-main"))

	// init tracer	// init tracer
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

	// init gin
	router := libgin.InitRootRouterGroup(ctx, cfg.Server.Gin, AppName)

	systemToken := authdomain.NewSystemToken()
	{
		auth := router.Group("auth")
		if err := authinit.Initialize(ctx, systemToken, auth, dbConn, cfg.Server.Gin.Log, cfg.App.Auth); err != nil {
			return 0, fmt.Errorf("initialize auth: %w", err)
		}
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
