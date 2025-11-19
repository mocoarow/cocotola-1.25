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

	mblibconfig "github.com/mocoarow/cocotola-1.25/moonbeam/lib/config"

	libconfig "github.com/mocoarow/cocotola-1.25/cocotola-lib/config"
	libmiddleware "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin/middleware"
)

type InitRouterGroupFunc func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc)

func InitRootRouterGroup(_ context.Context, corsConfig *mblibconfig.CORSConfig, logConfig *mblibconfig.LogConfig, debugConfig *libconfig.DebugConfig, appName string) *gin.Engine {
	if !debugConfig.Gin {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// cors
	ginCorsConfig := mblibconfig.InitCORS(corsConfig)

	router.Use(gin.Recovery())
	router.Use(cors.New(ginCorsConfig))
	router.Use(libmiddleware.PrometheusMiddleware())
	router.Use(otelgin.Middleware(appName, otelgin.WithFilter(func(req *http.Request) bool {
		return req.URL.Path != "/"
	})))

	if value, ok := logConfig.Enabled["accessLog"]; ok && value {
		withRequestBody := false
		if value, ok := logConfig.Enabled["accessLogRequestBody"]; ok && value {
			withRequestBody = true
		}
		withResponseBody := false
		if value, ok := logConfig.Enabled["accessLogResponseBody"]; ok && value {
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

	if debugConfig.Wait {
		router.Use(libmiddleware.NewWaitMiddleware(time.Second))
	}

	return router
}

func InitAPIRouterGroup(_ context.Context, parentRouterGroup gin.IRouter, appName string, logConfig *mblibconfig.LogConfig) *gin.RouterGroup {
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
