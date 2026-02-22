package auth

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type VerifyPasswordCommandGateway interface {
	service.UserRepositoryVerifyPassword
}

type VerifyPasswordCommand struct {
	gw     VerifyPasswordCommandGateway
	logger *slog.Logger
}

func NewVerifyPasswordCommand(gw VerifyPasswordCommandGateway) *VerifyPasswordCommand {
	return &VerifyPasswordCommand{
		gw:     gw,
		logger: slog.Default().With(slog.String(libdomain.LoggerNameKey, "VerifyPasswordCommand")),
	}
}

func (u *VerifyPasswordCommand) Execute(ctx context.Context, operator domain.SystemOwnerInterface, loginID, password string) error {
	ok, err := u.gw.VerifyPassword(ctx, operator, loginID, password)
	if err != nil {
		return fmt.Errorf("VerifyPassword: %w", err)
	}
	if !ok {
		return service.ErrUnauthenticated
	}
	return nil
}
