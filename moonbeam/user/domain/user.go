package domain

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type UserID struct {
	Value int `validate:"required,gte=0"`
}

func NewUserID(value int) (*UserID, error) {
	m := UserID{
		Value: value,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("validate user ID: %w", err)
	}

	return &m, nil
}

func (v *UserID) Int() int {
	return v.Value
}
