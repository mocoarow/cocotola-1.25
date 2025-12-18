package config

import (
	"embed"
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	libconfig "github.com/mocoarow/cocotola-1.25/cocotola-lib/config"
)

type InitConfig struct {
	OwnerLoginID  string `yaml:"ownerLoginId" validate:"required"`
	OwnerPassword string `yaml:"ownerPassword" validate:"required"`
}

type Config struct {
	App   *InitConfig            `yaml:"app" validate:"required"`
	DB    *libconfig.DBConfig    `yaml:"db" validate:"required"`
	Trace *libconfig.TraceConfig `yaml:"trace" validate:"required"`
	Log   *libconfig.LogConfig   `yaml:"log" validate:"required"`
}

//go:embed config.yml
var config embed.FS

func LoadConfig() (*Config, error) {
	filename := "config.yml"
	confContent, err := config.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config.ReadFile. filename: %s, err: %w", filename, err)
	}

	confContent = []byte(os.Expand(string(confContent), libconfig.ExpandEnvWithDefaults))
	var conf Config
	if err := yaml.Unmarshal(confContent, &conf); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal. filename: %s, err: %w", filename, err)
	}

	if err := libdomain.Validator.Struct(&conf); err != nil {
		return nil, fmt.Errorf("Validator.Struct. filename: %s, err: %w", filename, err)
	}

	return &conf, nil
}
