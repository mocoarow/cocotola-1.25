package auth

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type AuthVerifyPasswordCommandGateway interface {
	service.UserRepositoryVerifyPassword
}

type AuthVerifyPasswordCommand struct {
	gw     AuthVerifyPasswordCommandGateway
	logger *slog.Logger
}

func NewAuthVerifyPasswordCommand(gw AuthVerifyPasswordCommandGateway) *AuthVerifyPasswordCommand {
	return &AuthVerifyPasswordCommand{
		gw:     gw,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "VerifyPasswordCommand")),
	}
}

func (u *AuthVerifyPasswordCommand) Execute(ctx context.Context, operator domain.SystemOwnerInterface, loginID, password string) error {
	ok, err := u.gw.VerifyPassword(ctx, operator, loginID, password)
	if err != nil {
		return fmt.Errorf("VerifyPassword: %w", err)
	}
	if !ok {
		return service.ErrUnauthenticated
	}
	return nil
}
