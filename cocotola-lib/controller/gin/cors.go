package gin

import (
	"github.com/gin-contrib/cors"
)

type CORSConfig struct {
	AllowOrigins string `yaml:"allowOrigins" validate:"required"`
	AllowMethods string `yaml:"allowMethods" validate:"required"`
	AllowHeaders string `yaml:"allowHeaders"`
}

func InitCORS(cfg *CORSConfig) cors.Config {
	allowOrigins := SplitCommaSeparated(cfg.AllowOrigins)
	allowMethods := SplitCommaSeparated(cfg.AllowMethods)
	allowHeaders := SplitCommaSeparated(cfg.AllowHeaders)

	if len(allowOrigins) == 1 && allowOrigins[0] == "*" {
		return cors.Config{ //nolint:exhaustruct
			AllowAllOrigins: true,
			AllowMethods:    allowMethods,
			AllowHeaders:    allowHeaders,
		}
	}

	return cors.Config{ //nolint:exhaustruct
		AllowOrigins: allowOrigins,
		AllowMethods: allowMethods,
		AllowHeaders: allowHeaders,
	}
}
