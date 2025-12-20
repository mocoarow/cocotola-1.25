package config

type OwnerConfig struct {
	OwnerLoginID  string `yaml:"ownerLoginId" validate:"required"`
	OwnerPassword string `yaml:"ownerPassword" validate:"required"`
}
