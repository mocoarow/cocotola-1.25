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

	libconfig "github.com/mocoarow/cocotola-1.25/cocotola-lib/config"
	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libprocess "github.com/mocoarow/cocotola-1.25/cocotola-lib/process"
)

type ServerConfig struct {
	HTTPPort             int `yaml:"httpPort" validate:"required"`
	MetricsPort          int `yaml:"metricsPort" validate:"required"`
	ReadHeaderTimeoutSec int `yaml:"readHeaderTimeoutSec" validate:"gte=1"`
}

type Config struct {
	Server   *ServerConfig             `yaml:"server" validate:"required"`
	Trace    *libconfig.TraceConfig    `yaml:"trace" validate:"required"`
	CORS     *libconfig.CORSConfig     `yaml:"cors" validate:"required"`
	Shutdown *libconfig.ShutdownConfig `yaml:"shutdown" validate:"required"`
	Log      *libconfig.LogConfig      `yaml:"log" validate:"required"`
	Debug    *libconfig.DebugConfig    `yaml:"debug"`
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
		},
		Trace: &libconfig.TraceConfig{
			Exporter:           "google",
			SamplingPercentage: 100,
			OTLP:               nil,
			Uptrace:            nil,
			Google: &libconfig.GoogleTraceConfig{
				ProjectID: "mocoarow-25-08",
			},
		},
		CORS: &libconfig.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Content-Type"},
		},
		Shutdown: &libconfig.ShutdownConfig{
			TimeSec1: 10,
			TimeSec2: 10,
		},
		Log: &libconfig.LogConfig{
			Level:    "info",
			Platform: "gcp",
			Levels:   nil,
			Enabled:  nil,
			Exporter: "none",
			OTLP:     nil,
			Uptrace:  nil,
		},
		Debug: &libconfig.DebugConfig{
			Gin:  false,
			Wait: false,
		},
	}

	// init log
	shutdownlog, err := libconfig.InitLog(ctx, cfg.Log, AppName)
	if err != nil {
		return 0, fmt.Errorf("init log: %w", err)
	}
	defer shutdownlog()
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, AppName+"-main"))

	// init tracer
	shutdownTrace, err := libconfig.InitTracerProvider(ctx, cfg.Trace, AppName)
	if err != nil {
		return 0, fmt.Errorf("init trace: %w", err)
	}
	defer shutdownTrace()

	// init gin
	logConfig := libcontroller.LogConfig{
		Enabled: map[string]bool{
			"accessLog":             true,
			"accessLogRequestBody":  true,
			"accessLogResponseBody": true,
		},
	}
	ginConfig := libcontroller.GinConfig{
		CORS: libconfig.InitCORS(cfg.CORS),
		Log:  logConfig,
		Debug: libcontroller.DebugConfig{
			Gin:  cfg.Debug.Gin,
			Wait: cfg.Debug.Wait,
		},
	}
	router := libcontroller.InitRootRouterGroup(ctx, &ginConfig, AppName)

	// api
	api := libcontroller.InitAPIRouterGroup(ctx, router, AppName, &logConfig)
	// v1
	v1 := api.Group("v1")
	// public router
	libcontroller.InitPublicAPIRouterGroup(ctx, v1, []libcontroller.InitRouterGroupFunc{
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
	shutdownTime := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
	result := libprocess.Run(ctx,
		libprocess.WithAppServerProcess(router, cfg.Server.HTTPPort, readHeaderTimeout, shutdownTime),
		libprocess.WithSignalWatchProcess(),
		libprocess.WithMetricsServerProcess(cfg.Server.MetricsPort, cfg.Shutdown.TimeSec1),
	)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")

	return result, nil
}
