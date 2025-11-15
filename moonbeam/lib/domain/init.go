package domain

import (
	"github.com/go-playground/validator/v10"
)

const (
	LoggerNameKey = "logger_name"
)

var (
	Validator = validator.New()
)
