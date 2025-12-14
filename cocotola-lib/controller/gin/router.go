package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin/middleware"
)

type LogConfig struct {
	Enabled map[string]bool `yaml:"enabled"`
}
type DebugConfig struct {
	Gin  bool `yaml:"gin"`
	Wait bool `yaml:"wait"`
}
type GinConfig struct {
	CORS  cors.Config
	Log   LogConfig
	Debug DebugConfig
}
type InitRouterGroupFunc func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc)

func InitRootRouterGroup(_ context.Context, ginConfig *GinConfig, appName string) *gin.Engine {
	if !ginConfig.Debug.Gin {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(cors.New(ginConfig.CORS))
	router.Use(middleware.PrometheusMiddleware())
	router.Use(otelgin.Middleware(appName, otelgin.WithFilter(func(req *http.Request) bool {
		return req.URL.Path != "/"
	})))

	if value, ok := ginConfig.Log.Enabled["accessLog"]; ok && value {
		withRequestBody := false
		if value, ok := ginConfig.Log.Enabled["accessLogRequestBody"]; ok && value {
			withRequestBody = true
		}
		withResponseBody := false
		if value, ok := ginConfig.Log.Enabled["accessLogResponseBody"]; ok && value {
			withResponseBody = value
		}
		router.Use(sloggin.NewWithConfig(slog.Default(), sloggin.Config{ //nolint:exhaustruct
			DefaultLevel:     slog.LevelInfo,
			ClientErrorLevel: slog.LevelWarn,
			ServerErrorLevel: slog.LevelError,
			WithRequestBody:  withRequestBody,
			WithResponseBody: withResponseBody,
			Filters: []sloggin.Filter{
				func(c *gin.Context) bool {
					path := c.Request.URL.Path
					return path != "/"
				},
			},
		}))
	}

	if ginConfig.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware(time.Second))
	}

	return router
}

func InitAPIRouterGroup(_ context.Context, parentRouterGroup gin.IRouter, appName string, logConfig *LogConfig) *gin.RouterGroup {
	api := parentRouterGroup.Group("api")
	api.Use(otelgin.Middleware(appName))
	if value, ok := logConfig.Enabled["accessLog"]; ok && value {
		api.Use(sloggin.New(slog.Default()))
	}

	// if value, ok := logConfig.Enabled["traceLog"]; ok && value {
	// 	api.Use(middleware.NewTraceLogMiddleware(appName, true))
	// } else {
	// 	api.Use(middleware.NewTraceLogMiddleware(appName, false))
	// }

	return api
}

func InitPublicAPIRouterGroup(_ context.Context, parentRouterGroup gin.IRouter, initPublicRouterFunc []InitRouterGroupFunc, middleware ...gin.HandlerFunc) {
	for _, fn := range initPublicRouterFunc {
		fn(parentRouterGroup, middleware...)
	}
}

func InitPrivateAPIRouterGroup(_ context.Context, parentRouterGroup gin.IRouter, authMiddleware gin.HandlerFunc, initPrivateRouterFunc []InitRouterGroupFunc) {
	for _, fn := range initPrivateRouterFunc {
		fn(parentRouterGroup, authMiddleware)
	}
}
