package gin_test

import (
	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
)

var (
	config libgin.Config
)

func init() {
	config = libgin.Config{
		CORS: &libgin.CORSConfig{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Content-Type",
		},
		Log: &libgin.LogConfig{
			AccessLog:             false,
			AccessLogRequestBody:  false,
			AccessLogResponseBody: false,
		},
		Debug: &libgin.DebugConfig{
			Gin:  false,
			Wait: false,
		},
	}
}
