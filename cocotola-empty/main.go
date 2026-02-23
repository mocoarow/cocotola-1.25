package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller"
	libhandler "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/handler"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	libprocess "github.com/mocoarow/cocotola-1.25/cocotola-lib/process"
)

type ServerConfig struct {
	HTTPPort             int                           `yaml:"httpPort" validate:"required"`
	MetricsPort          int                           `yaml:"metricsPort" validate:"required"`
	ReadHeaderTimeoutSec int                           `yaml:"readHeaderTimeoutSec" validate:"gte=1"`
	Gin                  *libhandler.Config            `yaml:"gin" validate:"required"`
	Shutdown             *libcontroller.ShutdownConfig `yaml:"shutdown" validate:"required"`
}

type Config struct {
	Server *ServerConfig           `yaml:"server" validate:"required"`
	Trace  *libgateway.TraceConfig `yaml:"trace" validate:"required"`
	Log    *libgateway.LogConfig   `yaml:"log" validate:"required"`
}

const (
	AppName = "cocotola-empty"

	defaultHTTPPort             = 8080
	defaultMetricsPort          = 8081
	defaultReadHeaderTimeoutSec = 10
	defaultShutdownTimeSec      = 10
	defaultSamplingPercentage   = 100
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

	flag.Parse()

	cfg := &Config{
		Server: &ServerConfig{
			HTTPPort:             defaultHTTPPort,
			MetricsPort:          defaultMetricsPort,
			ReadHeaderTimeoutSec: defaultReadHeaderTimeoutSec,
			Gin: &libhandler.Config{
				CORS: &libhandler.CORSConfig{
					AllowOrigins: "*",
					AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
					AllowHeaders: "Content-Type",
				},
				Log: &libhandler.LogConfig{
					AccessLog:             true,
					AccessLogRequestBody:  true,
					AccessLogResponseBody: true,
				},
				Debug: &libhandler.DebugConfig{
					Gin:  false,
					Wait: false,
				},
			},
			Shutdown: &libcontroller.ShutdownConfig{
				TimeSec1: defaultShutdownTimeSec,
				TimeSec2: defaultShutdownTimeSec,
			},
		},
		Trace: &libgateway.TraceConfig{
			Exporter:           "google",
			SamplingPercentage: defaultSamplingPercentage,
			OTLP:               nil,
			Uptrace:            nil,
			Google: &libgateway.GoogleTraceConfig{
				ProjectID: "mocoarow-25-08",
			},
		},
		Log: &libgateway.LogConfig{
			Level:    "info",
			Platform: "gcp",
			Levels:   nil,
			Exporter: "none",
			OTLP:     nil,
			Uptrace:  nil,
		},
	}

	// init log
	shutdownlog, err := libgateway.InitLog(ctx, cfg.Log, AppName)
	if err != nil {
		return 0, fmt.Errorf("init log: %w", err)
	}
	defer shutdownlog()
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, AppName+"-main"))

	// init tracer
	shutdownTrace, err := libgateway.InitTracerProvider(ctx, cfg.Trace, AppName)
	if err != nil {
		return 0, fmt.Errorf("init trace: %w", err)
	}
	defer shutdownTrace()

	// init handler
	router := libhandler.InitRootRouterGroup(ctx, cfg.Server.Gin, AppName)

	// api
	api := libhandler.InitAPIRouterGroup(ctx, router, cfg.Server.Gin.Log, AppName)
	// v1
	v1 := api.Group("v1")
	// public router
	test := v1.Group("test")
	test.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	test.POST("/200", func(c *gin.Context) {
		logger.InfoContext(ctx, "POST /200")
		params := gin.H{}
		if err := c.BindJSON(&params); err != nil {
			logger.InfoContext(ctx, fmt.Sprintf("err: %+v", err))
			c.Status(http.StatusBadRequest)
			return
		}

		logger.InfoContext(ctx, fmt.Sprintf("params: %+v", params))
		c.Status(http.StatusOK)
	})

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
