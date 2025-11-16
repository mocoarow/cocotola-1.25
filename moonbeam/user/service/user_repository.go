package service

import (
	"fmt"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type AddUserParameter struct {
	LoginID              string `validate:"required,max=255"`
	Username             string `validate:"required,max=255"`
	Password             string `validate:"required,min=8,max=255"`
	Provider             string
	ProviderLoginID      string
	ProviderAuthToken    string
	providerRefreshToken string
}

func NewAddUserParameter(loginID, username, password, provider, providerLoginID, providerAuthToken, providerRefreshToken string) (*AddUserParameter, error) {
	m := AddUserParameter{
		LoginID:              loginID,
		Username:             username,
		Password:             password,
		Provider:             provider,
		ProviderLoginID:      providerLoginID,
		ProviderAuthToken:    providerAuthToken,
		providerRefreshToken: providerRefreshToken,
	}
	if err := libdomain.Validator.Struct(&m); err != nil {
		return nil, fmt.Errorf("new add user parameter: %w", err)
	}

	return &m, nil
}
