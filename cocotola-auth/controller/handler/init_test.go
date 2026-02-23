package handler_test

import (
	libhandler "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/handler"
)

var (
	config libhandler.Config
)

func init() {
	config = libhandler.Config{
		CORS: &libhandler.CORSConfig{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Content-Type",
		},
		Log: &libhandler.LogConfig{
			AccessLog:             false,
			AccessLogRequestBody:  false,
			AccessLogResponseBody: false,
		},
		Debug: &libhandler.DebugConfig{
			Gin:  false,
			Wait: false,
		},
	}
}
