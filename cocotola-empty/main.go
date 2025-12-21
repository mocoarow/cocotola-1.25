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
	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	libprocess "github.com/mocoarow/cocotola-1.25/cocotola-lib/process"
)

type ServerConfig struct {
	HTTPPort             int                           `yaml:"httpPort" validate:"required"`
	MetricsPort          int                           `yaml:"metricsPort" validate:"required"`
	ReadHeaderTimeoutSec int                           `yaml:"readHeaderTimeoutSec" validate:"gte=1"`
	Gin                  *libgin.Config                `yaml:"gin" validate:"required"`
	Shutdown             *libcontroller.ShutdownConfig `yaml:"shutdown" validate:"required"`
}

type Config struct {
	Server *ServerConfig           `yaml:"server" validate:"required"`
	Trace  *libgateway.TraceConfig `yaml:"trace" validate:"required"`
	Log    *libgateway.LogConfig   `yaml:"log" validate:"required"`
}

const AppName = "cocotola-empty"

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
			HTTPPort:             8080,
			MetricsPort:          8081,
			ReadHeaderTimeoutSec: 10,
			Gin: &libgin.Config{
				CORS: &libgin.CORSConfig{
					AllowOrigins: "*",
					AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
					AllowHeaders: "Content-Type",
				},
				Log: &libgin.LogConfig{
					AccessLog:             true,
					AccessLogRequestBody:  true,
					AccessLogResponseBody: true,
				},
				Debug: &libgin.DebugConfig{
					Gin:  false,
					Wait: false,
				},
			},
			Shutdown: &libcontroller.ShutdownConfig{
				TimeSec1: 10,
				TimeSec2: 10,
			},
		},
		Trace: &libgateway.TraceConfig{
			Exporter:           "google",
			SamplingPercentage: 100,
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

	// init gin
	router := libgin.InitRootRouterGroup(ctx, cfg.Server.Gin, AppName)

	// api
	api := libgin.InitAPIRouterGroup(ctx, router, cfg.Server.Gin.Log, AppName)
	// v1
	v1 := api.Group("v1")
	// public router
	libgin.InitPublicAPIRouterGroup(ctx, v1, []libgin.InitRouterGroupFunc{
		func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc) {
			test := parentRouterGroup.Group("test")
			for _, m := range middleware {
				test.Use(m)
			}
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
		},
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
