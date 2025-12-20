package gin

import (
	"github.com/gin-contrib/cors"
)

type CORSConfig struct {
	AllowOrigins []string `yaml:"allowOrigins" validate:"required"`
	AllowMethods []string `yaml:"allowMethods" validate:"required"`
	AllowHeaders []string `yaml:"allowHeaders"`
}

func InitCORS(cfg *CORSConfig) cors.Config {
	if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
		return cors.Config{ //nolint:exhaustruct
			AllowAllOrigins: true,
			AllowMethods:    cfg.AllowMethods,
			AllowHeaders:    cfg.AllowHeaders,
		}
	}

	return cors.Config{ //nolint:exhaustruct
		AllowOrigins: cfg.AllowOrigins,
		AllowMethods: cfg.AllowMethods,
		AllowHeaders: cfg.AllowHeaders,
	}
}
