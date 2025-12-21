package config

import (
	"embed"
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"

	libconfig "github.com/mocoarow/cocotola-1.25/cocotola-lib/config"
	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller"
	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	authconfig "github.com/mocoarow/cocotola-1.25/cocotola-auth/config"
)

type ServerConfig struct {
	HTTPPort             int                           `yaml:"httpPort" validate:"required"`
	MetricsPort          int                           `yaml:"metricsPort" validate:"required"`
	ReadHeaderTimeoutSec int                           `yaml:"readHeaderTimeoutSec" validate:"gte=1"`
	Gin                  *libgin.Config                `yaml:"gin" validate:"required"`
	Shutdown             *libcontroller.ShutdownConfig `yaml:"shutdown" validate:"required"`
}

type AppConfig struct {
	Auth *authconfig.AuthConfig `yaml:"auth" validate:"required"`
}

type Config struct {
	App    *AppConfig              `yaml:"app" validate:"required"`
	Server *ServerConfig           `yaml:"server" validate:"required"`
	DB     *libgateway.DBConfig    `yaml:"db" validate:"required"`
	Trace  *libgateway.TraceConfig `yaml:"trace" validate:"required"`
	Log    *libgateway.LogConfig   `yaml:"log" validate:"required"`
}

//go:embed config.yml
var config embed.FS

func LoadConfig() (*Config, error) {
	filename := "config.yml"
	confContent, err := config.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config file(%s): %w", filename, err)
	}

	confContent = []byte(os.Expand(string(confContent), libconfig.ExpandEnvWithDefaults))
	var conf Config
	if err := yaml.Unmarshal(confContent, &conf); err != nil {
		return nil, fmt.Errorf("unmarshal file(%s): %w", filename, err)
	}

	if err := libdomain.Validator.Struct(&conf); err != nil {
		return nil, fmt.Errorf("validate file(%s): %w", filename, err)
	}

	return &conf, nil
}
